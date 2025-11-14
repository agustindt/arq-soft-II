package repository

import (
	"arq-soft-II/backend/reservations-api/dao"
	"arq-soft-II/backend/reservations-api/domain"
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// MongoReservasRepository implementa ReservasRepository usando DB
type MongoReservasRepository struct {
	col *mongo.Collection // Referencia a la colección "Reservas" en DB
}

// NewMongoReservasRepository crea una nueva instancia del repository
// Recibe una referencia a la base de datos DB
func NewMongoReservasRepository(ctx context.Context, uri, dbName, collectionName string) *MongoReservasRepository {
	opt := options.Client().ApplyURI(uri)
	opt.SetServerSelectionTimeout(10 * time.Second)

	client, err := mongo.Connect(ctx, opt)
	if err != nil {
		log.Fatalf("Error connecting to DB: %v", err)
		return nil
	}

	pingCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := client.Ping(pingCtx, nil); err != nil {
		log.Fatalf("Error pinging DB: %v", err)
		return nil
	}

	return &MongoReservasRepository{
		col: client.Database(dbName).Collection(collectionName), // Conecta con la colección "Reservas"
	}
}

// List obtiene todos los Reservas de DB
func (r *MongoReservasRepository) List(ctx context.Context) ([]domain.Reserva, error) {
	// ⏰ Timeout para evitar que la operación se cuelgue
	// Esto es importante en producción para no bloquear indefinidamente
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// 🔍 Find() sin filtros retorna todos los documentos de la colección
	// bson.M{} es un filtro vacío (equivale a {} en DB shell)
	cur, err := r.col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx) // ⚠️ IMPORTANTE: Siempre cerrar el cursor para liberar recursos

	// 📦 Decodificar resultados en slice de DAO (modelo DB)
	// Usamos el modelo DAO porque maneja ObjectID y tags BSON
	var daoReservas []dao.Reserva
	if err := cur.All(ctx, &daoReservas); err != nil {
		return nil, err
	}

	// 🔄 Convertir de DAO a Domain (para la capa de negocio)
	// Separamos los modelos: DAO para DB, Domain para lógica de negocio
	domainReservas := make([]domain.Reserva, len(daoReservas))
	for i, daoReserva := range daoReservas {
		domainReservas[i] = daoReserva.ToDomain() // Función definida en dao/Reserva.go
	}

	return domainReservas, nil
}

// ListByUserID obtiene todas las reservas de un usuario específico
func (r *MongoReservasRepository) ListByUserID(ctx context.Context, userID int) ([]domain.Reserva, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Filtrar por user_id - buscar documentos donde users_id contenga el userID
	filter := bson.M{
		"users_id": bson.M{"$in": []int{userID}},
	}

	cur, err := r.col.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var daoReservas []dao.Reserva
	if err := cur.All(ctx, &daoReservas); err != nil {
		return nil, err
	}

	domainReservas := make([]domain.Reserva, len(daoReservas))
	for i, daoReserva := range daoReservas {
		domainReservas[i] = daoReserva.ToDomain()
	}

	return domainReservas, nil
}

// Create inserta un nuevo Reserva en DB
func (r *MongoReservasRepository) Create(ctx context.Context, Reserva domain.Reserva) (domain.Reserva, error) {
	ReservaDAO := dao.FromDomain(Reserva) // Convertir a modelo DAO

	ReservaDAO.ID = primitive.NewObjectID()
	ReservaDAO.CreatedAt = time.Now().UTC()
	ReservaDAO.UpdatedAt = time.Now().UTC()

	// Insertar en DB
	_, err := r.col.InsertOne(ctx, ReservaDAO)
	if err != nil {
		// Podemos manejar errores específicos de MongoDB, como claves duplicadas
		// https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo#IsDuplicateKeyError
		// Esto es útil si tenemos restricciones de unicidad en la colección
		if mongo.IsDuplicateKeyError(err) {
			return domain.Reserva{}, errors.New("reserva with the same ID already exists")
		}

		// Error genérico
		return domain.Reserva{}, err
	}

	return ReservaDAO.ToDomain(), nil // Convertir de vuelta a Domain
}

// GetByID busca un Reserva por su ID
// Consigna 2: Validar que el ID sea un ObjectID válido
func (r *MongoReservasRepository) GetByID(ctx context.Context, id string) (domain.Reserva, error) {
	// Validar que el ID es un ObjectID válido
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.Reserva{}, errors.New("invalid ID format")
	}

	// Buscar en DB
	var ReservaDAO dao.Reserva
	err = r.col.FindOne(ctx, bson.M{"_id": objID}).Decode(&ReservaDAO)
	if err != nil {
		// Manejar caso de no encontrado
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.Reserva{}, errors.New("reserva not found")
		}
		return domain.Reserva{}, err
	}

	return ReservaDAO.ToDomain(), nil
}

// Update actualiza un Reserva existente
// Consigna 3: Update parcial + actualizar updatedAt
func (r *MongoReservasRepository) Update(ctx context.Context, id string, Reserva domain.Reserva) (domain.Reserva, error) {
	// Validar que el ID es un ObjectID válido
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.Reserva{}, errors.New("invalid ID format")
	}

	// Preparar los campos a actualizar
	update := bson.M{
		"$set": bson.M{
			"actividad":  Reserva.Actividad,
			"users_id":   Reserva.UsersID,
			"cupo":       Reserva.Cupo,
			"date":       Reserva.Date,
			"status":     Reserva.Status,
			"updated_at": time.Now().UTC(), // Actualizar timestamp
		},
	}

	// Ejecutar la actualización
	result, err := r.col.UpdateByID(ctx, objID, update)
	if err != nil {
		return domain.Reserva{}, err
	}
	if result.MatchedCount == 0 {
		return domain.Reserva{}, errors.New("reserva not found")
	}

	// Retornar el Reserva actualizado
	return r.GetByID(ctx, id)
}

// Delete elimina un Reserva por ID
// Consigna 4: Eliminar documento de DB
func (r *MongoReservasRepository) Delete(ctx context.Context, id string) error {
	// Validar que el ID es un ObjectID válido
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid ID format")
	}

	// Ejecutar la eliminación
	result, err := r.col.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("reserva not found")
	}

	return nil
}
