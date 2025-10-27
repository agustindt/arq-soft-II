package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"search-api/clients"
	"search-api/config"
	"search-api/controllers"
	"search-api/middleware"
	"search-api/repository"
	"search-api/services"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	// Repo Solr + Cache dual
	repo := repository.NewSolrRepository(cfg.SolrURL, cfg.SolrCore)
	cache := services.NewDualCache(cfg.MemcachedURL, cfg.LocalCacheTTL)
	svc := services.NewSearchService(repo, cache)
	ctrl := controllers.NewSearchController(svc)

	// Consumer RabbitMQ: sincroniza índice
	ctx := context.Background()
	actCli := clients.NewActivitiesClient(cfg.ActivitiesAPI)
	consumer, err := clients.NewConsumer(cfg.RabbitURL, cfg.RabbitQueue, repo, actCli)
	if err != nil {
		log.Fatalf("rabbit consumer: %v", err)
	}
	if err := consumer.Start(ctx); err != nil {
		log.Fatalf("consumer start: %v", err)
	}

	// HTTP
	r := gin.Default()
	r.Use(middleware.CORS())
	r.GET("/health", ctrl.Health)
	r.GET("/search", ctrl.Search)

	srv := &http.Server{Addr: ":" + cfg.Port, Handler: r, ReadHeaderTimeout: 10 * time.Second}
	log.Printf("search-api listening on :%s", cfg.Port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}
