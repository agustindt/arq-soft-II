package services

import (
	"activities-api/clients"
	"activities-api/domain"
	"activities-api/utils"
	"context"
	"errors"
	"fmt"
	"os"
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
	usersAPI   string
}

func NewActivitiesService(repository ActivitiesRepository, publisher ActivityPublisher) *ActivitiesServiceImpl {
	usersAPI := os.Getenv("USERS_API_URL")
	if usersAPI == "" {
		usersAPI = "http://localhost:8081"
	}

	return &ActivitiesServiceImpl{
		repository: repository,
		publisher:  publisher,
		usersAPI:   usersAPI,
	}
}

var (
	ErrOwnerNotFound  = errors.New("owner_not_found")
	ErrOwnerForbidden = errors.New("owner_mismatch")
)

func (s *ActivitiesServiceImpl) requesterFromContext(ctx context.Context) (uint, string) {
	uid, _ := ctx.Value(utils.ContextUserIDKey).(uint)
	role, _ := ctx.Value(utils.ContextUserRoleKey).(string)
	return uid, role
}

func (s *ActivitiesServiceImpl) validateOwner(ctx context.Context, ownerID uint, requesterID uint, requesterRole string) error {
	if requesterRole == "admin" || requesterRole == "root" || requesterRole == "super_admin" {
		return nil
	}

	if ownerID == 0 {
		return ErrOwnerNotFound
	}

	userCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	user, err := clients.GetUserByIDWithContext(userCtx, s.usersAPI, ownerID)
	if err != nil || user == nil {
		return ErrOwnerNotFound
	}

	if user.ID != requesterID {
		return ErrOwnerForbidden
	}

	return nil
}

// List obtiene todas las actividades activas
func (s *ActivitiesServiceImpl) List(ctx context.Context) ([]domain.Activity, error) {
	return s.repository.List(ctx)
}

// ListAll obtiene todas las actividades (incluyendo inactivas)
func (s *ActivitiesServiceImpl) ListAll(ctx context.Context) ([]domain.Activity, error) {
	return s.repository.ListAll(ctx)
}

// Create valida y crea una nueva actividad
func (s *ActivitiesServiceImpl) Create(ctx context.Context, activity domain.Activity) (domain.Activity, error) {
	requesterID, requesterRole := s.requesterFromContext(ctx)
	if activity.CreatedBy == 0 {
		activity.CreatedBy = requesterID
	}

	if err := s.validateOwner(ctx, activity.CreatedBy, requesterID, requesterRole); err != nil {
		return domain.Activity{}, err
	}

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

	go func() {
		defer wg.Done()
		if err := s.validateActivity(activity); err != nil {
			resultsCh <- taskResult{name: "validate", err: err}
			return
		}
		resultsCh <- taskResult{name: "validate", err: nil}
	}()

	go func() {
		defer wg.Done()
		if err := s.enrichActivity(&activity); err != nil {
			resultsCh <- taskResult{name: "enrich", err: err}
			return
		}
		resultsCh <- taskResult{name: "enrich", err: nil}
	}()

	go func() {
		wg.Wait()
		close(resultsCh)
	}()

	for tr := range resultsCh {
		if tr.err != nil {
			cancel()
			return domain.Activity{}, fmt.Errorf("task %s failed: %w", tr.name, tr.err)
		}
	}

	created, err := s.repository.Create(ctx, activity)
	if err != nil {
		return domain.Activity{}, fmt.Errorf("error creating activity in repository: %w", err)
	}

	go func(id string) {
		pubCtx, pubCancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer pubCancel()
		if err := s.publisher.Publish(pubCtx, "created", id); err != nil {
			fmt.Printf("⚠ Warning: publish failed for activity %s: %v\n", id, err)
		}
	}(created.ID)

	return created, nil
}

// Update actualiza una actividad existente
func (s *ActivitiesServiceImpl) Update(ctx context.Context, id string, activity domain.Activity) (domain.Activity, error) {
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

	go func() {
		defer wg.Done()
		existing, err := s.repository.GetByID(tasksCtx, id)
		resultsCh <- taskResult{name: "fetch", err: err, data: existing}
	}()

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

	requesterID, requesterRole := s.requesterFromContext(ctx)
	if err := s.validateOwner(ctx, existing.CreatedBy, requesterID, requesterRole); err != nil {
		return domain.Activity{}, err
	}

	activity.CreatedBy = existing.CreatedBy

	fmt.Printf("📋 Updating activity: id=%s, name=%s -> %s\n", existing.ID, existing.Name, activity.Name)

	updated, err := s.repository.Update(ctx, id, activity)
	if err != nil {
		return domain.Activity{}, fmt.Errorf("error updating activity in repository: %w", err)
	}

	go func(id string) {
		pubCtx, pubCancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer pubCancel()
		if err := s.publisher.Publish(pubCtx, "updated", id); err != nil {
			fmt.Printf("⚠ Warning: publish failed for activity %s: %v\n", id, err)
		}
	}(updated.ID)

	return updated, nil
}

// Delete elimina una actividad (soft delete)
func (s *ActivitiesServiceImpl) Delete(ctx context.Context, id string) error {
	tasksCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	existing, err := s.repository.GetByID(tasksCtx, id)
	if err != nil {
		return fmt.Errorf("error fetching activity: %w", err)
	}

	requesterID, requesterRole := s.requesterFromContext(ctx)
	if err := s.validateOwner(ctx, existing.CreatedBy, requesterID, requesterRole); err != nil {
		return err
	}

	fmt.Printf("🗑️ Soft deleting activity: id=%s, name=%s\n", existing.ID, existing.Name)

	if err := s.repository.Delete(ctx, id); err != nil {
		return fmt.Errorf("error deleting activity in repository: %w", err)
	}

	go func(id string) {
		pubCtx, pubCancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer pubCancel()
		if err := s.publisher.Publish(pubCtx, "deleted", id); err != nil {
			fmt.Printf("⚠ Warning: publish failed for activity %s: %v\n", id, err)
		}
	}(id)

	return nil
}

// HardDelete elimina permanentemente una actividad
func (s *ActivitiesServiceImpl) HardDelete(ctx context.Context, id string) error {
	tasksCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	existing, err := s.repository.GetByID(tasksCtx, id)
	if err != nil {
		return fmt.Errorf("error fetching activity: %w", err)
	}

	requesterID, requesterRole := s.requesterFromContext(ctx)
	if err := s.validateOwner(ctx, existing.CreatedBy, requesterID, requesterRole); err != nil {
		return err
	}

	fmt.Printf("🗑️ PERMANENTLY deleting activity: id=%s, name=%s\n", existing.ID, existing.Name)

	if err := s.repository.HardDelete(ctx, id); err != nil {
		return fmt.Errorf("error hard deleting activity in repository: %w", err)
	}

	go func(id string) {
		pubCtx, pubCancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer pubCancel()
		if err := s.publisher.Publish(pubCtx, "deleted", id); err != nil {
			fmt.Printf("⚠ Warning: publish failed for activity %s: %v\n", id, err)
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

	go func(id string) {
		pubCtx, pubCancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer pubCancel()
		if err := s.publisher.Publish(pubCtx, "updated", id); err != nil {
			fmt.Printf("⚠ Warning: publish failed for activity %s: %v\n", id, err)
		}
	}(id)

	return toggled, nil
}
