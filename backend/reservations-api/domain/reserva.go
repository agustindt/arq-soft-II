package domain

import (
	"time"
)

type Reserva struct {
	ID        string    `json:"id"`
	UserID    int       `json:"user_id"`
	Actividad string    `json:"actividad"`
	Date      time.Time `json:"date"`
	Status    string    `json:"status"` //Pendiente, confirmada, cancelada
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
