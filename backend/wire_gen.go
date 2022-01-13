// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/harunalfat/chirpbird/backend/presentation/persistence"
	"github.com/harunalfat/chirpbird/backend/presentation/web/handlers"
	"github.com/harunalfat/chirpbird/backend/use_cases"
	"go.mongodb.org/mongo-driver/mongo"
)

// Injectors from wire.go:

func NewApp(mongoClient *mongo.Client) (*App, error) {
	channelRepository := persistence.NewMongodbChannelRepository(mongoClient)
	messageRepository := persistence.NewMongodbMessageRepository(mongoClient)
	messageUseCase := usecases.NewMessageUseCase(messageRepository)
	channelUseCase := usecases.NewChannelUseCase(channelRepository, messageUseCase)
	node, err := handlers.NewCentrifugeNode()
	if err != nil {
		return nil, err
	}
	nodeWrapper := usecases.NewNodeWrapperImpl(node)
	userRepository := persistence.NewMongodbUserRepository(mongoClient)
	userUseCase := usecases.NewUserUseCase(channelUseCase, nodeWrapper, userRepository)
	restHandler := handlers.NewRestHandler(channelUseCase, messageUseCase, userUseCase)
	wsHandler := handlers.NewWSHandler(channelUseCase, node, userUseCase)
	app := &App{
		restHandler: restHandler,
		wsHandler:   wsHandler,
	}
	return app, nil
}
