package domain

import "time"

// Activity representa una actividad deportiva (modelo de dominio para Search API)
type Activity struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	Difficulty  string    `json:"difficulty"`
	Location    string    `json:"location"`
	Price       float64   `json:"price"`
	Duration    int       `json:"duration"`
	MaxCapacity int       `json:"max_capacity"`
	Instructor  string    `json:"instructor"`
	Schedule    []string  `json:"schedule"`
	Equipment   []string  `json:"equipment"`
	ImageURL    string    `json:"image_url"`
	IsActive    bool      `json:"is_active"`
	CreatedBy   uint      `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
