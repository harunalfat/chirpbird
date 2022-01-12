package mongodb

import (
	"context"
	"log"

	"github.com/harunalfat/chirpbird/backend/entities"
	"github.com/harunalfat/chirpbird/backend/helpers"
	"github.com/harunalfat/chirpbird/backend/presentation/persistence"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	ID       = "_id"
	USERNAME = "username"
)

type MongodbUserRepository struct {
	client *mongo.Client
}

func NewMongodbUserRepository(client *mongo.Client) persistence.UserRepository {
	return &MongodbUserRepository{
		client,
	}
}

func (repo *MongodbUserRepository) Fetch(ctx context.Context, userIDHex string) (entities.User, error) {
	var user userDTO
	filter := bson.D{{Key: ID, Value: helpers.ObjectIDFromHex(userIDHex)}}
	err := repo.client.
		Database(DB_GENERAL).
		Collection(COLLECTION_USER).
		FindOne(ctx, filter).
		Decode(&user)
	log.Println(user)
	return user.ToEntity(), err
}

func (repo *MongodbUserRepository) FetchByUsername(ctx context.Context, username string) (entities.User, error) {
	var user userDTO
	filter := bson.D{{Key: USERNAME, Value: username}}
	err := repo.client.
		Database(DB_GENERAL).
		Collection(COLLECTION_USER).
		FindOne(ctx, filter).
		Decode(&user)

	return user.ToEntity(), err
}

func (repo *MongodbUserRepository) FetchMultiple(ctx context.Context, userIDHexes []string) ([]entities.User, error) {
	var users userDTOs
	filter := bson.D{{
		Key: ID,
		Value: bson.D{{
			Key:   "$in",
			Value: helpers.ObjectIDsFromHexes(userIDHexes),
		}},
	}}

	cursor, err := repo.client.
		Database(DB_GENERAL).
		Collection(COLLECTION_USER).
		Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &users)
	return users.ToEntity(), err
}

func (repo *MongodbUserRepository) Update(ctx context.Context, userIDHex string, updated entities.User) (entities.User, error) {
	var user userDTO
	user.FromEntity(updated)

	filter := bson.D{{Key: ID, Value: helpers.ObjectIDFromHex(userIDHex)}}
	err := repo.client.
		Database(DB_GENERAL).
		Collection(COLLECTION_USER).
		FindOneAndUpdate(ctx, filter, user).
		Decode(&user)

	return user.ToEntity(), err
}

func (repo *MongodbUserRepository) Insert(ctx context.Context, userArg entities.User) (entities.User, error) {
	var user userDTO
	user.FromEntity(userArg)

	res, err := repo.client.
		Database(DB_GENERAL).
		Collection(COLLECTION_USER).
		InsertOne(ctx, user)
	if err != nil {
		log.Printf("Failed to insert User\n%s", err)
		return entities.User{}, nil
	}

	user.ID = res.InsertedID.(primitive.ObjectID)
	return user.ToEntity(), nil
}
