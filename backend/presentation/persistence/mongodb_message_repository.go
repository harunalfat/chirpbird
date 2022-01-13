package persistence

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/harunalfat/chirpbird/backend/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	CHAT = "chat"
)

type MessageRepository interface {
	FetchFromChannel(ctx context.Context, channelID string) ([]entities.Message, error)
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

func getChannelKey(channelID string) string {
	return fmt.Sprintf("channel:%s", channelID)
}

func (repo *MongodbMessageRepository) FetchFromChannel(ctx context.Context, channelID string) (result []entities.Message, err error) {
	cursor, err := repo.client.Database(DB_NAME).Collection(getChannelKey(channelID)).Find(ctx, bson.D{})
	if err != nil {
		if err == mongo.ErrNilDocument || err == mongo.ErrNoDocuments {
			return make([]entities.Message, 0), nil
		}
		return
	}

	err = cursor.All(ctx, &result)
	return
}

func (repo *MongodbMessageRepository) Insert(ctx context.Context, msgArg entities.Message) (entities.Message, error) {
	msgArg.ID = uuid.New().String()
	msgArg.CreatedAt = time.Now()

	_, err := repo.client.Database(DB_NAME).
		Collection(getChannelKey(msgArg.ChannelID)).
		InsertOne(ctx, msgArg)
	if err != nil {
		return entities.Message{}, err
	}

	return msgArg, nil
}
