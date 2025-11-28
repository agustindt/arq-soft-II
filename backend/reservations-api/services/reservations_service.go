package services

import (
	"arq-soft-II/backend/reservations-api/clients"
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

	// ListByUserID retorna las reservas de un usuario específico
	ListByUserID(ctx context.Context, userID int) ([]domain.Reserva, error)

	// Create inserta un nuevo Reserva en DB
	Create(ctx context.Context, Reserva domain.Reserva) (domain.Reserva, error)

	// GetByID busca un Reserva por su ID
	GetByID(ctx context.Context, id string) (domain.Reserva, error)

	// Update actualiza un Reserva existente
	Update(ctx context.Context, id string, Reserva domain.Reserva) (domain.Reserva, error)

	// Delete elimina un Reserva por ID
	Delete(ctx context.Context, id string) error

	// CountActiveReservationsBySchedule cuenta las reservas activas para un horario específico
	CountActiveReservationsBySchedule(ctx context.Context, activityID string, schedule string, date time.Time) (int, error)

	// ExistsActiveReservation verifica si ya existe una reserva activa para un usuario, actividad, horario y fecha
	ExistsActiveReservation(ctx context.Context, userID int, activityID string, schedule string, date time.Time) (bool, error)

	// ExistsScheduleConflict verifica si el usuario ya tiene una reserva activa en el mismo horario (cualquier actividad)
	ExistsScheduleConflict(ctx context.Context, userID int, schedule string, date time.Time) (bool, string, error)

	// GetUserActiveReservationsByDate obtiene todas las reservas activas de un usuario en una fecha específica
	GetUserActiveReservationsByDate(ctx context.Context, userID int, date time.Time) ([]domain.Reserva, error)
} // ReservasServiceImpl implementa ReservasService

type ReservaPublisher interface {
	Publish(ctx context.Context, action string, reservaID string) error
}

type ReservasServiceImpl struct {
	repository       ReservasRepository
	publisher        ReservaPublisher
	activitiesAPIURL string
}

func NewReservasService(repository ReservasRepository, publisher ReservaPublisher, activitiesAPIURL string) ReservasServiceImpl {
	return ReservasServiceImpl{
		repository:       repository,
		publisher:        publisher,
		activitiesAPIURL: activitiesAPIURL,
	}
}

// List obtiene todos los Reservas (solo para admin)
func (s *ReservasServiceImpl) List(ctx context.Context) ([]domain.Reserva, error) {
	return s.repository.List(ctx)
}

