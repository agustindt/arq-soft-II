package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	Mongo       MongoConfig
	RabbitMQ    RabbitMQConfig
	UsersAPIURL string
	JWTSecret   string
}

type MongoConfig struct {
	URI        string
	DB         string
	Collection string
}

type RabbitMQConfig struct {
	Username string
	Password string
	Exchange string
	Host     string
	Port     string
}

func Load() Config {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found or error loading .env file")
	}

	return Config{
		Port:        getEnv("PORT", "8082"),
		UsersAPIURL: getEnv("USERS_API_URL", "http://localhost:8081"),
		JWTSecret:   getEnv("JWT_SECRET", "your-super-secret-jwt-key"),
		Mongo: MongoConfig{
			URI:        getEnv("MONGO_URI", "mongodb://localhost:27017"),
			DB:         getEnv("MONGO_DB", "activitiesdb"),
			Collection: getEnv("MONGO_COLLECTION", "activities"),
		},
		RabbitMQ: RabbitMQConfig{
			Username: getEnv("RABBITMQ_USER", "admin"),
			Password: getEnv("RABBITMQ_PASS", "admin123"),
			Exchange: getEnv("RABBITMQ_EXCHANGE", "entity.events"),
			Host:     getEnv("RABBITMQ_HOST", "localhost"),
			Port:     getEnv("RABBITMQ_PORT", "5672"),
		},
	}
}

func getEnv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
