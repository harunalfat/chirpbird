package mongodb

import (
	"context"
	"errors"

	"github.com/harunalfat/chirpbird/backend/entities"
	"github.com/harunalfat/chirpbird/backend/presentation/persistence"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	CHAT = "chat"
)

type MongodbMessageRepository struct {
	client *mongo.Client
}

func NewMongodbMessageRepository(client *mongo.Client) persistence.MessageRepository {
	return &MongodbMessageRepository{
		client,
	}
}

func (repo *MongodbMessageRepository) FetchFromGroup(ctx context.Context, groupID string) ([]entities.Message, error) {

	return nil, errors.New("NOT READY FUNCTION")
}

func (repo *MongodbMessageRepository) Insert(ctx context.Context, msgArg entities.Message) (entities.Message, error) {
	var msg messageDTO
	msg.FromEntity(msgArg)

	res, err := repo.client.Database(CHAT).Collection(msgArg.ChannelID).InsertOne(ctx, &msg)
	if err != nil {
		return entities.Message{}, err
	}

	msg.ID = res.InsertedID.(primitive.ObjectID)
	return msg.ToEntity(), nil
}
