package services

import (
	"arq-soft-II/backend/activities-api/domain"
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
)

type ActivitiesRepository interface {
	List(ctx context.Context) ([]domain.Activity, error)
	ListAll(ctx context.Context) ([]domain.Activity, error)
	Create(ctx context.Context, activity domain.Activity) (domain.Activity, error)
	GetByID(ctx context.Context, id string) (domain.Activity, error)
	Update(ctx context.Context, id string, activity domain.Activity) (domain.Activity, error)
	Delete(ctx context.Context, id string) error
	HardDelete(ctx context.Context, id string) error
	ToggleActive(ctx context.Context, id string) (domain.Activity, error)
	GetByCategory(ctx context.Context, category string) ([]domain.Activity, error)
}

type ActivityPublisher interface {
	Publish(ctx context.Context, action string, activityID string) error
}

type ActivitiesServiceImpl struct {
	repository ActivitiesRepository
	publisher  ActivityPublisher
}

func NewActivitiesService(repository ActivitiesRepository, publisher ActivityPublisher) *ActivitiesServiceImpl {
	return &ActivitiesServiceImpl{
		repository: repository,
		publisher:  publisher,
	}
}

// List obtiene todas las actividades activas
func (s *ActivitiesServiceImpl) List(ctx context.Context) ([]domain.Activity, error) {
	return s.repository.List(ctx)
}

// ListAll obtiene todas las actividades (incluyendo inactivas) - solo para admin
func (s *ActivitiesServiceImpl) ListAll(ctx context.Context) ([]domain.Activity, error) {
	return s.repository.ListAll(ctx)
}

// Create valida y crea una nueva actividad
func (s *ActivitiesServiceImpl) Create(ctx context.Context, activity domain.Activity) (domain.Activity, error) {
	// Ejecutar subtareas concurrentes: validación y enriquecimiento
	type taskResult struct {
		name string
		err  error
		data interface{}
	}

	tasksCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resultsCh := make(chan taskResult, 2)
	var wg sync.WaitGroup
	wg.Add(2)

	// Validación
	go func() {
		defer wg.Done()
		if err := s.validateActivity(activity); err != nil {
			resultsCh <- taskResult{name: "validate", err: err}
			return
		}
		resultsCh <- taskResult{name: "validate", err: nil}
	}()

	// Enriquecimiento de datos
	go func() {
		defer wg.Done()
		note, err := s.enrichActivity(tasksCtx, activity)
		resultsCh <- taskResult{name: "enrich", err: err, data: note}
	}()

	// Cerrar el channel cuando las tareas terminen
	go func() {
		wg.Wait()
		close(resultsCh)
	}()

	var enrichmentNote string

	for tr := range resultsCh {
		if tr.err != nil {
			cancel()
			return domain.Activity{}, fmt.Errorf("task %s failed: %w", tr.name, tr.err)
		}
		if tr.name == "enrich" {
			if v, ok := tr.data.(string); ok {
				enrichmentNote = v
			}
		}
	}

	if enrichmentNote != "" {
		fmt.Printf("📝 Enrichment note: %s\n", enrichmentNote)
	}

	// Crear la actividad en el repositorio
	created, err := s.repository.Create(ctx, activity)
	if err != nil {
		return domain.Activity{}, fmt.Errorf("error creating activity in repository: %w", err)
	}

	// Publicar evento de creación de forma asíncrona
	go func(id string) {
		pubCtx, pubCancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer pubCancel()
		if err := s.publisher.Publish(pubCtx, "created", id); err != nil {
			fmt.Printf("⚠️  Warning: publish failed for activity %s: %v\n", id, err)
		}
	}(created.ID)

	return created, nil
}

// GetByID obtiene una actividad por su ID
func (s *ActivitiesServiceImpl) GetByID(ctx context.Context, id string) (domain.Activity, error) {
	activity, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return domain.Activity{}, fmt.Errorf("error getting activity from repository: %w", err)
	}

	return activity, nil
}

// Update actualiza una actividad existente
func (s *ActivitiesServiceImpl) Update(ctx context.Context, id string, activity domain.Activity) (domain.Activity, error) {
	// Subtareas concurrentes: fetch existing y validate new data
	type taskResult struct {
		name string
		err  error
		data interface{}
	}

	tasksCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resultsCh := make(chan taskResult, 2)
	var wg sync.WaitGroup
	wg.Add(2)

	// Fetch existing
	go func() {
		defer wg.Done()
		existing, err := s.repository.GetByID(tasksCtx, id)
		resultsCh <- taskResult{name: "fetch", err: err, data: existing}
	}()

	// Validate new data
	go func() {
		defer wg.Done()
		if err := s.validateActivity(activity); err != nil {
			resultsCh <- taskResult{name: "validate", err: err}
			return
		}
		resultsCh <- taskResult{name: "validate", err: nil}
	}()

	go func() {
		wg.Wait()
		close(resultsCh)
	}()

	var existing domain.Activity
	for tr := range resultsCh {
		if tr.err != nil {
			cancel()
			return domain.Activity{}, fmt.Errorf("task %s failed: %w", tr.name, tr.err)
		}
		if tr.name == "fetch" {
			if v, ok := tr.data.(domain.Activity); ok {
				existing = v
			}
		}
	}

	fmt.Printf("📋 Updating activity: id=%s, name=%s -> %s\n", existing.ID, existing.Name, activity.Name)

	// Actualizar en el repositorio
	updated, err := s.repository.Update(ctx, id, activity)
	if err != nil {
		return domain.Activity{}, fmt.Errorf("error updating activity in repository: %w", err)
	}

	// Publicar evento de actualización de forma asíncrona
	go func(id string) {
		pubCtx, pubCancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer pubCancel()
		if err := s.publisher.Publish(pubCtx, "updated", id); err != nil {
			fmt.Printf("⚠️  Warning: publish failed for activity %s: %v\n", id, err)
		}
	}(updated.ID)

	return updated, nil
}

