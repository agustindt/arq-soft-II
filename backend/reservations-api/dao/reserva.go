package dao

import (
	"reservations-api/domain"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Reserva struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UsersID   []int              `json:"users_id" bson:"users_id"`
	Cupo      int                `json:"cupo" bson:"cupo"`
	Actividad string             `json:"actividad" bson:"actividad"`
	Date      time.Time          `json:"date" bson:"date"`
	Status    string             `json:"status" bson:"status"` //Pendiente, confirmada, cancelada
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

// ToDomain convierte de modelo DB a modelo de negocio
func (d Reserva) ToDomain() domain.Reserva {
	return domain.Reserva{
		ID:        d.ID.Hex(), // ObjectID -> string
		UsersID:   d.UsersID,
		Cupo:      d.Cupo,
		Actividad: d.Actividad,
		Date:      d.Date,
		Status:    d.Status,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
}

// FromDomain convierte de modelo de negocio a modelo DB
func FromDomain(domainItem domain.Reserva) Reserva {
	// Si el ID está vacío, DB generará uno automáticamente
	var objectID primitive.ObjectID
	if domainItem.ID != "" {
		objectID, _ = primitive.ObjectIDFromHex(domainItem.ID)
	}

	return Reserva{
		ID:        objectID,
		UsersID:   domainItem.UsersID,
		Cupo:      domainItem.Cupo,
		Actividad: domainItem.Actividad,
		Date:      domainItem.Date,
		Status:    domainItem.Status,
		CreatedAt: domainItem.CreatedAt,
		UpdatedAt: domainItem.UpdatedAt,
	}
}
