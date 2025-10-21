package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Reservation is the domain model
type Reservation struct {
	ID        primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	OwnerID   string                 `bson:"ownerId" json:"ownerId"`
	Resource  string                 `bson:"resource" json:"resource"`
	Start     time.Time              `bson:"start" json:"start"`
	End       time.Time              `bson:"end" json:"end"`
	Meta      map[string]interface{} `bson:"meta,omitempty" json:"meta,omitempty"`
	Score     float64                `bson:"score" json:"score"`
	CreatedAt time.Time              `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time              `bson:"updatedAt" json:"updatedAt"`
}