// Delete elimina (soft delete) una actividad
func (s *ActivitiesServiceImpl) Delete(ctx context.Context, id string) error {
	// Fetch existing para validar
	tasksCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	existing, err := s.repository.GetByID(tasksCtx, id)
	if err != nil {
		return fmt.Errorf("error fetching activity: %w", err)
	}

	fmt.Printf("🗑️  Soft deleting activity: id=%s, name=%s\n", existing.ID, existing.Name)

	// Soft delete
	if err := s.repository.Delete(ctx, id); err != nil {
		return fmt.Errorf("error deleting activity in repository: %w", err)
	}

	// Publicar evento de eliminación de forma asíncrona
	go func(id string) {
		pubCtx, pubCancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer pubCancel()
		if err := s.publisher.Publish(pubCtx, "deleted", id); err != nil {
			fmt.Printf("⚠️  Warning: publish failed for activity %s: %v\n", id, err)
		}
	}(id)

	return nil
}

// HardDelete elimina permanentemente una actividad de la base de datos
func (s *ActivitiesServiceImpl) HardDelete(ctx context.Context, id string) error {
	// Fetch existing para validar
	tasksCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	existing, err := s.repository.GetByID(tasksCtx, id)
	if err != nil {
		return fmt.Errorf("error fetching activity: %w", err)
	}

	fmt.Printf("🗑️  PERMANENTLY deleting activity: id=%s, name=%s\n", existing.ID, existing.Name)

	// Hard delete - eliminación permanente
	if err := s.repository.HardDelete(ctx, id); err != nil {
		return fmt.Errorf("error hard deleting activity in repository: %w", err)
	}

	// Publicar evento de eliminación de forma asíncrona
	go func(id string) {
		pubCtx, pubCancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer pubCancel()
		if err := s.publisher.Publish(pubCtx, "deleted", id); err != nil {
			fmt.Printf("⚠️  Warning: publish failed for activity %s: %v\n", id, err)
		}
	}(id)

	return nil
}

// ToggleActive activa/desactiva una actividad
func (s *ActivitiesServiceImpl) ToggleActive(ctx context.Context, id string) (domain.Activity, error) {
	toggled, err := s.repository.ToggleActive(ctx, id)
	if err != nil {
		return domain.Activity{}, fmt.Errorf("error toggling activity status: %w", err)
	}

	// Publicar evento de actualización
	go func(id string) {
		pubCtx, pubCancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer pubCancel()
		action := "updated"
		if err := s.publisher.Publish(pubCtx, action, id); err != nil {
			fmt.Printf("⚠️  Warning: publish failed for activity %s: %v\n", id, err)
		}
	}(toggled.ID)

	return toggled, nil
}

// GetByCategory obtiene actividades por categoría
func (s *ActivitiesServiceImpl) GetByCategory(ctx context.Context, category string) ([]domain.Activity, error) {
	if strings.TrimSpace(category) == "" {
		return nil, errors.New("category cannot be empty")
	}

	activities, err := s.repository.GetByCategory(ctx, category)
	if err != nil {
		return nil, fmt.Errorf("error getting activities by category: %w", err)
	}

	return activities, nil
}

// validateActivity aplica reglas de negocio para validar una actividad
func (s *ActivitiesServiceImpl) validateActivity(activity domain.Activity) error {
	if strings.TrimSpace(activity.Name) == "" {
		return errors.New("name is required and cannot be empty")
	}

	if strings.TrimSpace(activity.Description) == "" {
		return errors.New("description is required and cannot be empty")
	}

	if strings.TrimSpace(activity.Category) == "" {
		return errors.New("category is required and cannot be empty")
	}

	// Validar difficulty
	validDifficulties := map[string]bool{
		"beginner":     true,
		"intermediate": true,
		"advanced":     true,
	}
	if !validDifficulties[strings.ToLower(activity.Difficulty)] {
		return errors.New("difficulty must be one of: beginner, intermediate, advanced")
	}

	if activity.Duration <= 0 {
		return errors.New("duration must be greater than 0")
	}

	if activity.MaxCapacity <= 0 {
		return errors.New("max capacity must be greater than 0")
	}

	if activity.Price < 0 {
		return errors.New("price cannot be negative")
	}

	if strings.TrimSpace(activity.Location) == "" {
		return errors.New("location is required and cannot be empty")
	}

	return nil
}

// enrichActivity simula enriquecimiento de datos de la actividad
func (s *ActivitiesServiceImpl) enrichActivity(ctx context.Context, activity domain.Activity) (string, error) {
	select {
	case <-time.After(50 * time.Millisecond):
		// continue
	case <-ctx.Done():
		return "", ctx.Err()
	}

	return fmt.Sprintf("Activity '%s' processed at %s", activity.Name, time.Now().Format(time.RFC3339)), nil
}
