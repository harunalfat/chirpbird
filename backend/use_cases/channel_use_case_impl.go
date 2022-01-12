package usecases

import (
	"context"
	"time"

	"github.com/harunalfat/chirpbird/backend/entities"
	"github.com/harunalfat/chirpbird/backend/presentation/persistence"
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
		SenderID:  senderID,
		ChannelID: channelID,
		Text:      message,
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

func (uc *ChannelUseCase) Create(ctx context.Context, channel entities.Channel, creator entities.User) (entities.Channel, error) {
	channel.CreatorID = creator.ID
	return uc.channelRepo.Insert(ctx, channel)
}
