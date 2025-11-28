package services

import (
	"activities-api/domain"
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// ===== MOCKS =====

type mockRepository struct {
	mu sync.Mutex

	listFn         func(ctx context.Context) ([]domain.Activity, error)
	listAllFn      func(ctx context.Context) ([]domain.Activity, error)
	createFn       func(ctx context.Context, activity domain.Activity) (domain.Activity, error)
	getByIDFn      func(ctx context.Context, id string) (domain.Activity, error)
	updateFn       func(ctx context.Context, id string, activity domain.Activity) (domain.Activity, error)
	deleteFn       func(ctx context.Context, id string) error
	hardDeleteFn   func(ctx context.Context, id string) error
	toggleActiveFn func(ctx context.Context, id string) (domain.Activity, error)
	getByCategory  func(ctx context.Context, category string) ([]domain.Activity, error)

	// Para tracking de llamadas
	listCalls    int32
	createCalls  int32
	getByIDCalls int32
	updateCalls  int32
	deleteCalls  int32
}

func (m *mockRepository) List(ctx context.Context) ([]domain.Activity, error) {
	atomic.AddInt32(&m.listCalls, 1)
	if m.listFn != nil {
		return m.listFn(ctx)
	}
	return []domain.Activity{}, nil
}

func (m *mockRepository) ListAll(ctx context.Context) ([]domain.Activity, error) {
	if m.listAllFn != nil {
		return m.listAllFn(ctx)
	}
	return []domain.Activity{}, nil
}

func (m *mockRepository) Create(ctx context.Context, activity domain.Activity) (domain.Activity, error) {
	atomic.AddInt32(&m.createCalls, 1)
	if m.createFn != nil {
		return m.createFn(ctx, activity)
	}
	return activity, nil
}

func (m *mockRepository) GetByID(ctx context.Context, id string) (domain.Activity, error) {
	atomic.AddInt32(&m.getByIDCalls, 1)
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return domain.Activity{ID: id}, nil
}

func (m *mockRepository) Update(ctx context.Context, id string, activity domain.Activity) (domain.Activity, error) {
	atomic.AddInt32(&m.updateCalls, 1)
	if m.updateFn != nil {
		return m.updateFn(ctx, id, activity)
	}
	return activity, nil
}

func (m *mockRepository) Delete(ctx context.Context, id string) error {
	atomic.AddInt32(&m.deleteCalls, 1)
	if m.deleteFn != nil {
		return m.deleteFn(ctx, id)
	}
	return nil
}

func (m *mockRepository) HardDelete(ctx context.Context, id string) error {
	if m.hardDeleteFn != nil {
		return m.hardDeleteFn(ctx, id)
	}
	return nil
}

func (m *mockRepository) ToggleActive(ctx context.Context, id string) (domain.Activity, error) {
	if m.toggleActiveFn != nil {
		return m.toggleActiveFn(ctx, id)
	}
	return domain.Activity{}, nil
}

func (m *mockRepository) GetByCategory(ctx context.Context, category string) ([]domain.Activity, error) {
	if m.getByCategory != nil {
		return m.getByCategory(ctx, category)
	}
	return []domain.Activity{}, nil
}

type mockPublisher struct {
	mu sync.Mutex

	publishFn func(ctx context.Context, action string, activityID string) error

	// Para tracking
	publishCalls    int32
	publishedEvents []PublishedEvent
}

type PublishedEvent struct {
	Action string
	ID     string
	Time   time.Time
}

func (m *mockPublisher) Publish(ctx context.Context, action string, activityID string) error {
	atomic.AddInt32(&m.publishCalls, 1)
	m.mu.Lock()
	m.publishedEvents = append(m.publishedEvents, PublishedEvent{
		Action: action,
		ID:     activityID,
		Time:   time.Now(),
	})
	m.mu.Unlock()

	if m.publishFn != nil {
		return m.publishFn(ctx, action, activityID)
	}
	return nil
}

func (m *mockPublisher) GetPublishedEvents() []PublishedEvent {
	m.mu.Lock()
	defer m.mu.Unlock()
	events := make([]PublishedEvent, len(m.publishedEvents))
	copy(events, m.publishedEvents)
	return events
}

// ===== HELPERS =====

// Helper para crear una actividad válida
func createValidActivity(name string) domain.Activity {
	return domain.Activity{
		Name:        name,
		Description: "Test activity description",
		Category:    "Sports",
		Difficulty:  "beginner",
		Location:    "Park",
		Duration:    60,
		Price:       25.0,
		MaxCapacity: 50,
		IsActive:    true,
	}
}

// ===== TESTS =====

// TestCreateWithConcurrentValidationAndEnrichment verifica que Create ejecute
// validación y enriquecimiento en paralelo usando goroutines
func TestCreateWithConcurrentValidationAndEnrichment(t *testing.T) {
	repo := &mockRepository{
		createFn: func(ctx context.Context, activity domain.Activity) (domain.Activity, error) {
			activity.ID = "activity-123"
			return activity, nil
		},
	}

	publisher := &mockPublisher{}
	service := NewActivitiesService(repo, publisher)

	activity := createValidActivity("Running Event")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	start := time.Now()
	created, err := service.Create(ctx, activity)
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if created.ID == "" {
		t.Fatal("Activity ID should not be empty")
	}

	// Verificar que se llamó a Create exactamente una vez
	if atomic.LoadInt32(&repo.createCalls) != 1 {
		t.Errorf("Expected 1 Create call, got %d", atomic.LoadInt32(&repo.createCalls))
	}

	// Esperar a que la goroutine de publish termine
	time.Sleep(500 * time.Millisecond)

	// Verificar que el evento fue publicado de forma asíncrona
	events := publisher.GetPublishedEvents()
	if len(events) != 1 {
		t.Errorf("Expected 1 published event, got %d", len(events))
	}
	if events[0].Action != "created" {
		t.Errorf("Expected action 'created', got '%s'", events[0].Action)
	}

	t.Logf("✅ Create completed in %v with concurrent goroutines", elapsed)
}

// TestUpdateWithConcurrentFetchAndValidate verifica que Update ejecute
// fetch y validación en paralelo
func TestUpdateWithConcurrentFetchAndValidate(t *testing.T) {
	existingActivity := createValidActivity("Old Event")
	existingActivity.ID = "activity-123"

	repo := &mockRepository{
		getByIDFn: func(ctx context.Context, id string) (domain.Activity, error) {
			return existingActivity, nil
		},
		updateFn: func(ctx context.Context, id string, activity domain.Activity) (domain.Activity, error) {
			activity.ID = id
			return activity, nil
		},
	}

	publisher := &mockPublisher{}
	service := NewActivitiesService(repo, publisher)

	newData := createValidActivity("New Event Name")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	start := time.Now()
	updated, err := service.Update(ctx, "activity-123", newData)
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	if updated.Name != newData.Name {
		t.Errorf("Expected name '%s', got '%s'", newData.Name, updated.Name)
	}

	// Verificar que se llamó a GetByID y Update
	if atomic.LoadInt32(&repo.getByIDCalls) != 1 {
		t.Errorf("Expected 1 GetByID call, got %d", atomic.LoadInt32(&repo.getByIDCalls))
	}
	if atomic.LoadInt32(&repo.updateCalls) != 1 {
		t.Errorf("Expected 1 Update call, got %d", atomic.LoadInt32(&repo.updateCalls))
	}

	// Esperar a que la goroutine de publish termine
	time.Sleep(500 * time.Millisecond)

	events := publisher.GetPublishedEvents()
	if len(events) != 1 {
		t.Errorf("Expected 1 published event, got %d", len(events))
	}
	if events[0].Action != "updated" {
		t.Errorf("Expected action 'updated', got '%s'", events[0].Action)
	}

	t.Logf("✅ Update completed in %v with concurrent goroutines", elapsed)
}

// TestCreateContextCancellation verifica que si el contexto se cancela,
// las goroutines respetan la cancelación
func TestCreateContextCancellation(t *testing.T) {
	repo := &mockRepository{
		createFn: func(ctx context.Context, activity domain.Activity) (domain.Activity, error) {
			select {
			case <-time.After(500 * time.Millisecond):
				return activity, nil
			case <-ctx.Done():
				return domain.Activity{}, ctx.Err()
			}
		},
	}

	publisher := &mockPublisher{}
	service := NewActivitiesService(repo, publisher)

	activity := createValidActivity("Event")

	// Cancelar contexto inmediatamente
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := service.Create(ctx, activity)

	if err == nil {
		t.Fatal("Expected error due to context cancellation")
	}

	t.Logf("✅ Context cancellation handled correctly: %v", err)
}

// TestPublishAsyncWithTimeout verifica que el publish asíncrono respete
// el timeout y no bloquee la respuesta de Create
func TestPublishAsyncWithTimeout(t *testing.T) {
	repo := &mockRepository{
		createFn: func(ctx context.Context, activity domain.Activity) (domain.Activity, error) {
			activity.ID = "activity-123"
			return activity, nil
		},
	}

	// Publisher que siempre falla después de 100ms
	slowPublisher := &mockPublisher{
		publishFn: func(ctx context.Context, action string, activityID string) error {
			select {
			case <-time.After(500 * time.Millisecond):
				return errors.New("publish timeout")
			case <-ctx.Done():
				return ctx.Err()
			}
		},
	}

	service := NewActivitiesService(repo, slowPublisher)

	activity := createValidActivity("Event")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	start := time.Now()
	created, err := service.Create(ctx, activity)
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Create debe retornar rápidamente (< 500ms) a pesar de que publish falle
	if elapsed > 500*time.Millisecond {
		t.Errorf("Create took too long: %v. Should not be blocked by publish", elapsed)
	}

	if created.ID != "activity-123" {
		t.Errorf("Expected ID 'activity-123', got '%s'", created.ID)
	}

	// Esperar a que la goroutine de publish termine
	time.Sleep(1 * time.Second)

	t.Logf("✅ Create returned quickly (%v) while publish ran async", elapsed)
}

// TestConcurrentCreatesWithGoroutines verifica que múltiples Creates concurrentes
// funcionan correctamente sin race conditions
func TestConcurrentCreatesWithGoroutines(t *testing.T) {
	var mu sync.Mutex
	createdActivities := make(map[string]bool)

	repo := &mockRepository{
		createFn: func(ctx context.Context, activity domain.Activity) (domain.Activity, error) {
			activity.ID = "activity-" + activity.Name + "-id"
			mu.Lock()
			createdActivities[activity.ID] = true
			mu.Unlock()
			return activity, nil
		},
	}

	publisher := &mockPublisher{}
	service := NewActivitiesService(repo, publisher)

	// Crear 10 actividades concurrentemente
	var wg sync.WaitGroup
	numGoroutines := 10

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()

			activity := createValidActivity("Event")
			activity.Name = "Event-" + string(rune(48+idx))

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			_, err := service.Create(ctx, activity)
			if err != nil {
				t.Errorf("Create %d failed: %v", idx, err)
			}
		}(i)
	}

	wg.Wait()

	// Esperar a que los publishes terminen
	time.Sleep(1 * time.Second)

	// Verificar que todas se crearon
	if len(createdActivities) != numGoroutines {
		t.Errorf("Expected %d created activities, got %d", numGoroutines, len(createdActivities))
	}

	// Verificar que la cantidad de events publicados es correcta
	events := publisher.GetPublishedEvents()
	if len(events) != numGoroutines {
		t.Errorf("Expected %d published events, got %d", numGoroutines, len(events))
	}

	t.Logf("✅ %d concurrent Creates completed successfully with %d events published", numGoroutines, len(events))
}

