package implementations

import (
	"context"
	"time"

	"github.com/harunalfat/chirpbird/backend/entities"
	"github.com/harunalfat/chirpbird/backend/presentation/persistence"
	usecases "github.com/harunalfat/chirpbird/backend/use_cases"
)

type ChannelUseCaseImpl struct {
	channelRepo    persistence.ChannelRepository
	messageUseCase usecases.MessageUseCase
}

func NewChannelUseCaseImpl(channelRepo persistence.ChannelRepository, messageUseCase usecases.MessageUseCase) usecases.ChannelUseCase {
	return &ChannelUseCaseImpl{
		channelRepo,
		messageUseCase,
	}
}

func (uc *ChannelUseCaseImpl) Fetch(ctx context.Context, id string) (entities.Channel, error) {
	return uc.channelRepo.Fetch(ctx, id)
}

func (uc *ChannelUseCaseImpl) FetchByName(ctx context.Context, name string) (entities.Channel, error) {
	return uc.channelRepo.FetchByName(ctx, name)
}

func (uc *ChannelUseCaseImpl) UpdateChannelWithMessage(ctx context.Context, senderID string, channelID string, message string) error {
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

func (uc *ChannelUseCaseImpl) Create(ctx context.Context, channel entities.Channel, creator entities.User) (entities.Channel, error) {
	channel.CreatorID = creator.ID
	return uc.channelRepo.Insert(ctx, channel)
}
