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

func TestCreateNewUser(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	userRepo := mocks.NewMockUserRepository(mockCtrl)

	useCase := usecases.NewUserUseCase(
		&usecases.ChannelUseCase{},
		usecases.NodeWrapperImpl{},
		userRepo,
	)

	userRepo.EXPECT().
		FetchByUsername(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, username string) (entities.User, error) {
			newUser := entities.User{
				Username: username,
				Base: entities.Base{
					ID: uuid.NewString(),
				},
			}

			return newUser, nil
		}).MaxTimes(1)

	userRepo.EXPECT().
		Insert(gomock.Any(), gomock.Any()).
		Return(entities.User{}, nil).MaxTimes(0)

	user := entities.User{
		Username: "Myusername",
	}

	res, _ := useCase.CreateIfUsernameNotExist(context.Background(), user)
	if res.Username != user.Username {
		t.Error("Username is not the same")
	}

	if res.ID == "" {
		t.Error("ID not assigned")
	}
}

func TestCreateExistingUser(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	userRepo := mocks.NewMockUserRepository(mockCtrl)

	useCase := usecases.NewUserUseCase(
		&usecases.ChannelUseCase{},
		usecases.NodeWrapperImpl{},
		userRepo,
	)

	userRepo.EXPECT().
		FetchByUsername(gomock.Any(), gomock.Any()).
		Return(entities.User{
			Username: "",
		}, nil).MaxTimes(1)

	user := entities.User{
		Username: "Myusername",
	}

	mockedUser := entities.User{
		Username: user.Username,
		Base: entities.Base{
			ID: uuid.NewString(),
		},
	}

	userRepo.EXPECT().
		Insert(gomock.Any(), gomock.Any()).
		Return(mockedUser, nil).MaxTimes(1)

	res, _ := useCase.CreateIfUsernameNotExist(context.Background(), user)
	if res.Username != user.Username {
		t.Error("Username is not the same")
	}

	if res.ID == "" {
		t.Error("ID not assigned")
	}
}