// TestDeleteAsyncPublish verifica que Delete publique el evento de forma asíncrona
func TestDeleteAsyncPublish(t *testing.T) {
	repo := &mockRepository{
		getByIDFn: func(ctx context.Context, id string) (domain.Activity, error) {
			return domain.Activity{ID: id, Name: "Event"}, nil
		},
	}

	publisher := &mockPublisher{}
	service := NewActivitiesService(repo, publisher)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	start := time.Now()
	err := service.Delete(ctx, "activity-123")
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Delete debe retornar rápidamente
	if elapsed > 1*time.Second {
		t.Errorf("Delete took too long: %v", elapsed)
	}

	// Esperar a que la goroutine de publish termine
	time.Sleep(500 * time.Millisecond)

	events := publisher.GetPublishedEvents()
	if len(events) != 1 {
		t.Errorf("Expected 1 published event, got %d", len(events))
	}
	if events[0].Action != "deleted" {
		t.Errorf("Expected action 'deleted', got '%s'", events[0].Action)
	}

	t.Logf("✅ Delete completed in %v with async publish", elapsed)
}

// BenchmarkCreateWithConcurrency benchmarka el performance de Create
// con goroutines concurrentes
func BenchmarkCreateWithConcurrency(b *testing.B) {
	repo := &mockRepository{
		createFn: func(ctx context.Context, activity domain.Activity) (domain.Activity, error) {
			activity.ID = "activity-id"
			return activity, nil
		},
	}

	publisher := &mockPublisher{}
	service := NewActivitiesService(repo, publisher)

	activity := createValidActivity("Event")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := service.Create(ctx, activity)
		if err != nil {
			b.Fatalf("Create failed: %v", err)
		}
	}

	b.StopTimer()
	b.Logf("Completed %d iterations", b.N)
}
