package services

import (
	"context"
	"errors"
	"testing"

	"users-api/models"
	"users-api/repositories"
	"users-api/utils"

	"gorm.io/gorm"
)

type mockUserRepository struct {
	getByIDFn     func(ctx context.Context, id uint) (*models.User, error)
	getActiveByID func(ctx context.Context, id uint) (*models.User, error)
	updateFn      func(ctx context.Context, user *models.User) error
}

func (m *mockUserRepository) GetByID(ctx context.Context, id uint) (*models.User, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockUserRepository) GetActiveByID(ctx context.Context, id uint) (*models.User, error) {
	if m.getActiveByID != nil {
		return m.getActiveByID(ctx, id)
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	return nil, gorm.ErrRecordNotFound
}

func (m *mockUserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	return nil, gorm.ErrRecordNotFound
}

func (m *mockUserRepository) Create(ctx context.Context, user *models.User) error {
	return nil
}

func (m *mockUserRepository) Update(ctx context.Context, user *models.User) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, user)
	}
	return nil
}

func (m *mockUserRepository) Delete(ctx context.Context, user *models.User) error {
	return nil
}

func (m *mockUserRepository) ListActive(ctx context.Context, page, limit int) ([]models.User, int64, error) {
	return []models.User{}, 0, nil
}

func (m *mockUserRepository) List(ctx context.Context, filter repositories.UserFilter) ([]models.User, int64, error) {
	return []models.User{}, 0, nil
}

func (m *mockUserRepository) GetStats(ctx context.Context) (repositories.UserStats, error) {
	return repositories.UserStats{}, nil
}

func (m *mockUserRepository) ExistsRoot(ctx context.Context) (bool, error) {
	return false, nil
}

func TestGetProfileReturnsUser(t *testing.T) {
	repo := &mockUserRepository{
		getByIDFn: func(ctx context.Context, id uint) (*models.User, error) {
			return &models.User{ID: id, Email: "test@example.com", Username: "tester"}, nil
		},
	}

	service := NewUserService(repo)

	profile, err := service.GetProfile(context.Background(), 1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if profile == nil || profile.ID != 1 {
		t.Fatalf("expected profile with ID 1, got %+v", profile)
	}
}

func TestChangePasswordInvalidCurrent(t *testing.T) {
	hashed, _ := utils.HashPassword("correct-password")

	repo := &mockUserRepository{
		getByIDFn: func(ctx context.Context, id uint) (*models.User, error) {
			return &models.User{ID: id, Password: hashed}, nil
		},
		updateFn: func(ctx context.Context, user *models.User) error {
			return errors.New("update should not be called on invalid password")
		},
	}

	service := NewUserService(repo)

	err := service.ChangePassword(context.Background(), 2, "wrong-password", "newpass")
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}
