package mongodb

import (
	"context"
	"time"

	"github.com/harunalfat/chirpbird/backend/entities"
	"github.com/harunalfat/chirpbird/backend/helpers"
	"github.com/harunalfat/chirpbird/backend/presentation/persistence"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongodbChannelRepository struct {
	client *mongo.Client
}

func NewMongodbChannelRepository(client *mongo.Client) persistence.ChannelRepository {
	return &MongodbChannelRepository{
		client,
	}
}

func (repo *MongodbChannelRepository) Fetch(ctx context.Context, channelIDHex string) (entities.Channel, error) {
	var channel channelDTO
	filter := bson.D{{Key: "_id", Value: helpers.ObjectIDFromHex(channelIDHex)}}

	err := repo.client.
		Database(DB_GENERAL).
		Collection(COLLECTION_CHANNEL).
		FindOne(ctx, filter).
		Decode(&channel)

	return channel.ToEntity(), err
}

func (repo *MongodbChannelRepository) FetchByName(ctx context.Context, channelName string) (entities.Channel, error) {
	var channel channelDTO
	filter := bson.D{{Key: "name", Value: channelName}}

	err := repo.client.
		Database(DB_GENERAL).
		Collection(COLLECTION_CHANNEL).
		FindOne(ctx, filter).
		Decode(&channel)

	return channel.ToEntity(), err
}

func (repo *MongodbChannelRepository) Insert(ctx context.Context, channelArg entities.Channel) (entities.Channel, error) {
	var channel channelDTO
	channel.FromEntity(channelArg)

	result, err := repo.client.
		Database(DB_GENERAL).
		Collection(COLLECTION_CHANNEL).
		InsertOne(ctx, channel)
	if err != nil {
		return entities.Channel{}, err
	}

	channel.ID = result.InsertedID.(primitive.ObjectID)
	return channel.ToEntity(), err
}

func (repo *MongodbChannelRepository) Update(ctx context.Context, channelIDHex string, updated entities.Channel) (entities.Channel, error) {
	var channel channelDTO
	channel.FromEntity(updated)
	channel.UpdatedAt = time.Now()

	filter := bson.D{{
		Key:   "_id",
		Value: helpers.ObjectIDFromHex(channelIDHex),
	}}

	_, err := repo.client.Database(DB_GENERAL).Collection(COLLECTION_CHANNEL).ReplaceOne(ctx, filter, channel)
	return channel.ToEntity(), err
}
