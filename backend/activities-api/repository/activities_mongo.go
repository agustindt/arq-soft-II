package repository

import (
	"activities-api/dao"
	"activities-api/domain"
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoActivitiesRepository implementa ActivitiesRepository usando MongoDB
type MongoActivitiesRepository struct {
	col *mongo.Collection
}

// NewMongoActivitiesRepository crea una nueva instancia del repository
func NewMongoActivitiesRepository(ctx context.Context, uri, dbName, collectionName string) *MongoActivitiesRepository {
	opt := options.Client().ApplyURI(uri)
	opt.SetServerSelectionTimeout(10 * time.Second)

	client, err := mongo.Connect(ctx, opt)
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
		return nil
	}

	pingCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := client.Ping(pingCtx, nil); err != nil {
		log.Fatalf("Error pinging MongoDB: %v", err)
		return nil
	}

	log.Println("✁EConnected to MongoDB successfully")

	return &MongoActivitiesRepository{
		col: client.Database(dbName).Collection(collectionName),
	}
}

// List obtiene todas las actividades
func (r *MongoActivitiesRepository) List(ctx context.Context) ([]domain.Activity, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Solo actividades activas por defecto
	filter := bson.M{"is_active": true}

	cur, err := r.col.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var daoActivities []dao.Activity
	if err := cur.All(ctx, &daoActivities); err != nil {
		return nil, err
	}

	domainActivities := make([]domain.Activity, len(daoActivities))
	for i, daoActivity := range daoActivities {
		domainActivities[i] = daoActivity.ToDomain()
	}

	return domainActivities, nil
}

// ListAll obtiene todas las actividades (incluyendo inactivas)
func (r *MongoActivitiesRepository) ListAll(ctx context.Context) ([]domain.Activity, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	cur, err := r.col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var daoActivities []dao.Activity
	if err := cur.All(ctx, &daoActivities); err != nil {
		return nil, err
	}

	domainActivities := make([]domain.Activity, len(daoActivities))
	for i, daoActivity := range daoActivities {
		domainActivities[i] = daoActivity.ToDomain()
	}

	return domainActivities, nil
}

// Create inserta una nueva actividad
func (r *MongoActivitiesRepository) Create(ctx context.Context, activity domain.Activity) (domain.Activity, error) {
	activityDAO := dao.FromDomain(activity)

	activityDAO.ID = primitive.NewObjectID()
	activityDAO.CreatedAt = time.Now().UTC()
	activityDAO.UpdatedAt = time.Now().UTC()
	activityDAO.IsActive = true // Por defecto activa

	_, err := r.col.InsertOne(ctx, activityDAO)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return domain.Activity{}, errors.New("activity with the same ID already exists")
		}
		return domain.Activity{}, err
	}

	return activityDAO.ToDomain(), nil
}

// GetByID busca una actividad por su ID
func (r *MongoActivitiesRepository) GetByID(ctx context.Context, id string) (domain.Activity, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.Activity{}, errors.New("invalid ID format")
	}

	var activityDAO dao.Activity
	err = r.col.FindOne(ctx, bson.M{"_id": objID}).Decode(&activityDAO)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.Activity{}, errors.New("activity not found")
		}
		return domain.Activity{}, err
	}

	return activityDAO.ToDomain(), nil
}

// Update actualiza una actividad existente
func (r *MongoActivitiesRepository) Update(ctx context.Context, id string, activity domain.Activity) (domain.Activity, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.Activity{}, errors.New("invalid ID format")
	}

	update := bson.M{
		"$set": bson.M{
			"name":         activity.Name,
			"description":  activity.Description,
			"category":     activity.Category,
			"difficulty":   activity.Difficulty,
			"location":     activity.Location,
			"price":        activity.Price,
			"duration":     activity.Duration,
			"max_capacity": activity.MaxCapacity,
			"instructor":   activity.Instructor,
			"schedule":     activity.Schedule,
			"equipment":    activity.Equipment,
			"image_url":    activity.ImageURL,
			"updated_at":   time.Now().UTC(),
		},
	}

	result, err := r.col.UpdateByID(ctx, objID, update)
	if err != nil {
		return domain.Activity{}, err
	}
	if result.MatchedCount == 0 {
		return domain.Activity{}, errors.New("activity not found")
	}

	return r.GetByID(ctx, id)
}

// Delete elimina una actividad (soft delete - marca como inactiva)
func (r *MongoActivitiesRepository) Delete(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid ID format")
	}

	// Soft delete - solo marcamos como inactiva
	update := bson.M{
		"$set": bson.M{
			"is_active":  false,
			"updated_at": time.Now().UTC(),
		},
	}

	result, err := r.col.UpdateByID(ctx, objID, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("activity not found")
	}

	return nil
}

// HardDelete elimina una actividad permanentemente de la base de datos
func (r *MongoActivitiesRepository) HardDelete(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid ID format")
	}

	result, err := r.col.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("activity not found")
	}

	return nil
}

// ToggleActive activa/desactiva una actividad
func (r *MongoActivitiesRepository) ToggleActive(ctx context.Context, id string) (domain.Activity, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.Activity{}, errors.New("invalid ID format")
	}

	// Primero obtener el estado actual
	var activityDAO dao.Activity
	err = r.col.FindOne(ctx, bson.M{"_id": objID}).Decode(&activityDAO)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.Activity{}, errors.New("activity not found")
		}
		return domain.Activity{}, err
	}

	// Toggle el estado
	update := bson.M{
		"$set": bson.M{
			"is_active":  !activityDAO.IsActive,
			"updated_at": time.Now().UTC(),
		},
	}

	_, err = r.col.UpdateByID(ctx, objID, update)
	if err != nil {
		return domain.Activity{}, err
	}

	return r.GetByID(ctx, id)
}

// GetByCategory obtiene actividades por categoría
func (r *MongoActivitiesRepository) GetByCategory(ctx context.Context, category string) ([]domain.Activity, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	filter := bson.M{
		"category":  category,
		"is_active": true,
	}

	cur, err := r.col.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var daoActivities []dao.Activity
	if err := cur.All(ctx, &daoActivities); err != nil {
		return nil, err
	}

	domainActivities := make([]domain.Activity, len(daoActivities))
	for i, daoActivity := range daoActivities {
		domainActivities[i] = daoActivity.ToDomain()
	}

	return domainActivities, nil
}
