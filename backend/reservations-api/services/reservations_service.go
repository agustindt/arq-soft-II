package services

import (
	"context"
	"errors"
	"fmt"
	"reservations/domain"
	"strings"
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
	// Validar campos del Reserva
	if err := s.validateReserva(Reserva); err != nil {
		return domain.Reserva{}, fmt.Errorf("validation error: %w", err)
	}

	created, err := s.repository.Create(ctx, Reserva)
	if err != nil {
		return domain.Reserva{}, fmt.Errorf("error creating Reserva in repository: %w", err)
	}

	if err := s.publisher.Publish(ctx, "create", created.ID); err != nil {
		return domain.Reserva{}, fmt.Errorf("error publishing Reserva creation: %w", err)
	}

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
	_, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return domain.Reserva{}, fmt.Errorf("reserva does not exists: %w", err)
	}

	// Validar campos del Reserva
	if err := s.validateReserva(Reserva); err != nil {
		return domain.Reserva{}, fmt.Errorf("validation error: %w", err)
	}

	// Actualizar en la BD
	updated, err := s.repository.Update(ctx, id, Reserva)
	if err != nil {
		return domain.Reserva{}, fmt.Errorf("error updating Reserva in repository: %w", err)
	}

	if err := s.publisher.Publish(ctx, "create", updated.ID); err != nil {
		return domain.Reserva{}, fmt.Errorf("error publishing Reserva creation: %w", err)
	}

	return updated, nil
}

// Delete elimina un Reserva por ID
func (s *ReservasServiceImpl) Delete(ctx context.Context, id string) error {
	_, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("reserva does not exists: %w", err)
	}

	err = s.repository.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("error deleting Reserva in repository: %w", err)
	}

	if err := s.publisher.Publish(ctx, "create", id); err != nil {
		return fmt.Errorf("error publishing Reserva creation: %w", err)
	}

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
	return nil
}
