package usecases

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/harunalfat/chirpbird/backend/entities"
	"github.com/harunalfat/chirpbird/backend/presentation/persistence"
	"go.mongodb.org/mongo-driver/mongo"
)

type ChannelUseCase struct {
	channelRepo    persistence.ChannelRepository
	messageUseCase *MessageUseCase
}

func NewChannelUseCase(channelRepo persistence.ChannelRepository, messageUseCase *MessageUseCase) *ChannelUseCase {
	return &ChannelUseCase{
		channelRepo,
		messageUseCase,
	}
}

func (uc *ChannelUseCase) Fetch(ctx context.Context, id string) (entities.Channel, error) {
	return uc.channelRepo.Fetch(ctx, id)
}

func (uc *ChannelUseCase) FetchByName(ctx context.Context, name string) (entities.Channel, error) {
	return uc.channelRepo.FetchByName(ctx, name)
}

func (uc *ChannelUseCase) UpdateChannelWithMessage(ctx context.Context, senderID string, channelID string, message string) error {
	input := entities.Message{
		Sender: entities.User{
			Base: entities.Base{
				ID: senderID,
			},
		},
		ChannelID: channelID,
		Data:      message,
		Base: entities.Base{
			CreatedAt: time.Now(),
		},
	}

	if _, err := uc.messageUseCase.Store(ctx, input); err != nil {
		return err
	}

	channel, err := uc.channelRepo.Fetch(ctx, channelID)
	if err != nil {
		return err
	}

	_, err = uc.channelRepo.Update(ctx, channelID, channel)
	return err
}

func (uc *ChannelUseCase) CreateIfNameNotExist(ctx context.Context, channel entities.Channel, creator entities.User) (result entities.Channel, err error) {
	result, err = uc.channelRepo.FetchByName(ctx, channel.Name)
	if err != nil && err != mongo.ErrNoDocuments {
		return entities.Channel{}, err
	}

	if result.Name == "" {
		channel.CreatorID = creator.ID
		channel.CreatedAt = time.Now()
		hash := sha256.Sum256([]byte(channel.Name))
		channel.HashIdentifier = fmt.Sprintf("%x", hash[:])

		result, err = uc.channelRepo.Insert(ctx, channel)
		if err != nil {
			return entities.Channel{}, err
		}
	}

	return result, err
}
