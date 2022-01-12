package implementations

import (
	"context"

	"github.com/harunalfat/chirpbird/backend/entities"
	"github.com/harunalfat/chirpbird/backend/presentation/persistence"
	usecases "github.com/harunalfat/chirpbird/backend/use_cases"
)

type MessageUseCaseImpl struct {
	messageRepo persistence.MessageRepository
}

func NewMessageUseCaseImpl(messageRepo persistence.MessageRepository) usecases.MessageUseCase {
	return &MessageUseCaseImpl{
		messageRepo,
	}
}

func (uc *MessageUseCaseImpl) Store(ctx context.Context, msg entities.Message) (entities.Message, error) {
	return uc.messageRepo.Insert(ctx, msg)
}
