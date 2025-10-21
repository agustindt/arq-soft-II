package services

import (
	"context"
	"errors"
	"fmt"
	"reservations/domain"
	"strings"
	"time"
)

// ReservasRepository define las operaciones de datos para Reservas
// Patrón Repository: abstrae el acceso a datos del resto de la aplicación
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

type ReservasServiceImpl struct {
	repository ReservasRepository // Inyección de dependencia
}

// NewReservasService crea una nueva instancia del service
// Pattern: Dependency Injection - recibe dependencies como parámetros
func NewReservasService(repository ReservasRepository) ReservasServiceImpl {
	return ReservasServiceImpl{
		repository: repository,
	}
}

// List obtiene todos los Reservas
// ✅ IMPLEMENTADO - Delegación simple al repository
func (s *ReservasServiceImpl) List(ctx context.Context) ([]domain.Reserva, error) {
	// En este caso, no hay lógica de negocio especial
	// Solo delegamos al repository
	return s.repository.List(ctx)
}

// Create valida y crea un nuevo Reserva
// Consigna 1: Validar name no vacío y price >= 0
func (s *ReservasServiceImpl) Create(ctx context.Context, Reserva domain.Reserva) (domain.Reserva, error) {
	// Validar campos del Reserva
	if err := s.validateReserva(Reserva); err != nil {
		return domain.Reserva{}, fmt.Errorf("validation error: %w", err)
	}

	created, err := s.repository.Create(ctx, Reserva)
	if err != nil {
		return domain.Reserva{}, fmt.Errorf("error creating Reserva in repository: %w", err)
	}

	return created, nil
}

// GetByID obtiene un Reserva por su ID
// Consigna 2: Validar formato de ID antes de consultar DB
func (s *ReservasServiceImpl) GetByID(ctx context.Context, id string) (domain.Reserva, error) {
	Reserva, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return domain.Reserva{}, fmt.Errorf("error getting Reserva from repository: %w", err)
	}

	return Reserva, nil
}

// Update actualiza un Reserva existente
// Consigna 3: Validar campos antes de actualizar
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

	return updated, nil
}

// Delete elimina un Reserva por ID
// Consigna 4: Validar ID antes de eliminar
func (s *ReservasServiceImpl) Delete(ctx context.Context, id string) error {
	_, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("reserva does not exists: %w", err)
	}

	err = s.repository.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("error deleting Reserva in repository: %w", err)
	}

	return nil
}

// validateReserva aplica reglas de negocio para validar un Reserva
func (s *ReservasServiceImpl) validateReserva(Reserva domain.Reserva) error {
	// 📝 Name es obligatorio y no puede estar vacío
	if strings.TrimSpace(Reserva.Actividad) == "" {
		return errors.New("name is required and cannot be empty")
	}
	// Date > date actual TODO
	if Reserva.Date.Before(time.Now()) {
		return errors.New("la fecha de la reserva debe ser posterior a la fecha actual")
	}

	// ✅ Todas las validaciones pasaron
	return nil
}
