package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"

	"reservations/config"
	"reservations/handlers"
	"reservations/messaging"
	"reservations/repository"
	"reservations/services"
)

func main() {
	// context for startup
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// connect mongo
	client, err := config.ConnectMongo(ctx)
	if err != nil {
		log.Fatalf("mongo connect failed: %v", err)
	}
	db := client.Database(config.DatabaseName)
	col := db.Collection("reservations")

	// rabbit
	rabbitURL := getenvDefault("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/")
	exchange := getenvDefault("RABBITMQ_EXCHANGE", "reservations")
	pub, err := messaging.NewRabbitPublisher(rabbitURL, exchange)
	if err != nil {
		log.Fatalf("rabbitmq init failed: %v", err)
	}
	defer pub.Close()

	// repo
	repo := repository.NewMongoReservationRepo(col)

	// users service
	usersBase := getenvDefault("USERS_SERVICE_URL", "http://localhost:8081")
	users := services.NewHTTPUserService(usersBase)

	// business service (tasks from env)
	tasks := 5
	svc := services.NewReservationService(repo, pub, users, tasks)

	// handlers
	rh := handlers.NewReservationHandler(svc)

	r := mux.NewRouter()
	rh.Register(r)

	addr := getenvDefault("HTTP_ADDR", ":8080")
	srv := &http.Server{Addr: addr, Handler: r}

	log.Printf("listening %s", addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}

func getenvDefault(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
