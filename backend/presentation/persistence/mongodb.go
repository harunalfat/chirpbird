package persistence

import (
	"context"
	"log"
	"os"

	"github.com/harunalfat/chirpbird/backend/env"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	DB_NAME = "chirpbird"

	COLLECTION_USER    = "users"
	COLLECTION_CHANNEL = "channels"
)

var (
	usersUniqueUsername = "users_unique_username"
	usersUniqueId       = "users_unique_id"
	channelsUniqueId    = "channels_unique_id"
	channelsUniqueHash  = "channels_unique_hash"
)

var MONGO_INDEXES = map[string]map[string]*mongo.IndexModel{
	COLLECTION_USER: {
		usersUniqueUsername: &mongo.IndexModel{
			Keys: bson.M{
				"username": 1,
			},
			Options: &options.IndexOptions{
				Name: &usersUniqueUsername,
			},
		},
		usersUniqueId: &mongo.IndexModel{
			Keys: bson.M{
				"id": 1,
			},
			Options: &options.IndexOptions{
				Name: &usersUniqueId,
			},
		},
	},

	COLLECTION_CHANNEL: {
		channelsUniqueId: &mongo.IndexModel{
			Keys: bson.M{
				"id": 1,
			},
			Options: &options.IndexOptions{
				Name: &channelsUniqueId,
			},
		},
		channelsUniqueHash: &mongo.IndexModel{
			Keys: bson.M{
				"hashIdentifier": 1,
			},
			Options: &options.IndexOptions{
				Name: &channelsUniqueHash,
			},
		},
	},
}

func newMongoClient() *mongo.Client {
	clientOps := options.Client().ApplyURI(os.Getenv(env.MONGODB_CONN_URL))
	client, err := mongo.Connect(context.Background(), clientOps)
	if err != nil {
		log.Fatalf("Cannot connect to mongo instance\n%s", err)
	}

	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatalf("Cannot ping mongo instance\n%s", err)
	}
	log.Println("Successfully ping mongo instance")
	return client
}

func buildIndexesIfNotExist(client *mongo.Client) {
	for collectionName, _ := range MONGO_INDEXES {
		collection := client.Database(DB_NAME).Collection(collectionName)
		indexSpecs, err := collection.Indexes().ListSpecifications(context.Background())
		if err != nil {
			log.Fatalf("Cannot fetch indexes\n%s", err)
		}

		existingIndexNamesMap := make(map[string]bool, len(indexSpecs))
		for _, indexSpec := range indexSpecs {
			existingIndexNamesMap[indexSpec.Name] = true
		}

		for indexName, model := range MONGO_INDEXES[collectionName] {
			// if index is included in existing indexes list, bypass the build process
			if existingIndexNamesMap[indexName] {
				continue
			}

			// build index
			_, err := collection.Indexes().CreateOne(context.Background(), *model)
			if err != nil {
				log.Fatalf("Cannot create index\n%s", err)
			}
		}
	}

}

func MongoDBInit() *mongo.Client {
	client := newMongoClient()
	buildIndexesIfNotExist(client)

	return client
}
