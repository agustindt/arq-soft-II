package config

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoURI string
var DatabaseName string
var MongoClient *mongo.Client

func init() {
	MongoURI = getenvDefault("MONGO_URI", "mongodb://localhost:27017")
	DatabaseName = getenvDefault("DATABASE", "reservationsdb")
}

func ConnectMongo(ctx context.Context) (*mongo.Client, error) {
	cctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(cctx, options.Client().ApplyURI(MongoURI))
	if err != nil {
		return nil, err
	}
	// ping
	if err := client.Ping(cctx, nil); err != nil {
		return nil, err
	}
	MongoClient = client
	log.Printf("connected to mongo %s/%s", MongoURI, DatabaseName)
	return client, nil
}

func getenvDefault(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
