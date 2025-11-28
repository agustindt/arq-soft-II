package domain

import (
	"time"
)

type Reserva struct {
	ID        string    `json:"id"`
	UsersID   []int     `json:"users_id"`
	Cupo      int       `json:"cupo"`
	Actividad string    `json:"actividad"`
	Schedule  string    `json:"schedule"` // Horario específico (ej: "Lunes 20:00")
	Date      time.Time `json:"date"`
	Status    string    `json:"status"` //Pendiente, confirmada, cancelada
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
