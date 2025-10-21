package repository

import (
	"context"
	"time"

	errorspkg "reservations/errors"
	"reservations/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// --- INTERFAZ ---
// Esto es lo que el service espera usar.
type ReservationRepository interface {
	Insert(ctx context.Context, r *models.Reservation) (string, error)
	GetByID(ctx context.Context, id string) (*models.Reservation, error)
	UpdateByID(ctx context.Context, id string, r *models.Reservation) error
	DeleteByID(ctx context.Context, id string) error
}

// --- IMPLEMENTACIÓN MONGO ---

type mongoReservationRepo struct {
	col *mongo.Collection
}

// Constructor
func NewMongoReservationRepo(db *mongo.Database) ReservationRepository {
	return &mongoReservationRepo{
		col: db.Collection("reservations"),
	}
}

// Insert crea una nueva reserva
func (m *mongoReservationRepo) Insert(ctx context.Context, r *models.Reservation) (string, error) {
	r.ID = primitive.NewObjectID()
	r.CreatedAt = time.Now().UTC()
	r.UpdatedAt = r.CreatedAt
	_, err := m.col.InsertOne(ctx, r)
	if err != nil {
		return "", err
	}
	return r.ID.Hex(), nil
}

// GetByID obtiene una reserva por ID
func (m *mongoReservationRepo) GetByID(ctx context.Context, id string) (*models.Reservation, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var r models.Reservation
	if err := m.col.FindOne(ctx, bson.M{"_id": oid}).Decode(&r); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errorspkg.ErrNotFound
		}
		return nil, err
	}
	return &r, nil
}

// UpdateByID actualiza una reserva
func (m *mongoReservationRepo) UpdateByID(ctx context.Context, id string, r *models.Reservation) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	r.UpdatedAt = time.Now().UTC()
	_, err = m.col.UpdateByID(ctx, oid, bson.M{"$set": r})
	return err
}

// DeleteByID elimina una reserva
func (m *mongoReservationRepo) DeleteByID(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	res, err := m.col.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errorspkg.ErrNotFound
	}
	return nil
}
