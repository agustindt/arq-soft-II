package main

import (
	"context"
	"log"
	"net/http"
	"reservations/clients"
	"reservations/config"
	"reservations/controllers"
	"reservations/middleware"
	"reservations/repository"
	"reservations/services"
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

	// services
	ReservaService := services.NewReservasService(ReservasMongoRepo, reservasQueue, nil)

	// controllers
	ReservaController := controllers.NewReservasController(&ReservaService)

	// Configurar router HTTP con Gin
	router := gin.Default()

	// Middleware
	router.Use(middleware.CORSMiddleware)

	// Health check endpoint
	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Router
	// GET /Reservas - listar todos los Reservas
	router.GET("/reservas", ReservaController.GetReservas)

	// POST /Reservas - crear nuevo Reserva
	router.POST("/reservas", ReservaController.CreateReserva)

	// GET /Reservas/:id - obtener Reserva por ID
	router.GET("/reservas/:id", ReservaController.GetReservaByID)

	// PUT /Reservas/:id - actualizar Reserva existente
	router.PUT("/reservas/:id", ReservaController.UpdateReserva)

	// DELETE /Reservas/:id - eliminar Reserva
	router.DELETE("/reservas/:id", ReservaController.DeleteReserva)

	// Configuración del server
	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
	}

	log.Printf(" API listening on port %s", cfg.Port)
	log.Printf(" Health check: http://localhost:%s/healthz", cfg.Port)
	log.Printf(" Reservas API: http://localhost:%s/Reservas", cfg.Port)

	// Iniciar servidor
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}
