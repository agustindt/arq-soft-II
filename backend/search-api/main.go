package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"search-api/clients"
	"search-api/config"
	"search-api/controllers"
	"search-api/services"
	"search-api/utils"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}

	// Conectarse a Memcached
	memcachedURL := os.Getenv("MEMCACHED_URL")
	if memcachedURL == "" {
		memcachedURL = "memcached:11211"
	}

	ttlSeconds := 60
	if ttlStr := os.Getenv("CACHE_TTL_SECONDS"); ttlStr != "" {
		if parsed, err := strconv.Atoi(ttlStr); err == nil && parsed > 0 {
			ttlSeconds = parsed
		}
	}

	cacheTTL := time.Duration(ttlSeconds) * time.Second

	memc, err := utils.NewCache(memcachedURL, cacheTTL)
	if err != nil {
		log.Fatal("Error conectando con Memcached:", err)
	}
	defer log.Println("Conexión Memcached cerrada.")

	// Conectarse a Solr
	solrURL := os.Getenv("SOLR_URL")
	if solrURL == "" {
		solrURL = "http://solr:8983/solr"
	}
	solrCore := os.Getenv("SOLR_CORE")
	if solrCore == "" {
		solrCore = "activities"
	}

	solrClient := clients.NewSolrClient(solrURL, solrCore)

	// Health check de Solr
	if err := solrClient.HealthCheck(); err != nil {
		log.Printf("Warning: Solr health check failed: %v", err)
	} else {
		log.Printf("Conexión con Solr OK: %s/%s", solrURL, solrCore)
	}

	// Crear service y controller
	service := services.NewSearchService(memc, solrClient, cacheTTL)
	controller := controllers.NewSearchController(service)

	// Conectarse a RabbitMQ
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://admin:admin@rabbitmq:5672/"
	}

	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		log.Fatal("Error conectando a RabbitMQ:", err)
	}
	defer conn.Close()

	// Iniciar consumer
	config.StartRabbitConsumer(conn, service)

	// CORS middleware
	corsMiddleware := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next(w, r)
		}
	}

	// Rutas
	http.HandleFunc("/search", corsMiddleware(controller.HandleSearch))

	fmt.Printf("Search API running on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
