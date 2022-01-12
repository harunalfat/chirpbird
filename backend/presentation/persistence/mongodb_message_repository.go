package persistence

import (
	"context"
	"errors"

	"github.com/harunalfat/chirpbird/backend/entities"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/x/mongo/driver/uuid"
)

const (
	CHAT = "chat"
)

type MessageRepository interface {
	FetchFromGroup(ctx context.Context, groupID uuid.UUID) ([]entities.Message, error)
	Insert(context.Context, entities.Message) (entities.Message, error)
}

type MongodbMessageRepository struct {
	client *mongo.Client
}

func NewMongodbMessageRepository(client *mongo.Client) MessageRepository {
	return &MongodbMessageRepository{
		client,
	}
}

func (repo *MongodbMessageRepository) FetchFromGroup(ctx context.Context, groupID uuid.UUID) ([]entities.Message, error) {

	return nil, errors.New("NOT READY FUNCTION")
}

func (repo *MongodbMessageRepository) Insert(ctx context.Context, msgArg entities.Message) (entities.Message, error) {
	_, err := repo.client.Database(CHAT).Collection(msgArg.ChannelID).InsertOne(ctx, msgArg)
	if err != nil {
		return entities.Message{}, err
	}

	return msgArg, nil
}
