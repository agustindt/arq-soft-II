package dao

import (
	"arq-soft-II/backend/activities-api/domain"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Activity struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	Category    string             `json:"category" bson:"category"`
	Difficulty  string             `json:"difficulty" bson:"difficulty"`
	Location    string             `json:"location" bson:"location"`
	Price       float64            `json:"price" bson:"price"`
	Duration    int                `json:"duration" bson:"duration"`
	MaxCapacity int                `json:"max_capacity" bson:"max_capacity"`
	Instructor  string             `json:"instructor" bson:"instructor"`
	Schedule    []string           `json:"schedule" bson:"schedule"`
	Equipment   []string           `json:"equipment" bson:"equipment"`
	ImageURL    string             `json:"image_url" bson:"image_url"`
	IsActive    bool               `json:"is_active" bson:"is_active"`
	CreatedBy   uint               `json:"created_by" bson:"created_by"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

// ToDomain convierte de modelo DB a modelo de negocio
func (d Activity) ToDomain() domain.Activity {
	return domain.Activity{
		ID:          d.ID.Hex(),
		Name:        d.Name,
		Description: d.Description,
		Category:    d.Category,
		Difficulty:  d.Difficulty,
		Location:    d.Location,
		Price:       d.Price,
		Duration:    d.Duration,
		MaxCapacity: d.MaxCapacity,
		Instructor:  d.Instructor,
		Schedule:    d.Schedule,
		Equipment:   d.Equipment,
		ImageURL:    d.ImageURL,
		IsActive:    d.IsActive,
		CreatedBy:   d.CreatedBy,
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
	}
}

// FromDomain convierte de modelo de negocio a modelo DB
func FromDomain(domainItem domain.Activity) Activity {
	// Si el ID está vacío, MongoDB generará uno automáticamente
	var objectID primitive.ObjectID
	if domainItem.ID != "" {
		objectID, _ = primitive.ObjectIDFromHex(domainItem.ID)
	}

	return Activity{
		ID:          objectID,
		Name:        domainItem.Name,
		Description: domainItem.Description,
		Category:    domainItem.Category,
		Difficulty:  domainItem.Difficulty,
		Location:    domainItem.Location,
		Price:       domainItem.Price,
		Duration:    domainItem.Duration,
		MaxCapacity: domainItem.MaxCapacity,
		Instructor:  domainItem.Instructor,
		Schedule:    domainItem.Schedule,
		Equipment:   domainItem.Equipment,
		ImageURL:    domainItem.ImageURL,
		IsActive:    domainItem.IsActive,
		CreatedBy:   domainItem.CreatedBy,
		CreatedAt:   domainItem.CreatedAt,
		UpdatedAt:   domainItem.UpdatedAt,
	}
}
