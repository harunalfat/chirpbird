package usecases

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"strings"
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

func (uc *ChannelUseCase) UpdateChannelWithMessage(ctx context.Context, message entities.Message) error {
	if _, err := uc.messageUseCase.Store(ctx, message); err != nil {
		return err
	}

	channel, err := uc.channelRepo.Fetch(ctx, message.ChannelID)
	if err != nil {
		return err
	}

	_, err = uc.channelRepo.Update(ctx, message.ChannelID, channel)
	return err
}

func (uc *ChannelUseCase) CreateIfNameNotExist(ctx context.Context, channel entities.Channel, creator entities.User) (result entities.Channel, err error) {
	result, err = uc.channelRepo.FetchByName(ctx, channel.Name)
	if err != nil && err != mongo.ErrNoDocuments {
		return entities.Channel{}, err
	}

	if result.IsPrivate {
		isParticipant := false
		for _, participant := range channel.Participants {
			if participant.ID == creator.ID {
				isParticipant = true
			}
		}

		if !isParticipant {
			err = errors.New("cannot recreate private channel")
			return
		}
	}

	if result.Name == "" {
		channel.CreatorID = creator.ID
		channel.CreatedAt = time.Now()
		hash := sha256.Sum256([]byte(strings.ToLower(channel.Name)))
		channel.HashIdentifier = fmt.Sprintf("%x", hash[:])

		result, err = uc.channelRepo.Insert(ctx, channel)
		if err != nil {
			return entities.Channel{}, err
		}
	}

	return result, err
}
