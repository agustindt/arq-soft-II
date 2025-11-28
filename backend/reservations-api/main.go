// Package main implements the Reservations API microservice.
//
// The Reservations API manages activity reservations and bookings. It provides endpoints
// for creating, viewing, updating, and managing reservations with capacity validation
// and status workflow management.
//
// Key Features:
//   - Create reservations for activities (authenticated users)
//   - Multi-user group reservations support
//   - Capacity validation against activity max_capacity
//   - Status workflow (pendiente ↁEconfirmada ↁEcancelada)
//   - Admin endpoints for reservation management
//   - Integration with Activities API for capacity checks
//   - JWT-based authentication and authorization
//   - Graceful shutdown with proper resource cleanup
//
// Reservation Status Flow:
//   - pendiente: Initial state, awaiting confirmation
//   - confirmada: Confirmed reservation
//   - cancelada: Cancelled reservation (terminal state)
//
// Database: MongoDB 6.0
// Message Queue: RabbitMQ (optional, for future event publishing)
// Port: 8080
//
// For complete API documentation, see docs/api/reservations-api.md
package main

import (
	"reservations-api/clients"
	"reservations-api/config"
	"reservations-api/controllers"
	"reservations-api/middleware"
	"reservations-api/repository"
	"reservations-api/services"
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
	// configuracion inicial
	cfg := config.Load()

	// Context
	ctx := context.Background()

	// Capa de datos: maneja operaciones DB
	ReservasMongoRepo := repository.NewMongoReservasRepository(ctx, cfg.Mongo.URI, cfg.Mongo.DB, "Reservas")

	reservasQueue := clients.NewRabbitMQClient(
		cfg.RabbitMQ.Username,
		cfg.RabbitMQ.Password,
		cfg.RabbitMQ.QueueName,
		cfg.RabbitMQ.Host,
		cfg.RabbitMQ.Port,
	)

	// crear cola de publicación con workers y retries
	publishQueue := services.NewPublishQueue(reservasQueue, 200, 3, 200*time.Millisecond)
	// start workers (use ctx so they can be cancelled on shutdown)
	publishQueue.Start(ctx, 2)

	// activities API URL from config
	activitiesAPIURL := cfg.ActivitiesAPIURL
	if activitiesAPIURL == "" {
		activitiesAPIURL = "http://localhost:8082"
	}

	// services
	ReservaService := services.NewReservasService(ReservasMongoRepo, publishQueue, activitiesAPIURL)

	// controllers
	ReservaController := controllers.NewReservasController(&ReservaService)

	// Configurar router HTTP con Gin
	router := gin.Default()

	// Middleware
	router.Use(middleware.CORSMiddleware)

	// Health check endpoint
	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Reservations API is running",
			"service": "reservations-api",
		})
	})
	usersAPI := cfg.UsersAPIURL

	// Router
	// GET /Reservas - listar reservas (usuarios autenticados ven solo las suyas, admins ven todas)
	router.GET("/reservas", middleware.AuthRequired(usersAPI), ReservaController.GetReservas)
	// POST /Reservas - crear nuevo Reserva (cualquier usuario autenticado)
	router.POST("/reservas", middleware.AuthRequired(usersAPI), ReservaController.CreateReserva)

	// GET /Reservas/:id - obtener Reserva por ID
	router.GET("/reservas/:id", ReservaController.GetReservaByID)
	// PUT /Reservas/:id - actualizar Reserva existente
	router.PUT("/reservas/:id", middleware.AdminOnly(usersAPI), ReservaController.UpdateReserva)

	// DELETE /Reservas/:id - eliminar Reserva (usuario autenticado puede eliminar sus propias reservas)
	router.DELETE("/reservas/:id", middleware.AuthRequired(usersAPI), ReservaController.DeleteReserva)

	// GET /activities/:id/availability - obtener disponibilidad por horario
	router.GET("/activities/:id/availability", ReservaController.GetScheduleAvailability)

	// Configuración del server
	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
	}

	log.Printf(" API listening on port %s", cfg.Port)
	log.Printf(" Health check: http://localhost:%s/healthz", cfg.Port)
	log.Printf(" Reservas API: http://localhost:%s/Reservas", cfg.Port)

	// Iniciar servidor en goroutine para poder manejar shutdown tranquilamente
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// Escuchar señales del sistema para shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Tiempo máximo para completar shutdown
	ctxShutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctxShutdown); err != nil {
		log.Fatalf("server Shutdown Failed:%+v", err)
	}

	// Detener publish queue y cerrar clientes externos (RabbitMQ)
	if publishQueue != nil {
		publishQueue.Stop()
	}
	if reservasQueue != nil {
		if err := reservasQueue.Close(); err != nil {
			log.Printf("error closing rabbitmq client: %v", err)
		}
	}

	log.Println("Server exited properly")
}
