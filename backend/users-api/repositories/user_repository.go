package repositories

import (
	"context"
	"strings"

	"arq-soft-II/backend/users-api/models"

	"gorm.io/gorm"
)

// UserFilter holds optional filters for admin user listings.
type UserFilter struct {
	Role   string
	Status string
	Search string
	Page   int
	Limit  int
}

// UserStats represents aggregated counters for different user segments.
type UserStats struct {
	TotalUsers     int64
	ActiveUsers    int64
	InactiveUsers  int64
	RootUsers      int64
	AdminUsers     int64
	ModeratorUsers int64
	RegularUsers   int64
}

// UserRepository exposes persistence operations for users.
type UserRepository interface {
	GetByID(ctx context.Context, id uint) (*models.User, error)
	GetActiveByID(ctx context.Context, id uint) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	Create(ctx context.Context, user *models.User) error
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, user *models.User) error
	ListActive(ctx context.Context, page, limit int) ([]models.User, int64, error)
	List(ctx context.Context, filter UserFilter) ([]models.User, int64, error)
	GetStats(ctx context.Context) (UserStats, error)
	ExistsRoot(ctx context.Context) (bool, error)
}

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository instantiates a UserRepository backed by GORM.
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetActiveByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).
		Where("email = ?", strings.ToLower(email)).
		First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).
		Where("username = ?", username).
		First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *userRepository) Delete(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Delete(user).Error
}

func (r *userRepository) ListActive(ctx context.Context, page, limit int) ([]models.User, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	var (
		users []models.User
		total int64
	)

	query := r.db.WithContext(ctx).Model(&models.User{}).Where("is_active = ?", true)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *userRepository) List(ctx context.Context, filter UserFilter) ([]models.User, int64, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}

	offset := (filter.Page - 1) * filter.Limit

	var (
		users []models.User
		total int64
	)

	query := r.db.WithContext(ctx).Model(&models.User{})

	if filter.Role != "" {
		query = query.Where("role = ?", filter.Role)
	}

	switch filter.Status {
	case "active":
		query = query.Where("is_active = ?", true)
	case "inactive":
		query = query.Where("is_active = ?", false)
	}

	if filter.Search != "" {
		search := "%" + strings.ToLower(filter.Search) + "%"
		query = query.Where(
			"LOWER(email) LIKE ? OR LOWER(username) LIKE ? OR LOWER(first_name) LIKE ? OR LOWER(last_name) LIKE ?",
			search, search, search, search,
		)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Offset(offset).
		Limit(filter.Limit).
		Order("created_at DESC").
		Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *userRepository) GetStats(ctx context.Context) (UserStats, error) {
	var stats UserStats

	if err := r.db.WithContext(ctx).Model(&models.User{}).Count(&stats.TotalUsers).Error; err != nil {
		return stats, err
	}
	if err := r.db.WithContext(ctx).Model(&models.User{}).Where("is_active = ?", true).Count(&stats.ActiveUsers).Error; err != nil {
		return stats, err
	}
	if err := r.db.WithContext(ctx).Model(&models.User{}).Where("is_active = ?", false).Count(&stats.InactiveUsers).Error; err != nil {
		return stats, err
	}
	if err := r.db.WithContext(ctx).Model(&models.User{}).Where("role = ?", "root").Count(&stats.RootUsers).Error; err != nil {
		return stats, err
	}
	if err := r.db.WithContext(ctx).Model(&models.User{}).Where("role = ?", "admin").Count(&stats.AdminUsers).Error; err != nil {
		return stats, err
	}
	if err := r.db.WithContext(ctx).Model(&models.User{}).Where("role = ?", "moderator").Count(&stats.ModeratorUsers).Error; err != nil {
		return stats, err
	}
	if err := r.db.WithContext(ctx).Model(&models.User{}).Where("role = ?", "user").Count(&stats.RegularUsers).Error; err != nil {
		return stats, err
	}

	return stats, nil
}

func (r *userRepository) ExistsRoot(ctx context.Context) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("role = ?", "root").
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
