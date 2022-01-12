package usecases

import (
	"context"

	"github.com/harunalfat/chirpbird/backend/entities"
)

type ChannelUseCase interface {
	Create(context.Context, entities.Channel, entities.User) (entities.Channel, error)
	Fetch(ctx context.Context, channelID string) (entities.Channel, error)
	FetchByName(ctx context.Context, name string) (entities.Channel, error)
	UpdateChannelWithMessage(ctx context.Context, senderID string, channelID string, message string) error
}

type MessageUseCase interface {
	Store(context.Context, entities.Message) (entities.Message, error)
}

type UserUseCase interface {
	CreateIfUsernameNotExist(context.Context, entities.User) (entities.User, error)
	EmbedChannelIfNotExist(ctx context.Context, user entities.User, channel entities.Channel) (entities.User, error)
	EmbedChannelToMultipleUsersIfNotExist(ctx context.Context, usernames []string, channelName string) (err error)
	Fetch(ctx context.Context, userID string) (entities.User, error)
	SubsribeUserConnectionToChannel(ctx context.Context, userID string, channelID string) error
}
