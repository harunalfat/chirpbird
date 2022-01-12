package usecases

import (
	"context"

	"github.com/harunalfat/chirpbird/backend/entities"
	"github.com/harunalfat/chirpbird/backend/presentation/persistence"
)

type MessageUseCase struct {
	messageRepo persistence.MessageRepository
}

func NewMessageUseCase(messageRepo persistence.MessageRepository) *MessageUseCase {
	return &MessageUseCase{
		messageRepo,
	}
}

func (uc *MessageUseCase) Store(ctx context.Context, msg entities.Message) (entities.Message, error) {
	return uc.messageRepo.Insert(ctx, msg)
}
