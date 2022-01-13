package persistence

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/harunalfat/chirpbird/backend/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	ID       = "id"
	USERNAME = "username"
)

type UserRepository interface {
	Fetch(ctx context.Context, userID string) (entities.User, error)
	FetchByUsername(ctx context.Context, username string) (entities.User, error)
	FetchMultiple(context.Context, []string) ([]entities.User, error)
	Update(ctx context.Context, userID string, updated entities.User) (entities.User, error)
	Insert(context.Context, entities.User) (entities.User, error)
}

type MongodbUserRepository struct {
	client *mongo.Client
}

func NewMongodbUserRepository(client *mongo.Client) UserRepository {
	return &MongodbUserRepository{
		client,
	}
}

func (repo *MongodbUserRepository) Fetch(ctx context.Context, userID string) (entities.User, error) {
	var user entities.User
	filter := bson.D{{Key: ID, Value: userID}}
	err := repo.client.
		Database(DB_NAME).
		Collection(COLLECTION_USER).
		FindOne(ctx, filter).
		Decode(&user)

	return user, err
}

func (repo *MongodbUserRepository) FetchByUsername(ctx context.Context, username string) (entities.User, error) {
	var user entities.User
	filter := bson.D{{Key: USERNAME, Value: username}}
	err := repo.client.
		Database(DB_NAME).
		Collection(COLLECTION_USER).
		FindOne(ctx, filter).
		Decode(&user)

	return user, err
}

func (repo *MongodbUserRepository) FetchMultiple(ctx context.Context, userIDs []string) ([]entities.User, error) {
	var users []entities.User
	filter := bson.D{{
		Key: ID,
		Value: bson.D{{
			Key:   "$in",
			Value: userIDs,
		}},
	}}

	cursor, err := repo.client.
		Database(DB_NAME).
		Collection(COLLECTION_USER).
		Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &users)
	return users, err
}

func (repo *MongodbUserRepository) Update(ctx context.Context, userID string, updated entities.User) (entities.User, error) {
	filter := bson.D{{Key: ID, Value: userID}}
	update := bson.D{{
		Key: "$set", Value: updated,
	}}
	err := repo.client.
		Database(DB_NAME).
		Collection(COLLECTION_USER).
		FindOneAndUpdate(ctx, filter, update).
		Decode(&updated)

	return updated, err
}

func (repo *MongodbUserRepository) Insert(ctx context.Context, userArg entities.User) (entities.User, error) {
	userArg.CreatedAt = time.Now()
	userArg.ID = uuid.New().String()

	_, err := repo.client.
		Database(DB_NAME).
		Collection(COLLECTION_USER).
		InsertOne(ctx, userArg)
	if err != nil {
		log.Printf("Failed to insert User\n%s", err)
		return entities.User{}, nil
	}

	return userArg, nil
}
