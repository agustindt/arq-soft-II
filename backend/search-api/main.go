package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"arq-soft-II/backend/search-api/config"
	"arq-soft-II/backend/search-api/controllers"
	"arq-soft-II/backend/search-api/services"
	"arq-soft-II/config/cache"
	"arq-soft-II/config/httpx"
	"arq-soft-II/config/rabbitmq"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}

	// 🔹 Conectarse a Memcached
	memcachedURL := os.Getenv("MEMCACHED_URL")
	if memcachedURL == "" {
		memcachedURL = "memcached:11211" // valor por defecto si no hay env
	}

	memc, err := cache.New(memcachedURL)
	if err != nil {
		log.Fatal("❌ Error conectando con Memcached:", err)
	}
	defer log.Println("✅ Conexión Memcached cerrada.")

	// 🔹 Crear service y controller
	service := services.NewSearchService(memc)
	controller := controllers.NewSearchController(service)

	// 🔹 Conectarse a RabbitMQ
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://admin:admin@rabbitmq:5672/" // valor por defecto
	}

	mq, err := rabbitmq.New(rabbitURL)
	if err != nil {
		log.Fatal("❌ Error conectando a RabbitMQ:", err)
	}
	defer mq.Close()

	// Declarar exchange y queue
	err = mq.DeclareSetup("entity.events", "search-sync", "activities.*")
	if err != nil {
		log.Fatal("❌ Error declarando exchange/queue:", err)
	}

	// Iniciar el consumer (escucha eventos de Activities)
	config.StartRabbitConsumer(mq, service)

	// 🔹 Endpoints públicos
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Search API is running")
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Search API - Sports Activities Platform")
	})

	http.HandleFunc("/search", controller.HandleSearch)

	// 🔹 Endpoint interno protegido por X-Service-Token
	http.Handle("/internal/reindex",
		httpx.RequireServiceToken(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "Petición interna autorizada (reindex OK)")
		})),
	)

	log.Printf("🚀 Search API corriendo en puerto %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
