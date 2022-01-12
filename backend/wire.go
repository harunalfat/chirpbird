//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/harunalfat/chirpbird/backend/presentation/persistence"

	"github.com/harunalfat/chirpbird/backend/presentation/web/handlers"
	usecases "github.com/harunalfat/chirpbird/backend/use_cases"
)

func NewApp() (*App, error) {
	wire.Build(
		wire.Struct(new(App), "*"),
		handlers.NewRestHandler,
		handlers.NewWSHandler,
		handlers.NewCentrifugeNode,

		persistence.NewMongoClient,
		persistence.NewMongodbMessageRepository,
		persistence.NewMongodbUserRepository,
		persistence.NewMongodbChannelRepository,

		usecases.NewNodeWrapperImpl,
		usecases.NewMessageUseCase,
		usecases.NewUserUseCase,
		usecases.NewChannelUseCase,
	)
	return &App{}, nil
}
