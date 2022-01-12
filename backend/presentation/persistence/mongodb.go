package persistence

import (
	"context"
	"log"
	"os"

	"github.com/harunalfat/chirpbird/backend/env"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	DB_GENERAL = "general"

	COLLECTION_USER    = "user"
	COLLECTION_CHANNEL = "channel"
)

func NewMongoClient() (*mongo.Client, error) {
	clientOps := options.Client().ApplyURI(os.Getenv(env.MONGODB_CONN_URL))
	client, err := mongo.Connect(context.Background(), clientOps)
	if err != nil {
		log.Fatalf("Could not connect to mongo instance\n%s", err)
		return nil, err
	}

	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatalf("Could not ping mongo instance\n%s", err)
	} else {
		log.Println("Successfully ping mongo instance")
	}
	return client, err
}
