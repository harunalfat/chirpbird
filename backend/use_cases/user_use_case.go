package usecases

import (
	"context"

	"github.com/google/uuid"
	"github.com/harunalfat/chirpbird/backend/entities"
	"github.com/harunalfat/chirpbird/backend/helpers"
	"github.com/harunalfat/chirpbird/backend/presentation/persistence"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserUseCase struct {
	channelUseCase *ChannelUseCase
	node           NodeWrapper
	userRepo       persistence.UserRepository
}

func NewUserUseCase(channelUseCase *ChannelUseCase, node NodeWrapper, userRepo persistence.UserRepository) *UserUseCase {
	return &UserUseCase{
		channelUseCase,
		node,
		userRepo,
	}
}

func (uc *UserUseCase) Fetch(ctx context.Context, userID uuid.UUID) (entities.User, error) {
	return uc.userRepo.Fetch(ctx, userID)
}

func (uc *UserUseCase) EmbedChannelIfNotExist(ctx context.Context, user entities.User, channel entities.Channel) (res entities.User, err error) {
	if !helpers.IsExistsInUUIDArray(user.ChannelIDs, channel.ID) {
		user.ChannelIDs = append(user.ChannelIDs, channel.ID)
		res, err = uc.userRepo.Update(ctx, user.ID, user)
	}

	return
}

func (uc *UserUseCase) EmbedChannelToMultipleUsersIfNotExist(ctx context.Context, userIDs []uuid.UUID, channelID uuid.UUID) (err error) {
	users, err := uc.userRepo.FetchMultiple(ctx, userIDs)
	if err != nil {
		return
	}

	channel, err := uc.channelUseCase.Fetch(ctx, channelID)
	if err != nil {
		return
	}

	for _, user := range users {
		_, err = uc.EmbedChannelIfNotExist(ctx, user, channel)
		if err != nil {
			return
		}
	}
	return
}

func (uc *UserUseCase) CreateIfUsernameNotExist(ctx context.Context, user entities.User) (result entities.User, err error) {
	result, err = uc.userRepo.FetchByUsername(ctx, user.Username)
	if err != nil && err != mongo.ErrNoDocuments {
		return
	}

	if result.Username == "" {
		result, err = uc.userRepo.Insert(ctx, user)
	}
	return
}

func (uc *UserUseCase) SubsribeUserConnectionToChannel(ctx context.Context, userID uuid.UUID, channelID uuid.UUID) error {
	return uc.node.SubscribeClientToChannel(ctx, userID, channelID)
}
