package services

import (
	"activities-api/domain"
	"activities-api/utils"
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

	listCalls    int32
	createCalls  int32
	getByIDCalls int32
	updateCalls  int32
	deleteCalls  int32
}

type mockPublisher struct {
	mu sync.Mutex

	publishedEvents []PublishedEvent
	publishFn       func(ctx context.Context, action string, activityID string) error
}

type PublishedEvent struct {
	Action string
	ID     string
	Time   time.Time
}

func (m *mockPublisher) Publish(ctx context.Context, action string, activityID string) error {
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

func adminCtx(base context.Context) context.Context {
	ctx := context.WithValue(base, utils.ContextUserIDKey, uint(1))
	return context.WithValue(ctx, utils.ContextUserRoleKey, "admin")
}

// ===== TESTS =====

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

	baseCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ctx := adminCtx(baseCtx)

	start := time.Now()
	created, err := service.Create(ctx, activity)
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if created.ID == "" {
		t.Fatal("Activity ID should not be empty")
	}

	if atomic.LoadInt32(&repo.createCalls) != 1 {
		t.Errorf("Expected 1 Create call, got %d", atomic.LoadInt32(&repo.createCalls))
	}

	time.Sleep(500 * time.Millisecond)

	events := publisher.GetPublishedEvents()
	if len(events) != 1 {
		t.Errorf("Expected 1 published event, got %d", len(events))
	}
	if events[0].Action != "created" {
		t.Errorf("Expected action 'created', got '%s'", events[0].Action)
	}

	t.Logf("✅ Create completed in %v with concurrent goroutines", elapsed)
}

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

	baseCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ctx := adminCtx(baseCtx)

	start := time.Now()
	updated, err := service.Update(ctx, "activity-123", newData)
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	if updated.Name != newData.Name {
		t.Errorf("Expected name '%s', got '%s'", newData.Name, updated.Name)
	}

	if atomic.LoadInt32(&repo.getByIDCalls) != 1 {
		t.Errorf("Expected 1 GetByID call, got %d", atomic.LoadInt32(&repo.getByIDCalls))
	}
	if atomic.LoadInt32(&repo.updateCalls) != 1 {
		t.Errorf("Expected 1 Update call, got %d", atomic.LoadInt32(&repo.updateCalls))
	}

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

	ctx, cancel := context.WithCancel(context.Background())
	ctx = adminCtx(ctx)
	cancel()

	_, err := service.Create(ctx, activity)

	if err == nil {
		t.Fatal("Expected error due to context cancellation")
	}

	t.Logf("✅ Context cancellation handled correctly: %v", err)
}

func TestPublishAsyncWithTimeout(t *testing.T) {
	repo := &mockRepository{
		createFn: func(ctx context.Context, activity domain.Activity) (domain.Activity, error) {
			activity.ID = "activity-123"
			return activity, nil
		},
	}

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

	baseCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ctx := adminCtx(baseCtx)

	start := time.Now()
	created, err := service.Create(ctx, activity)
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if elapsed > 500*time.Millisecond {
		t.Errorf("Create took too long: %v. Should not be blocked by publish", elapsed)
	}

	if created.ID != "activity-123" {
		t.Errorf("Expected ID 'activity-123', got '%s'", created.ID)
	}

	time.Sleep(1 * time.Second)

	t.Logf("✅ Create returned quickly (%v) while publish ran async", elapsed)
}

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

	var wg sync.WaitGroup
	numGoroutines := 10

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()

			activity := createValidActivity("Event")
			activity.Name = "Event-" + string(rune(48+idx))

			baseCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			ctx := adminCtx(baseCtx)

			_, err := service.Create(ctx, activity)
			if err != nil {
				t.Errorf("Create %d failed: %v", idx, err)
			}
		}(i)
	}

	wg.Wait()

	time.Sleep(1 * time.Second)

	if len(createdActivities) != numGoroutines {
		t.Errorf("Expected %d created activities, got %d", numGoroutines, len(createdActivities))
	}

	events := publisher.GetPublishedEvents()
	if len(events) != numGoroutines {
		t.Errorf("Expected %d published events, got %d", numGoroutines, len(events))
	}

	t.Logf("✅ %d concurrent Creates completed successfully with %d events published", numGoroutines, len(events))
}

func TestDeleteAsyncPublish(t *testing.T) {
	repo := &mockRepository{
		getByIDFn: func(ctx context.Context, id string) (domain.Activity, error) {
			return domain.Activity{ID: id, Name: "Event"}, nil
		},
	}

	publisher := &mockPublisher{}
	service := NewActivitiesService(repo, publisher)

	baseCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ctx := adminCtx(baseCtx)

	start := time.Now()
	err := service.Delete(ctx, "activity-123")
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	if elapsed > 1*time.Second {
		t.Errorf("Delete took too long: %v", elapsed)
	}

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

	baseCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	ctx := adminCtx(baseCtx)

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
