package models

import (
	"time"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strings"
)

// SocialLinks estructura para enlaces de redes sociales
type SocialLinks struct {
	Instagram string `json:"instagram,omitempty"`
	Twitter   string `json:"twitter,omitempty"`
	Facebook  string `json:"facebook,omitempty"`
	LinkedIn  string `json:"linkedin,omitempty"`
	YouTube   string `json:"youtube,omitempty"`
	Website   string `json:"website,omitempty"`
}

// Value implementa driver.Valuer para GORM
func (s SocialLinks) Value() (driver.Value, error) {
	if s == (SocialLinks{}) {
		return nil, nil
	}
	return json.Marshal(s)
}

// Scan implementa sql.Scanner para GORM
func (s *SocialLinks) Scan(value interface{}) error {
	if value == nil {
		*s = SocialLinks{}
		return nil
	}
	
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	
	return json.Unmarshal(bytes, s)
}

// DateOnly estructura para manejar fechas sin tiempo
type DateOnly struct {
	time.Time
}

// UnmarshalJSON implementa json.Unmarshaler para manejar fechas en formato "YYYY-MM-DD"
func (d *DateOnly) UnmarshalJSON(data []byte) error {
	// Eliminar comillas del JSON
	str := strings.Trim(string(data), `"`)
	
	// Si es string vacío o null, retornar sin error
	if str == "" || str == "null" {
		return nil
	}
	
	// Parsear la fecha en formato YYYY-MM-DD
	parsed, err := time.Parse("2006-01-02", str)
	if err != nil {
		return err
	}
	
	d.Time = parsed
	return nil
}

// MarshalJSON implementa json.Marshaler para devolver la fecha en formato "YYYY-MM-DD"
func (d DateOnly) MarshalJSON() ([]byte, error) {
	if d.Time.IsZero() {
		return []byte("null"), nil
	}
	return []byte(`"` + d.Time.Format("2006-01-02") + `"`), nil
}

// Value implementa driver.Valuer para GORM
func (d DateOnly) Value() (driver.Value, error) {
	if d.Time.IsZero() {
		return nil, nil
	}
	return d.Time.Format("2006-01-02"), nil
}

// Scan implementa sql.Scanner para GORM
func (d *DateOnly) Scan(value interface{}) error {
	if value == nil {
		d.Time = time.Time{}
		return nil
	}
	
	switch v := value.(type) {
	case time.Time:
		d.Time = v
	case string:
		parsed, err := time.Parse("2006-01-02", v)
		if err != nil {
			return err
		}
		d.Time = parsed
	case []byte:
		parsed, err := time.Parse("2006-01-02", string(v))
		if err != nil {
			return err
		}
		d.Time = parsed
	default:
		return errors.New("cannot scan into DateOnly")
	}
	
	return nil
}

type User struct {
	ID               uint         `json:"id" gorm:"primaryKey"`
	Email            string       `json:"email" gorm:"type:varchar(255);uniqueIndex;not null"`
	Username         string       `json:"username" gorm:"type:varchar(100);uniqueIndex;not null"`
	Password         string       `json:"-" gorm:"not null"` // El "-" evita que se serialice en JSON
	FirstName        string       `json:"first_name" gorm:"type:varchar(100)"`
	LastName         string       `json:"last_name" gorm:"type:varchar(100)"`
	
	// Campos de perfil extendido
	AvatarURL        *string      `json:"avatar_url" gorm:"type:varchar(500)"`
	Bio              *string      `json:"bio" gorm:"type:text"`
	Phone            *string      `json:"phone" gorm:"type:varchar(20)"`
	BirthDate        *DateOnly    `json:"birth_date" gorm:"type:date"`
	Location         *string      `json:"location" gorm:"type:varchar(200)"`
	Gender           *string      `json:"gender" gorm:"type:enum('male','female','other','prefer_not_to_say');default:null"`
	
	// Campos específicos para plataforma deportiva
	Height           *float32     `json:"height" gorm:"type:decimal(5,2);comment:Height in cm"`
	Weight           *float32     `json:"weight" gorm:"type:decimal(5,2);comment:Weight in kg"`
	SportsInterests  *string      `json:"sports_interests" gorm:"type:json;comment:Array of sports as JSON"`
	FitnessLevel     *string      `json:"fitness_level" gorm:"type:enum('beginner','intermediate','advanced','professional');default:null"`
	
	// Enlaces sociales
	SocialLinks      SocialLinks  `json:"social_links" gorm:"type:json"`
	
	// Campos de sistema
	Role             string       `json:"role" gorm:"type:varchar(50);default:'user'"`
	EmailVerified    bool         `json:"email_verified" gorm:"default:false"`
	EmailVerifiedAt  *time.Time   `json:"email_verified_at"`
	IsActive         bool         `json:"is_active" gorm:"default:true"`
	LastLoginAt      *time.Time   `json:"last_login_at"`
	
	// Timestamps
	CreatedAt        time.Time    `json:"created_at"`
	UpdatedAt        time.Time    `json:"updated_at"`
}

