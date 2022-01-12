//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/harunalfat/chirpbird/backend/presentation/persistence/implementations/mongodb"

	"github.com/harunalfat/chirpbird/backend/presentation/web/handlers"
	usecases "github.com/harunalfat/chirpbird/backend/use_cases/implementations"
)

func NewApp() (*App, error) {
	wire.Build(
		wire.Struct(new(App), "*"),
		handlers.NewRestHandlerImpl,
		handlers.NewCentrifugeHandler,
		handlers.NewCentrifugeNode,

		mongodb.NewMongoClient,
		mongodb.NewMongodbMessageRepository,
		mongodb.NewMongodbUserRepository,
		mongodb.NewMongodbChannelRepository,

		usecases.NewNodeWrapperImpl,
		usecases.NewMessageUseCaseImpl,
		usecases.NewUserUseCaseImpl,
		usecases.NewChannelUseCaseImpl,
	)
	return &App{}, nil
}
