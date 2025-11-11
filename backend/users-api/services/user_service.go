package services

import (
	"context"
	"errors"

	"arq-soft-II/backend/users-api/models"
	"arq-soft-II/backend/users-api/repositories"
	"arq-soft-II/backend/users-api/utils"

	"gorm.io/gorm"
)

// PaginatedUsers contains a slice of users with the total count for pagination.
type PaginatedUsers struct {
	Users []models.UserResponse
	Total int64
}

// PaginatedPublicUsers contains public user data with pagination metadata.
type PaginatedPublicUsers struct {
	Users []models.PublicUserResponse
	Total int64
}

// UserService orchestrates user related business logic.
type UserService struct {
	repo repositories.UserRepository
}

// NewUserService creates a UserService instance.
func NewUserService(repo repositories.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// GetProfile returns the full profile for a user.
func (s *UserService) GetProfile(ctx context.Context, userID uint) (*models.UserResponse, error) {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	response := user.ToResponse()
	return &response, nil
}

// UpdateProfile updates the mutable fields of a user profile.
func (s *UserService) UpdateProfile(ctx context.Context, userID uint, req *models.UpdateProfileRequest) (*models.UserResponse, error) {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	if req.FirstName != nil {
		user.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		user.LastName = *req.LastName
	}
	if req.AvatarURL != nil {
		user.AvatarURL = req.AvatarURL
	}
	if req.Bio != nil {
		user.Bio = req.Bio
	}
	if req.Phone != nil {
		user.Phone = req.Phone
	}
	if req.BirthDate != nil {
		user.BirthDate = req.BirthDate
	}
	if req.Location != nil {
		user.Location = req.Location
	}
	if req.Gender != nil {
		user.Gender = req.Gender
	}
	if req.Height != nil {
		user.Height = req.Height
	}
	if req.Weight != nil {
		user.Weight = req.Weight
	}
	if req.SportsInterests != nil {
		user.SportsInterests = req.SportsInterests
	}
	if req.FitnessLevel != nil {
		user.FitnessLevel = req.FitnessLevel
	}

	user.SocialLinks = req.SocialLinks

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	response := user.ToResponse()
	return &response, nil
}

// ChangePassword updates a user's password after verifying the current credentials.
func (s *UserService) ChangePassword(ctx context.Context, userID uint, currentPassword, newPassword string) error {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	if !utils.CheckPassword(currentPassword, user.Password) {
		return ErrInvalidCredentials
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	user.Password = hashedPassword

	return s.repo.Update(ctx, user)
}

// GetPublicUserByID returns the public profile for an active user.
func (s *UserService) GetPublicUserByID(ctx context.Context, userID uint) (*models.PublicUserResponse, error) {
	user, err := s.repo.GetActiveByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	response := user.ToPublicResponse()
	return &response, nil
}

// ListPublicUsers returns active users with pagination.
func (s *UserService) ListPublicUsers(ctx context.Context, page, limit int) (*PaginatedPublicUsers, error) {
	users, total, err := s.repo.ListActive(ctx, page, limit)
	if err != nil {
		return nil, err
	}

	responses := make([]models.PublicUserResponse, 0, len(users))
	for _, user := range users {
		responses = append(responses, user.ToPublicResponse())
	}

	return &PaginatedPublicUsers{
		Users: responses,
		Total: total,
	}, nil
}

// UpdateAvatar sets the avatar URL for a user and returns the previous path for cleanup.
func (s *UserService) UpdateAvatar(ctx context.Context, userID uint, avatarURL string) (*models.UserResponse, *string, error) {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, ErrUserNotFound
		}
		return nil, nil, err
	}

	var previous *string
	if user.AvatarURL != nil && *user.AvatarURL != "" {
		prev := *user.AvatarURL
		previous = &prev
	}

	user.AvatarURL = &avatarURL

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, previous, err
	}

	response := user.ToResponse()
	return &response, previous, nil
}

// DeleteAvatar removes the avatar URL association and returns the removed path for cleanup.
func (s *UserService) DeleteAvatar(ctx context.Context, userID uint) (*models.UserResponse, *string, error) {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, ErrUserNotFound
		}
		return nil, nil, err
	}

	var previous *string
	if user.AvatarURL != nil && *user.AvatarURL != "" {
		prev := *user.AvatarURL
		previous = &prev
	}

	user.AvatarURL = nil

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, previous, err
	}

	response := user.ToResponse()
	return &response, previous, nil
}
