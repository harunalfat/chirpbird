//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/harunalfat/chirpbird/backend/presentation/persistence"

	"github.com/harunalfat/chirpbird/backend/presentation/web/handlers"
	usecases "github.com/harunalfat/chirpbird/backend/use_cases"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

func NewApp(mongoClient *mongo.Client) (*App, error) {
	wire.Build(
		wire.Struct(new(App), "*"),
		wire.Struct(new(http.Server)),
		handlers.NewRestHandler,
		handlers.NewWSHandler,
		handlers.NewCentrifugeNode,

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
