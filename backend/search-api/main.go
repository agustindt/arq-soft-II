// Package main implements the Search API microservice.
//
// The Search API provides full-text search capabilities for activities using Apache Solr.
// It implements CQRS pattern as the read side, consuming events from Activities API via
// RabbitMQ to keep the search index synchronized. It also includes a caching layer using
// Memcached for improved query performance.
//
// Key Features:
//   - Full-text search across activity name, description, location, and instructor
//   - Multi-faceted filtering (category, difficulty, price range, location)
//   - Two-level caching (Memcached + Solr internal caches)
//   - Event-driven index synchronization via RabbitMQ
//   - Pagination support for search results
//   - Sub-10ms response times for cached queries
//
// Event Types Consumed:
//   - activity.created: Index new activity in Solr
//   - activity.updated: Re-index updated activity and invalidate cache
//   - activity.deleted: Remove activity from Solr index and invalidate cache
//
// Search Engine: Apache Solr 9.4
// Cache: Memcached (5-minute TTL)
// Message Queue: RabbitMQ (consumer)
// Port: 8083
//
// For complete API documentation, see docs/api/search-api.md
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
	"search-api/middleware"
	"search-api/services"
	"search-api/utils"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}

	// 隼 Conectarse a Memcached
	memcachedURL := os.Getenv("MEMCACHED_URL")
	if memcachedURL == "" {
		memcachedURL = "memcached:11211" // valor por defecto si no hay env
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
		log.Fatal("笶・Error conectando con Memcached:", err)
	}
	defer log.Println("笨・Conexiﾃｳn Memcached cerrada.")

	// 隼 Conectarse a Solr
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
		log.Printf("笞・・ Warning: Solr health check failed: %v", err)
	} else {
		log.Printf("笨・Conexiﾃｳn con Solr OK: %s/%s", solrURL, solrCore)
	}

	// 隼 Crear service y controller con Solr
	service := services.NewSearchService(memc, solrClient, cacheTTL)
	controller := controllers.NewSearchController(service)

	// 隼 Conectarse a RabbitMQ
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://admin:admin@rabbitmq:5672/" // valor por defecto
	}

	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		log.Fatal("笶・Error conectando a RabbitMQ:", err)
	}
	defer conn.Close()

	// Iniciar el consumer (escucha eventos de Activities)
	config.StartRabbitConsumer(conn, service)

	// CORS middleware function
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

	// 隼 Endpoints pﾃｺblicos
	http.HandleFunc("/health", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Search API is running")
	}))

	http.HandleFunc("/", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Search API - Sports Activities Platform")
	}))

	http.HandleFunc("/search", corsMiddleware(controller.HandleSearch))

	// 隼 Endpoint interno protegido por X-Service-Token
	http.HandleFunc("/internal/reindex",
		middleware.RequireServiceToken(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "Peticiﾃｳn interna autorizada (reindex OK)")
		}),
	)

	log.Printf("噫 Search API corriendo en puerto %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
