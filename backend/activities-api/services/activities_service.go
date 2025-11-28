package services

import (
	"activities-api/clients"
	"activities-api/domain"
	"activities-api/utils"
	"context"
	"errors"
	"fmt"
	"os"
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
	usersAPI   string
}

var (
	ErrOwnerForbidden = errors.New("requester is not allowed to modify this resource")
	ErrOwnerNotFound  = errors.New("owner user not found")
)

// Constructor
func NewActivitiesService(repo ActivitiesRepository, publisher ActivityPublisher) *ActivitiesServiceImpl {
	return &ActivitiesServiceImpl{
		repository: repo,
		publisher:  publisher,
		usersAPI:   os.Getenv("USERS_API_URL"),
	}
}

func (s *ActivitiesServiceImpl) requesterFromContext(ctx context.Context) (uint, string) {
	id, role := utils.RequesterFromContext(ctx)
	return id, role
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

func (s *ActivitiesServiceImpl) validateActivity(activity domain.Activity) error {
	if strings.TrimSpace(activity.Name) == "" {
		return errors.New("activity name is required")
	}
	if strings.TrimSpace(activity.Category) == "" {
		return errors.New("activity category is required")
	}
	if strings.TrimSpace(activity.Difficulty) == "" {
		return errors.New("activity difficulty is required")
	}
	if activity.Duration <= 0 {
		return errors.New("activity duration must be greater than zero")
	}
	if activity.MaxCapacity <= 0 {
		return errors.New("max capacity must be greater than zero")
	}
	if activity.Price < 0 {
		return errors.New("price cannot be negative")
	}

	return nil
}

func (s *ActivitiesServiceImpl) enrichActivity(activity *domain.Activity) error {
	now := time.Now()

	if activity.CreatedAt.IsZero() {
		activity.CreatedAt = now
	}

	activity.UpdatedAt = now

	if activity.Schedule == nil {
		activity.Schedule = []string{}
	}
	if activity.Equipment == nil {
		activity.Equipment = []string{}
	}

	return nil
}

// List activities (active only)
func (s *ActivitiesServiceImpl) List(ctx context.Context) ([]domain.Activity, error) {
	return s.repository.List(ctx)
}

// List all activities (active + inactive)
func (s *ActivitiesServiceImpl) ListAll(ctx context.Context) ([]domain.Activity, error) {
	return s.repository.ListAll(ctx)
}

// Get by category
func (s *ActivitiesServiceImpl) GetByCategory(ctx context.Context, category string) ([]domain.Activity, error) {
	return s.repository.GetByCategory(ctx, category)
}

// Get by ID
func (s *ActivitiesServiceImpl) GetByID(ctx context.Context, id string) (domain.Activity, error) {
	return s.repository.GetByID(ctx, id)
}

// Create activity
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

	created, err := s.repository.Create(tasksCtx, activity)
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

// Update activity
func (s *ActivitiesServiceImpl) Update(ctx context.Context, id string, activity domain.Activity) (domain.Activity, error) {
	if err := s.validateActivity(activity); err != nil {
		return domain.Activity{}, err
	}

	if err := s.enrichActivity(&activity); err != nil {
		return domain.Activity{}, err
	}

	updated, err := s.repository.Update(ctx, id, activity)
	if err != nil {
		return domain.Activity{}, err
	}

	go func(id string) {
		pubCtx, pubCancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer pubCancel()
		s.publisher.Publish(pubCtx, "updated", id)
	}(updated.ID)

	return updated, nil
}

// Delete (soft delete)
func (s *ActivitiesServiceImpl) Delete(ctx context.Context, id string) error {
	err := s.repository.Delete(ctx, id)
	if err != nil {
		return err
	}

	go func() {
		pubCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		s.publisher.Publish(pubCtx, "deleted", id)
	}()

	return nil
}

// Hard delete
func (s *ActivitiesServiceImpl) HardDelete(ctx context.Context, id string) error {
	err := s.repository.HardDelete(ctx, id)
	if err != nil {
		return err
	}

	go func() {
		pubCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		s.publisher.Publish(pubCtx, "hard_deleted", id)
	}()

	return nil
}

// Toggle active
func (s *ActivitiesServiceImpl) ToggleActive(ctx context.Context, id string) (domain.Activity, error) {
	activity, err := s.repository.ToggleActive(ctx, id)
	if err != nil {
		return domain.Activity{}, err
	}

	go func() {
		pubCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		s.publisher.Publish(pubCtx, "toggled", id)
	}()

	return activity, nil
}
