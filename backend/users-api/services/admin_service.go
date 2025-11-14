package services

import (
	"context"
	"errors"
	"strings"

	"arq-soft-II/backend/users-api/models"
	"arq-soft-II/backend/users-api/repositories"
	"arq-soft-II/backend/users-api/utils"

	"gorm.io/gorm"
)

// Admin pagination response for full user data.
type AdminUserList struct {
	Users []models.UserResponse
	Total int64
}

// AdminService encapsulates administrative operations.
type AdminService struct {
	repo repositories.UserRepository
}

// NewAdminService constructs an AdminService.
func NewAdminService(repo repositories.UserRepository) *AdminService {
	return &AdminService{repo: repo}
}

// RootExistsError carries information about an existing root account.
type RootExistsError struct {
	Existing models.User
}

// Error implements the error interface.
func (e *RootExistsError) Error() string {
	return ErrRootAlreadyExists.Error()
}

// CreateRootInput holds the data needed to bootstrap the root user.
type CreateRootInput struct {
	Email     string
	Username  string
	Password  string
	FirstName string
	LastName  string
	SecretKey string
}

// CreateUserInput defines admin-created user data.
type CreateUserInput struct {
	Email     string
	Username  string
	Password  string
	FirstName string
	LastName  string
	Role      string
}

// UpdateUserRoleInput carries data for role changes.
type UpdateUserRoleInput struct {
	UserID uint
	Role   string
}

// UpdateUserStatusInput carries activation toggles.
type UpdateUserStatusInput struct {
	UserID uint
	Active bool
}

const rootSecretKey = "SPORTS_PLATFORM_ROOT_2024"

// CreateRoot bootstraps the first root user if none exists.
func (s *AdminService) CreateRoot(ctx context.Context, input CreateRootInput) (*models.UserResponse, error) {
	if input.SecretKey != rootSecretKey {
		return nil, ErrInvalidSecretKey
	}

	exists, err := s.repo.ExistsRoot(ctx)
	if err != nil {
		return nil, err
	}
	if exists {
		existingRoots, _, listErr := s.repo.List(ctx, repositories.UserFilter{
			Role:  "root",
			Page:  1,
			Limit: 1,
		})
		if listErr != nil {
			return nil, listErr
		}
		if len(existingRoots) > 0 {
			return nil, &RootExistsError{Existing: existingRoots[0]}
		}
		return nil, ErrRootAlreadyExists
	}

	email := strings.ToLower(strings.TrimSpace(input.Email))
	username := strings.TrimSpace(input.Username)

	if _, err := s.repo.GetByEmail(ctx, email); err == nil {
		return nil, ErrEmailAlreadyExists
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if _, err := s.repo.GetByUsername(ctx, username); err == nil {
		return nil, ErrUsernameAlreadyExists
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	user := models.User{
		Email:         email,
		Username:      username,
		Password:      hashedPassword,
		FirstName:     input.FirstName,
		LastName:      input.LastName,
		Role:          "root",
		EmailVerified: true,
		IsActive:      true,
	}

	if err := s.repo.Create(ctx, &user); err != nil {
		return nil, err
	}

	response := user.ToResponse()
	return &response, nil
}

// CreateUser creates a user on behalf of an administrator.
func (s *AdminService) CreateUser(ctx context.Context, input CreateUserInput) (*models.UserResponse, error) {
	email := strings.ToLower(strings.TrimSpace(input.Email))
	username := strings.TrimSpace(input.Username)

	if _, err := s.repo.GetByEmail(ctx, email); err == nil {
		return nil, ErrEmailAlreadyExists
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if _, err := s.repo.GetByUsername(ctx, username); err == nil {
		return nil, ErrUsernameAlreadyExists
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	user := models.User{
		Email:     email,
		Username:  username,
		Password:  hashedPassword,
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Role:      input.Role,
		IsActive:  true,
	}

	if err := s.repo.Create(ctx, &user); err != nil {
		return nil, err
	}

	response := user.ToResponse()
	return &response, nil
}

// ListUsers returns a paginated list of users for admin views.
func (s *AdminService) ListUsers(ctx context.Context, filter repositories.UserFilter) (*AdminUserList, error) {
	users, total, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	responses := make([]models.UserResponse, 0, len(users))
	for _, user := range users {
		responses = append(responses, user.ToResponse())
	}

	return &AdminUserList{
		Users: responses,
		Total: total,
	}, nil
}

// UpdateUserRole sets a user's role, respecting root guardrails.
func (s *AdminService) UpdateUserRole(ctx context.Context, input UpdateUserRoleInput, actorRole string) (*models.UserResponse, error) {
	user, err := s.repo.GetByID(ctx, input.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	if user.Role == "root" && actorRole != "root" {
		return nil, ErrCannotModifyRoot
	}

	user.Role = input.Role

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	response := user.ToResponse()
	return &response, nil
}

// UpdateUserStatus toggles a user's active state.
func (s *AdminService) UpdateUserStatus(ctx context.Context, input UpdateUserStatusInput) (*models.UserResponse, error) {
	user, err := s.repo.GetByID(ctx, input.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	if user.Role == "root" && !input.Active {
		return nil, ErrCannotDeactivateRoot
	}

	user.IsActive = input.Active

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	response := user.ToResponse()
	return &response, nil
}

// DeleteUser removes a user from the system.
func (s *AdminService) DeleteUser(ctx context.Context, userID uint) (*models.User, error) {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	if user.Role == "root" {
		return nil, ErrCannotDeleteRoot
	}

	if err := s.repo.Delete(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// GetSystemStats returns aggregated statistics for users.
func (s *AdminService) GetSystemStats(ctx context.Context) (repositories.UserStats, error) {
	return s.repo.GetStats(ctx)
}
