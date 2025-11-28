package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port             string
	Mongo            MongoConfig
	RabbitMQ         RabbitMQConfig
	UsersAPIURL      string
	ActivitiesAPIURL string
}

type MongoConfig struct {
	URI string
	DB  string
}

type RabbitMQConfig struct {
	Username  string
	Password  string
	QueueName string
	Host      string
	Port      string
}

func Load() Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading .env file")
	}

	return Config{
		Port:             getEnv("PORT", "8080"),
		UsersAPIURL:      getEnv("USERS_API_URL", "http://users-api:8081"),
		ActivitiesAPIURL: getEnv("ACTIVITIES_API_URL", "http://activities-api:8082"),
		Mongo: MongoConfig{
			URI: getEnv("MONGO_URI", "mongodb://localhost:27017"),
			DB:  getEnv("MONGO_DB", "reservas"),
		},
		RabbitMQ: RabbitMQConfig{
			Username:  getEnv("RABBITMQ_USER", "admin"),
			Password:  getEnv("RABBITMQ_PASS", "admin"),
			QueueName: getEnv("RABBITMQ_QUEUE_NAME", "reservas-news"),
			Host:      getEnv("RABBITMQ_HOST", "localhost"),
			Port:      getEnv("RABBITMQ_PORT", "5672"),
		},
	}
}

func getEnv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
