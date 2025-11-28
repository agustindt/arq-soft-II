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
	"activities-api/clients"
	"activities-api/config"
	"activities-api/controllers"
	"activities-api/middleware"
	"activities-api/repository"
	"activities-api/services"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	)xq

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
		public.GET("", activitiesController.GetActivities)                              // Listar activas
		public.GET("/:id", activitiesController.GetActivityByID)                        // Obtener por ID
		public.GET("/category/:category", activitiesController.GetActivitiesByCategory) // Filtrar por categoría
	}

	// Endpoints protegidos (requieren admin role)
	admin := router.Group("/activities")
	admin.Use(middleware.AdminOnly(usersAPI))
	{
		admin.GET("/all", activitiesController.GetAllActivities)              // Listar todas (incluyendo inactivas)
		admin.POST("", activitiesController.CreateActivity)                   // Crear actividad
		admin.PUT("/:id", activitiesController.UpdateActivity)                // Actualizar actividad
		admin.DELETE("/:id", activitiesController.DeleteActivity)             // Eliminar actividad (soft delete)
		admin.PATCH("/:id/toggle", activitiesController.ToggleActiveActivity) // Activar/desactivar
	}

	// Configuración del server HTTP
	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Printf("✁EActivities API listening on port %s", cfg.Port)
	log.Printf("📊 Health check: http://localhost:%s/healthz", cfg.Port)
	log.Printf("🏃 Activities API: http://localhost:%s/activities", cfg.Port)
	log.Printf("🗄�E�E MongoDB: %s/%s", cfg.Mongo.URI, cfg.Mongo.DB)
	log.Printf("🐰 RabbitMQ: %s:%s (exchange: %s)", cfg.RabbitMQ.Host, cfg.RabbitMQ.Port, cfg.RabbitMQ.Exchange)

	// Iniciar servidor en goroutine para poder manejar shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❁EServer error: %v", err)
		}
	}()

	// Escuchar señales del sistema para graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("⏳ Shutting down server...")

	// Tiempo máximo para completar shutdown
	ctxShutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Cerrar servidor HTTP
	if err := srv.Shutdown(ctxShutdown); err != nil {
		log.Printf("❁EServer shutdown failed: %v", err)
	}

	// Detener publish queue
	if publishQueue != nil {
		log.Println("⏳ Stopping publish queue...")
		publishQueue.Stop()
	}

	// Cerrar cliente RabbitMQ
	if rabbitClient != nil {
		log.Println("⏳ Closing RabbitMQ connection...")
		if err := rabbitClient.Close(); err != nil {
			log.Printf("⚠ Error closing RabbitMQ client: %v", err)
		}
	}

	log.Println("✁EServer exited properly")
}
