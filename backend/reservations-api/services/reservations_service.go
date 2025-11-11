package services

import (
	"arq-soft-II/backend/reservations-api/domain"
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
)

type ReservasRepository interface {
	// List retorna todos los Reservas de la base de datos
	List(ctx context.Context) ([]domain.Reserva, error)

	// Create inserta un nuevo Reserva en DB
	Create(ctx context.Context, Reserva domain.Reserva) (domain.Reserva, error)

	// GetByID busca un Reserva por su ID
	GetByID(ctx context.Context, id string) (domain.Reserva, error)

	// Update actualiza un Reserva existente
	Update(ctx context.Context, id string, Reserva domain.Reserva) (domain.Reserva, error)

	// Delete elimina un Reserva por ID
	Delete(ctx context.Context, id string) error
} // ReservasServiceImpl implementa ReservasService

type ReservaPublisher interface {
	Publish(ctx context.Context, action string, reservaID string) error
}

type ReservasServiceImpl struct {
	repository ReservasRepository
	publisher  ReservaPublisher
}

func NewReservasService(repository ReservasRepository, publisher ReservaPublisher) ReservasServiceImpl {
	return ReservasServiceImpl{
		repository: repository,
		publisher:  publisher,
	}
}

// List obtiene todos los Reservas
func (s *ReservasServiceImpl) List(ctx context.Context) ([]domain.Reserva, error) {
	return s.repository.List(ctx)
}

// Create valida y crea un nuevo Reserva
func (s *ReservasServiceImpl) Create(ctx context.Context, Reserva domain.Reserva) (domain.Reserva, error) {
	// Ejecutar subtareas concurrentes: validación, cálculo de precio y enriquecimiento
	type taskResult struct {
		name string
		err  error
		data interface{}
	}

	tasksCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resultsCh := make(chan taskResult, 3)
	var wg sync.WaitGroup
	wg.Add(3)

	// Validación
	go func() {
		defer wg.Done()
		if err := s.validateReserva(Reserva); err != nil {
			resultsCh <- taskResult{name: "validate", err: err}
			return
		}
		resultsCh <- taskResult{name: "validate", err: nil}
	}()

	// Cálculo de precio (simulado)
	go func() {
		defer wg.Done()
		price, err := s.calculatePrice(tasksCtx, Reserva)
		resultsCh <- taskResult{name: "price", err: err, data: price}
	}()

	// Enriquecimiento de datos (simulado)
	go func() {
		defer wg.Done()
		note, err := s.enrichData(tasksCtx, Reserva)
		resultsCh <- taskResult{name: "enrich", err: err, data: note}
	}()

	// cerrar el channel cuando las tareas terminen
	go func() {
		wg.Wait()
		close(resultsCh)
	}()

	var finalPrice float64
	var finalNote string

	for tr := range resultsCh {
		if tr.err != nil {
			cancel()
			return domain.Reserva{}, fmt.Errorf("task %s failed: %w", tr.name, tr.err)
		}
		switch tr.name {
		case "price":
			if v, ok := tr.data.(float64); ok {
				finalPrice = v
			}
		case "enrich":
			if v, ok := tr.data.(string); ok {
				finalNote = v
			}
		}
	}

	if finalPrice > 0 {
		fmt.Printf("calculated price: %f\n", finalPrice)
	}
	if finalNote != "" {
		fmt.Printf("enrichment note: %s\n", finalNote)
	}

	created, err := s.repository.Create(ctx, Reserva)
	if err != nil {
		return domain.Reserva{}, fmt.Errorf("error creating Reserva in repository: %w", err)
	}

	go func(id string) {
		pubCtx, pubCancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer pubCancel()
		if err := s.publisher.Publish(pubCtx, "create", id); err != nil {
			fmt.Printf("warning: publish failed for reserva %s: %v\n", id, err)
		}
	}(created.ID)

	return created, nil
}

// GetByID obtiene un Reserva por su ID
func (s *ReservasServiceImpl) GetByID(ctx context.Context, id string) (domain.Reserva, error) {
	Reserva, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return domain.Reserva{}, fmt.Errorf("error getting Reserva from repository: %w", err)
	}

	return Reserva, nil
}

