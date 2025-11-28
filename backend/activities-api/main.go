// Package main implements the Activities API microservice.
//
// The Activities API manages sports activities with full CRUD operations and event-driven
// communication. It publishes events to RabbitMQ whenever activities are created, updated,
// or deleted, enabling other services (like Search API) to stay synchronized.
//
// Key Features:
//   - Activity CRUD operations (admin only for write operations)
//   - MongoDB document storage for flexible activity schema
//   - Event publishing to RabbitMQ for activity state changes
//   - Public endpoints for listing and viewing activities
//   - Category and status-based filtering
//   - JWT validation with admin role requirement for mutations
//   - Graceful shutdown with proper resource cleanup
//
// Event Types Published:
//   - activity.created: When a new activity is created
//   - activity.updated: When an activity is modified
//   - activity.deleted: When an activity is soft-deleted
//
// Database: MongoDB 6.0
// Message Queue: RabbitMQ (publisher)
// Port: 8082
//
// For complete API documentation, see docs/api/activities-api.md
package main

import (
	"context"
	"log"
	"time"

	"activities-api/clients"
	"activities-api/config"
	"activities-api/controllers"
	"activities-api/middleware"
	"activities-api/repository"
	"activities-api/services"

	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("🚀 Starting Activities API...")

	// Cargar configuración
	cfg := config.Load()

	// Context principal
	ctx := context.Background()

	// Repository: maneja operaciones en MongoDB
	activitiesRepo := repository.NewMongoActivitiesRepository(
		ctx,
		cfg.Mongo.URI,
		cfg.Mongo.DB,
		cfg.Mongo.Collection,
	)

	// RabbitMQ Client para publicar eventos
	rabbitClient := clients.NewRabbitMQClient(
		cfg.RabbitMQ.Username,
		cfg.RabbitMQ.Password,
		cfg.RabbitMQ.Host,
		cfg.RabbitMQ.Port,
		cfg.RabbitMQ.Exchange,
		activitiesRepo, // necesario para obtener datos completos al publicar
	)

	// Crear cola de publicación con workers y retries
	publishQueue := services.NewPublishQueue(rabbitClient, 200, 3, 200*time.Millisecond)
	publishQueue.Start(ctx, 2) // 2 workers

	// Service layer: lógica de negocio
	activitiesService := services.NewActivitiesService(activitiesRepo, publishQueue)

	// Controller: maneja requests HTTP
	activitiesController := controllers.NewActivitiesController(activitiesService)

	// Configurar router HTTP con Gin
	router := gin.Default()

	// Middleware global
	router.Use(middleware.CORSMiddleware)

	// Health check endpoint
	router.GET("/healthz", activitiesController.HealthCheck)

	usersAPI := cfg.UsersAPIURL

	// Endpoints públicos (sin autenticación)
	public := router.Group("/activities")
	{
		public.GET("/", activitiesController.GetActivities)
		public.GET("/all", activitiesController.GetAllActivities)
		public.GET("/category/:category", activitiesController.GetActivitiesByCategory)
		public.GET("/:id", activitiesController.GetActivityByID)
	}

	// Endpoints protegidos (requieren JWT Admin)
	admin := router.Group("/activities")
	admin.Use(middleware.AdminOnly(usersAPI))
	{
		admin.POST("/", activitiesController.CreateActivity)
		admin.PUT("/:id", activitiesController.UpdateActivity)
		admin.DELETE("/:id", activitiesController.DeleteActivity)
		admin.DELETE("/:id/hard", activitiesController.HardDelete)
		admin.PUT("/:id/toggle", activitiesController.ToggleActive)
	}

	// Iniciar servidor
	log.Println("🌐 Activities API listening on port 8082...")
	if err := router.Run(":8082"); err != nil {
		log.Fatalf("❌ Could not start server: %v", err)
	}
}
