package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"reservations/clients"
	"reservations/config"
	"reservations/controllers"
	"reservations/middleware"
	"reservations/repository"
	"reservations/services"
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

	// services
	ReservaService := services.NewReservasService(ReservasMongoRepo, publishQueue)

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
	usersAPI := cfg.UsersAPIURL

	// Router
	// GET /Reservas - listar todos los Reservas
	router.GET("/reservas", ReservaController.GetReservas)

	// POST /Reservas - crear nuevo Reserva
	router.POST("/reservas", middleware.AdminOnly(usersAPI), ReservaController.CreateReserva)

	// GET /Reservas/:id - obtener Reserva por ID
	router.GET("/reservas/:id", ReservaController.GetReservaByID)

	// PUT /Reservas/:id - actualizar Reserva existente
	router.PUT("/reservas/:id", middleware.AdminOnly(usersAPI), ReservaController.UpdateReserva)

	// DELETE /Reservas/:id - eliminar Reserva
	router.DELETE("/reservas/:id", middleware.AdminOnly(usersAPI), ReservaController.DeleteReserva)

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
