package persistence

import (
	"context"

	"github.com/harunalfat/chirpbird/backend/entities"
)

type UserRepository interface {
	Fetch(ctx context.Context, userID string) (entities.User, error)
	FetchByUsername(ctx context.Context, username string) (entities.User, error)
	FetchMultiple(context.Context, []string) ([]entities.User, error)
	Update(ctx context.Context, userID string, updated entities.User) (entities.User, error)
	Insert(context.Context, entities.User) (entities.User, error)
}

type ChannelRepository interface {
	Fetch(ctx context.Context, channelID string) (entities.Channel, error)
	FetchByName(ctx context.Context, channelName string) (entities.Channel, error)
	Insert(context.Context, entities.Channel) (entities.Channel, error)
	Update(ctx context.Context, channelID string, updated entities.Channel) (entities.Channel, error)
}

type MessageRepository interface {
	FetchFromGroup(ctx context.Context, groupID string) ([]entities.Message, error)
	Insert(context.Context, entities.Message) (entities.Message, error)
}
