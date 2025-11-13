package domain

import (
	"time"
)

type Activity struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Category    string    `json:"category"`    // deportes, fitness, yoga, natación, etc.
	Difficulty  string    `json:"difficulty"`  // beginner, intermediate, advanced
	Location    string    `json:"location"`    // Ubicación/sede
	Price       float64   `json:"price"`       // Precio de la actividad
	Duration    int       `json:"duration"`    // Duración en minutos
	MaxCapacity int       `json:"max_capacity"` // Cupo máximo
	Instructor  string    `json:"instructor"`  // Nombre del instructor
	Schedule    []string  `json:"schedule"`    // Horarios disponibles (ej: ["Lunes 18:00", "Miércoles 18:00"])
	Equipment   []string  `json:"equipment"`   // Equipamiento necesario
	ImageURL    string    `json:"image_url"`   // URL de la imagen
	IsActive    bool      `json:"is_active"`   // Actividad activa/inactiva
	CreatedBy   uint      `json:"created_by"`  // ID del usuario admin que creó la actividad
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
