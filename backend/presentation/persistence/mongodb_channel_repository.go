package persistence

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/harunalfat/chirpbird/backend/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ChannelRepository interface {
	Fetch(ctx context.Context, channelID uuid.UUID) (entities.Channel, error)
	FetchByName(ctx context.Context, channelName string) (entities.Channel, error)
	Insert(context.Context, entities.Channel) (entities.Channel, error)
	Update(ctx context.Context, channelID uuid.UUID, updated entities.Channel) (entities.Channel, error)
}

type MongodbChannelRepository struct {
	client *mongo.Client
}

func NewMongodbChannelRepository(client *mongo.Client) ChannelRepository {
	return &MongodbChannelRepository{
		client,
	}
}

func (repo *MongodbChannelRepository) Fetch(ctx context.Context, channelID uuid.UUID) (entities.Channel, error) {
	var channel entities.Channel
	filter := bson.D{{Key: "id", Value: channelID}}

	err := repo.client.
		Database(DB_GENERAL).
		Collection(COLLECTION_CHANNEL).
		FindOne(ctx, filter).
		Decode(&channel)

	return channel, err
}

func (repo *MongodbChannelRepository) FetchByName(ctx context.Context, channelName string) (entities.Channel, error) {
	var channel entities.Channel
	filter := bson.D{{Key: "name", Value: channelName}}

	err := repo.client.
		Database(DB_GENERAL).
		Collection(COLLECTION_CHANNEL).
		FindOne(ctx, filter).
		Decode(&channel)

	return channel, err
}

func (repo *MongodbChannelRepository) Insert(ctx context.Context, channelArg entities.Channel) (entities.Channel, error) {
	_, err := repo.client.
		Database(DB_GENERAL).
		Collection(COLLECTION_CHANNEL).
		InsertOne(ctx, channelArg)
	if err != nil {
		return entities.Channel{}, err
	}

	return channelArg, err
}

func (repo *MongodbChannelRepository) Update(ctx context.Context, channelID uuid.UUID, channel entities.Channel) (entities.Channel, error) {
	channel.UpdatedAt = time.Now()

	filter := bson.D{{
		Key:   "_id",
		Value: channelID,
	}}
	update := bson.D{{
		Key: "$set", Value: channel,
	}}

	_, err := repo.client.Database(DB_GENERAL).Collection(COLLECTION_CHANNEL).UpdateOne(ctx, filter, update)
	return channel, err
}
