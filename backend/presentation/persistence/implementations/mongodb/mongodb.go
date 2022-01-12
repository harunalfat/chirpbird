package mongodb

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/harunalfat/chirpbird/backend/env"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	DB_GENERAL = "general"

	COLLECTION_USER    = "user"
	COLLECTION_CHANNEL = "channel"
)

type BaseDTO struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	CreatedAt time.Time          `bson:"created,omitempty"`
	UpdatedAt time.Time          `bson:"updatedAt,omitempty"`
}

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
