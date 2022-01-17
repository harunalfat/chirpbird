package usecases_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/harunalfat/chirpbird/backend/entities"
	"github.com/harunalfat/chirpbird/backend/presentation/persistence/mocks"
	usecases "github.com/harunalfat/chirpbird/backend/use_cases"
)

func TestStoreMessage(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	messageRepo := mocks.NewMockMessageRepository(mockCtrl)
	messageUseCase := usecases.NewMessageUseCase(messageRepo)

	message := entities.Message{
		Data: "something",
	}

	messageRepo.EXPECT().
		Insert(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, data entities.Message) (entities.Message, error) {
			data.ID = uuid.NewString()
			return data, nil
		})

	res, _ := messageUseCase.Store(context.Background(), message)
	if res.ID == "" {
		t.Error("Should be assigned with UUID")
	}
}
