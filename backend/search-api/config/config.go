package config

import (
	"log"
	"os"
)

type Config struct {
	Port          string
	SolrURL       string
	SolrCore      string
	MemcachedURL  string
	RabbitURL     string
	RabbitQueue   string
	ActivitiesAPI string
	LocalCacheTTL int // segundos
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func Load() Config {
	cfg := Config{
		Port:          getenv("PORT", "8083"),
		SolrURL:       getenv("SOLR_URL", "http://solr:8983/solr"),
		SolrCore:      getenv("SOLR_CORE", "activities"),
		MemcachedURL:  getenv("MEMCACHED_URL", "memcached:11211"),
		RabbitURL:     getenv("RABBITMQ_URL", "amqp://admin:admin@rabbitmq:5672/"),
		RabbitQueue:   getenv("RABBITMQ_QUEUE", "entity.events"),
		ActivitiesAPI: getenv("ACTIVITIES_API_URL", "http://activities-api:8082"),
		LocalCacheTTL: 60,
	}
	log.Printf("[cfg] solr=%s core=%s memcached=%s queue=%s", cfg.SolrURL, cfg.SolrCore, cfg.MemcachedURL, cfg.RabbitQueue)
	return cfg
}