// Update actualiza un Reserva existente
func (s *ReservasServiceImpl) Update(ctx context.Context, id string, Reserva domain.Reserva) (domain.Reserva, error) {
	// Subtareas concurrentes: fetch existing and validate new data
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

	// fetch existing
	go func() {
		defer wg.Done()
		existing, err := s.repository.GetByID(tasksCtx, id)
		resultsCh <- taskResult{name: "fetch", err: err, data: existing}
	}()

	// validate new data
	go func() {
		defer wg.Done()
		if err := s.validateReserva(Reserva); err != nil {
			resultsCh <- taskResult{name: "validate", err: err}
			return
		}
		resultsCh <- taskResult{name: "validate", err: nil}
	}()

	go func() { wg.Wait(); close(resultsCh) }()

	var existing domain.Reserva
	for tr := range resultsCh {
		if tr.err != nil {
			cancel()
			return domain.Reserva{}, fmt.Errorf("task %s failed: %w", tr.name, tr.err)
		}
		if tr.name == "fetch" {
			if v, ok := tr.data.(domain.Reserva); ok {
				existing = v
			}
		}
	}

	// usar existing para evitar variable sin uso, y poder comparar si es necesario
	fmt.Printf("fetched existing reserva id=%s actividad=%s\n", existing.ID, existing.Actividad)

	// aquí podríamos comparar existing vs Reserva y aplicar lógicas de negocio
	updated, err := s.repository.Update(ctx, id, Reserva)
	if err != nil {
		return domain.Reserva{}, fmt.Errorf("error updating Reserva in repository: %w", err)
	}

	go func(id string) {
		pubCtx, pubCancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer pubCancel()
		if err := s.publisher.Publish(pubCtx, "update", id); err != nil {
			fmt.Printf("warning: publish failed for reserva %s: %v\n", id, err)
		}
	}(updated.ID)

	return updated, nil
}

// Delete elimina un Reserva por ID
func (s *ReservasServiceImpl) Delete(ctx context.Context, id string) error {
	// Subtareas concurrentes: fetch existing (to check rules) y prepare audit note
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
		// preparar nota de auditoría (simulada)
		select {
		case <-time.After(20 * time.Millisecond):
			resultsCh <- taskResult{name: "audit", err: nil, data: "ready"}
		case <-tasksCtx.Done():
			resultsCh <- taskResult{name: "audit", err: tasksCtx.Err()}
		}
	}()

	go func() { wg.Wait(); close(resultsCh) }()

	var existing domain.Reserva
	for tr := range resultsCh {
		if tr.err != nil {
			cancel()
			return fmt.Errorf("task %s failed: %w", tr.name, tr.err)
		}
		if tr.name == "fetch" {
			if v, ok := tr.data.(domain.Reserva); ok {
				existing = v
			}
		}
	}

	// ejemplo de regla: no permitir eliminar si status es confirmada
	if existing.Status == "confirmada" {
		return fmt.Errorf("cannot delete a confirmed reserva")
	}

	if err := s.repository.Delete(ctx, id); err != nil {
		return fmt.Errorf("error deleting Reserva in repository: %w", err)
	}

	go func(id string) {
		pubCtx, pubCancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer pubCancel()
		if err := s.publisher.Publish(pubCtx, "delete", id); err != nil {
			fmt.Printf("warning: publish failed for reserva %s: %v\n", id, err)
		}
	}(id)

	return nil
}

// validateReserva aplica reglas de negocio para validar un Reserva
func (s *ReservasServiceImpl) validateReserva(Reserva domain.Reserva) error {
	if strings.TrimSpace(Reserva.Actividad) == "" {
		return errors.New("name is required and cannot be empty")
	}
	if Reserva.Date.Before(time.Now()) {
		return errors.New("la fecha de la reserva debe ser posterior a la fecha actual")
	}
	if Reserva.Cupo < 0 {
		return errors.New("el cupo no puede ser menor a cero")
	}
	if len(Reserva.UsersID) == 0 {
		return errors.New("debe haber al menos un usuario en la reserva")
	}
	return nil
}

// calculatePrice realiza algún cálculo simulado de precio y respeta el context
func (s *ReservasServiceImpl) calculatePrice(ctx context.Context, r domain.Reserva) (float64, error) {
	base := 10.0
	select {
	case <-time.After(100 * time.Millisecond):
		// continue
	case <-ctx.Done():
		return 0, ctx.Err()
	}
	if strings.Contains(strings.ToLower(r.Actividad), "premium") {
		base += 20.0
	}
	return base, nil
}

// enrichData simula una llamada externa para enriquecer la reserva
func (s *ReservasServiceImpl) enrichData(ctx context.Context, r domain.Reserva) (string, error) {
	select {
	case <-time.After(50 * time.Millisecond):
		// continue
	case <-ctx.Done():
		return "", ctx.Err()
	}
	return fmt.Sprintf("Reserva para %s procesada el %s", r.Actividad, time.Now().Format(time.RFC3339)), nil
}
