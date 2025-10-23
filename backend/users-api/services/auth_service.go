package services

import (
	"context"
	"errors"
	"strings"

	"users-api/models"
	"users-api/repositories"
	"users-api/utils"

	"gorm.io/gorm"
)

// AuthResult represents the payload returned after successful authentication.
type AuthResult struct {
	Token string
	User  models.UserResponse
}

// RegisterInput contains the user supplied details for registration.
type RegisterInput struct {
	Email     string
	Username  string
	Password  string
	FirstName string
	LastName  string
}

// LoginInput captures login credentials.
type LoginInput struct {
	Email    string
	Password string
}

// AuthService encapsulates authentication use cases.
type AuthService struct {
	repo repositories.UserRepository
}

// NewAuthService builds an AuthService instance.
func NewAuthService(repo repositories.UserRepository) *AuthService {
	return &AuthService{repo: repo}
}

// Register creates a new user, hashing the password and returning an auth token.
func (s *AuthService) Register(ctx context.Context, input RegisterInput) (*AuthResult, error) {
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
		IsActive:  true,
	}

	if err := s.repo.Create(ctx, &user); err != nil {
		return nil, err
	}

	token, err := utils.GenerateJWT(user.ID, user.Email, user.Username)
	if err != nil {
		return nil, err
	}

	return &AuthResult{
		Token: token,
		User:  user.ToResponse(),
	}, nil
}

// Login validates credentials and returns a fresh auth token.
func (s *AuthService) Login(ctx context.Context, input LoginInput) (*AuthResult, error) {
	email := strings.ToLower(strings.TrimSpace(input.Email))

	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if !user.IsActive {
		return nil, ErrAccountDisabled
	}

	if !utils.CheckPassword(input.Password, user.Password) {
		return nil, ErrInvalidCredentials
	}

	token, err := utils.GenerateJWT(user.ID, user.Email, user.Username)
	if err != nil {
		return nil, err
	}

	return &AuthResult{
		Token: token,
		User:  user.ToResponse(),
	}, nil
}

// Refresh validates an existing token and issues a new one.
func (s *AuthService) Refresh(_ context.Context, token string) (string, error) {
	return utils.RefreshJWT(token)
}