// ListByUserID obtiene las reservas de un usuario específico
func (s *ReservasServiceImpl) ListByUserID(ctx context.Context, userID int) ([]domain.Reserva, error) {
	return s.repository.ListByUserID(ctx, userID)
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
		if err := s.validateReserva(tasksCtx, Reserva); err != nil {
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
		if err := s.validateReserva(tasksCtx, Reserva); err != nil {
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
func (s *ReservasServiceImpl) validateReserva(ctx context.Context, Reserva domain.Reserva) error {
	if strings.TrimSpace(Reserva.Actividad) == "" {
		return errors.New("actividad is required and cannot be empty")
	}
	// Verificar que la fecha no sea zero time
	if Reserva.Date.IsZero() {
		return errors.New("date is required and cannot be empty")
	}
	// Verificar que la fecha sea posterior a la fecha actual (con un margen de 1 minuto para evitar problemas de tiempo)
	now := time.Now()
	oneMinuteAgo := now.Add(-1 * time.Minute)
	if Reserva.Date.Before(oneMinuteAgo) {
		return errors.New("la fecha de la reserva debe ser posterior a la fecha actual")
	}
	if Reserva.Cupo <= 0 {
		return errors.New("el cupo debe ser mayor a cero")
	}
	if len(Reserva.UsersID) == 0 {
		return errors.New("debe haber al menos un usuario en la reserva")
	}

	// Validar schedule
	if strings.TrimSpace(Reserva.Schedule) == "" {
		return errors.New("schedule is required and cannot be empty")
	}

	// Obtener la actividad desde la API de actividades (necesitamos la duración para verificar solapamiento)
	activity, err := clients.GetActivityByID(s.activitiesAPIURL, Reserva.Actividad)
	if err != nil {
		return fmt.Errorf("error fetching activity: %w", err)
	}

	// Verificar conflictos de horario considerando solapamiento de tiempo
	if len(Reserva.UsersID) > 0 {
		for _, userID := range Reserva.UsersID {
			// Obtener todas las reservas activas del usuario en esa fecha
			existingReservations, err := s.repository.GetUserActiveReservationsByDate(ctx, userID, Reserva.Date)
			if err != nil {
				return fmt.Errorf("error checking existing reservations: %w", err)
			}

			fmt.Printf("DEBUG: Found %d existing reservations for user %d on date %s\n",
				len(existingReservations), userID, Reserva.Date.Format("2006-01-02"))

			// Verificar si alguna reserva existente se solapa con la nueva
			for _, existingReservation := range existingReservations {
				fmt.Printf("DEBUG: Checking overlap - New: %s (%s, %d min) vs Existing: %s (%s, activity: %s)\n",
					Reserva.Schedule, Reserva.Actividad, activity.Duration,
					existingReservation.Schedule, existingReservation.Actividad, existingReservation.Actividad)

				// Verificar si es la misma actividad y horario (duplicado exacto)
				if existingReservation.Actividad == Reserva.Actividad && existingReservation.Schedule == Reserva.Schedule {
					return fmt.Errorf("ya existe una reserva activa para esta actividad, horario y fecha")
				}

				// Obtener la actividad de la reserva existente para conocer su duración
				existingActivity, err := clients.GetActivityByID(s.activitiesAPIURL, existingReservation.Actividad)
				if err != nil {
					fmt.Printf("DEBUG: Error fetching existing activity %s: %v\n", existingReservation.Actividad, err)
					// Si no podemos obtener la actividad existente, continuamos (puede haber sido eliminada)
					continue
				}

				fmt.Printf("DEBUG: Existing activity duration: %d minutes\n", existingActivity.Duration)

				// Verificar solapamiento de horarios
				overlaps, err := schedulesOverlap(
					Reserva.Schedule, activity.Duration,
					existingReservation.Schedule, existingActivity.Duration,
				)
				if err != nil {
					fmt.Printf("DEBUG: Error checking overlap: %v\n", err)
					return fmt.Errorf("error checking schedule overlap: %w", err)
				}

				fmt.Printf("DEBUG: Overlaps? %v\n", overlaps)

				if overlaps {
					return fmt.Errorf("el horario '%s' se solapa con tu reserva existente de '%s' en el horario '%s'. No puedes estar en dos lugares al mismo tiempo",
						Reserva.Schedule, existingActivity.Name, existingReservation.Schedule)
				}
			}
		}
	}

	// Verificar que el schedule existe en la actividad
	scheduleExists := false
	for _, sch := range activity.Schedule {
		if sch == Reserva.Schedule {
			scheduleExists = true
			break
		}
	}
	if !scheduleExists {
		return fmt.Errorf("schedule '%s' does not exist for this activity", Reserva.Schedule)
	}

	// Verificar capacidad disponible para el horario específico
	currentReservations, err := s.repository.CountActiveReservationsBySchedule(
		ctx, Reserva.Actividad, Reserva.Schedule, Reserva.Date,
	)
	if err != nil {
		return fmt.Errorf("error checking capacity: %w", err)
	}

	availableCapacity := activity.MaxCapacity - currentReservations
	if Reserva.Cupo > availableCapacity {
		return fmt.Errorf("insufficient capacity: requested %d, available %d", Reserva.Cupo, availableCapacity)
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

// GetScheduleAvailability retorna la disponibilidad de cada horario para una actividad en una fecha
func (s *ReservasServiceImpl) GetScheduleAvailability(ctx context.Context, activityID string, date time.Time) (map[string]int, error) {
	// Obtener la actividad desde la API de actividades
	activity, err := clients.GetActivityByID(s.activitiesAPIURL, activityID)
	if err != nil {
		return nil, fmt.Errorf("error fetching activity: %w", err)
	}

	// Calcular disponibilidad para cada horario
	availability := make(map[string]int)
	for _, schedule := range activity.Schedule {
		currentReservations, err := s.repository.CountActiveReservationsBySchedule(
			ctx, activityID, schedule, date,
		)
		if err != nil {
			// En caso de error, asumimos 0 disponibles
			availability[schedule] = 0
			continue
		}
		availability[schedule] = activity.MaxCapacity - currentReservations
	}

	return availability, nil
}

// parseScheduleTime extrae las horas y minutos de un horario como "Lunes 20:00"
// Retorna los minutos desde medianoche (ej: 20:00 = 1200 minutos)
func parseScheduleTime(schedule string) (int, error) {
	// El formato es "Día HH:MM", ej: "Lunes 20:00"
	parts := strings.Split(schedule, " ")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid schedule format: %s", schedule)
	}

	timePart := parts[1] // "20:00"
	timeParts := strings.Split(timePart, ":")
	if len(timeParts) != 2 {
		return 0, fmt.Errorf("invalid time format: %s", timePart)
	}

	hour := 0
	minute := 0

	// Convertir hora y minutos a enteros
	if h, err := fmt.Sscanf(timeParts[0], "%d", &hour); err != nil || h != 1 {
		return 0, fmt.Errorf("invalid hour: %s", timeParts[0])
	}

	if m, err := fmt.Sscanf(timeParts[1], "%d", &minute); err != nil || m != 1 {
		return 0, fmt.Errorf("invalid minute: %s", timeParts[1])
	}

	// Convertir a minutos desde medianoche
	totalMinutes := hour*60 + minute
	return totalMinutes, nil
}

// extractDayFromSchedule extrae el día de la semana de un horario como "Lunes 20:00"
// Normaliza el día eliminando acentos para hacer la comparación más robusta
func extractDayFromSchedule(schedule string) string {
	parts := strings.Split(schedule, " ")
	if len(parts) >= 1 {
		day := parts[0]
		// Normalizar: eliminar acentos y convertir a minúsculas
		day = strings.ToLower(day)
		day = strings.ReplaceAll(day, "á", "a")
		day = strings.ReplaceAll(day, "é", "e")
		day = strings.ReplaceAll(day, "í", "i")
		day = strings.ReplaceAll(day, "ó", "o")
		day = strings.ReplaceAll(day, "ú", "u")
		return day
	}
	return ""
}

// schedulesOverlap verifica si dos horarios se solapan considerando sus duraciones
// schedule1/duration1: primer horario y su duración en minutos
// schedule2/duration2: segundo horario y su duración en minutos
func schedulesOverlap(schedule1 string, duration1 int, schedule2 string, duration2 int) (bool, error) {
	// Verificar si son el mismo día de la semana
	day1 := extractDayFromSchedule(schedule1)
	day2 := extractDayFromSchedule(schedule2)

	if day1 != day2 {
		// Días diferentes, no hay solapamiento
		return false, nil
	}

	// Parsear los horarios (minutos desde medianoche)
	start1, err := parseScheduleTime(schedule1)
	if err != nil {
		return false, err
	}

	start2, err := parseScheduleTime(schedule2)
	if err != nil {
		return false, err
	}

	// Calcular los rangos de tiempo
	end1 := start1 + duration1
	end2 := start2 + duration2

	// Verificar solapamiento: dos rangos [a,b] y [c,d] se solapan si a < d AND c < b
	overlap := start1 < end2 && start2 < end1

	return overlap, nil
}
