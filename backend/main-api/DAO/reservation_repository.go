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

func (m *mongoReservationRepo) UpdateByID(ctx context.Context, id string, r *models.Reservation) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	r.UpdatedAt = time.Now().UTC()
	_, err = m.col.UpdateByID(ctx, oid, bson.M{"$set": r})
	return err
}

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