// Response para no exponer datos sensibles
type UserResponse struct {
	ID               uint         `json:"id"`
	Email            string       `json:"email"`
	Username         string       `json:"username"`
	FirstName        string       `json:"first_name"`
	LastName         string       `json:"last_name"`
	AvatarURL        *string      `json:"avatar_url"`
	Bio              *string      `json:"bio"`
	Phone            *string      `json:"phone"`
	BirthDate        *DateOnly    `json:"birth_date"`
	Location         *string      `json:"location"`
	Gender           *string      `json:"gender"`
	Height           *float32     `json:"height"`
	Weight           *float32     `json:"weight"`
	SportsInterests  *string      `json:"sports_interests"`
	FitnessLevel     *string      `json:"fitness_level"`
	SocialLinks      SocialLinks  `json:"social_links"`
	Role             string       `json:"role"`
	EmailVerified    bool         `json:"email_verified"`
	EmailVerifiedAt  *time.Time   `json:"email_verified_at"`
	IsActive         bool         `json:"is_active"`
	LastLoginAt      *time.Time   `json:"last_login_at"`
	CreatedAt        time.Time    `json:"created_at"`
	UpdatedAt        time.Time    `json:"updated_at"`
}

// PublicUserResponse para listados públicos (menos información)
type PublicUserResponse struct {
	ID           uint        `json:"id"`
	Username     string      `json:"username"`
	Email        string      `json:"email"`
	FirstName    string      `json:"first_name"`
	LastName     string      `json:"last_name"`
	AvatarURL    *string     `json:"avatar_url"`
	Bio          *string     `json:"bio"`
	Location     *string     `json:"location"`
	SocialLinks  SocialLinks `json:"social_links"`
	FitnessLevel *string     `json:"fitness_level"`
	Role         string      `json:"role"`
	CreatedAt    time.Time   `json:"created_at"`
}

// UpdateProfileRequest estructura para actualizar perfil
type UpdateProfileRequest struct {
	FirstName       *string     `json:"first_name" validate:"omitempty,min=2,max=100"`
	LastName        *string     `json:"last_name" validate:"omitempty,min=2,max=100"`
	AvatarURL       *string     `json:"avatar_url" validate:"omitempty,url"`
	Bio             *string     `json:"bio" validate:"omitempty,max=500"`
	Phone           *string     `json:"phone" validate:"omitempty,e164"`
	BirthDate       *DateOnly   `json:"birth_date"`
	Location        *string     `json:"location" validate:"omitempty,max=200"`
	Gender          *string     `json:"gender" validate:"omitempty,oneof=male female other prefer_not_to_say"`
	Height          *float32    `json:"height" validate:"omitempty,min=50,max=300"`
	Weight          *float32    `json:"weight" validate:"omitempty,min=20,max=500"`
	SportsInterests *string     `json:"sports_interests"`
	FitnessLevel    *string     `json:"fitness_level" validate:"omitempty,oneof=beginner intermediate advanced professional"`
	SocialLinks     SocialLinks `json:"social_links"`
}

// ToResponse convierte User a UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:               u.ID,
		Email:            u.Email,
		Username:         u.Username,
		FirstName:        u.FirstName,
		LastName:         u.LastName,
		AvatarURL:        u.AvatarURL,
		Bio:              u.Bio,
		Phone:            u.Phone,
		BirthDate:        u.BirthDate,
		Location:         u.Location,
		Gender:           u.Gender,
		Height:           u.Height,
		Weight:           u.Weight,
		SportsInterests:  u.SportsInterests,
		FitnessLevel:     u.FitnessLevel,
		SocialLinks:      u.SocialLinks,
		Role:             u.Role,
		EmailVerified:    u.EmailVerified,
		EmailVerifiedAt:  u.EmailVerifiedAt,
		IsActive:         u.IsActive,
		LastLoginAt:      u.LastLoginAt,
		CreatedAt:        u.CreatedAt,
		UpdatedAt:        u.UpdatedAt,
	}
}

// ToPublicResponse convierte User a PublicUserResponse
func (u *User) ToPublicResponse() PublicUserResponse {
	return PublicUserResponse{
		ID:           u.ID,
		Username:     u.Username,
		Email:        u.Email,
		FirstName:    u.FirstName,
		LastName:     u.LastName,
		AvatarURL:    u.AvatarURL,
		Bio:          u.Bio,
		Location:     u.Location,
		SocialLinks:  u.SocialLinks,
		FitnessLevel: u.FitnessLevel,
		Role:         u.Role,
		CreatedAt:    u.CreatedAt,
	}
}
