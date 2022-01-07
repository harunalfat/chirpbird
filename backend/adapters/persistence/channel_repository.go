package persistence

import "github.com/harunalfat/chirpbird/backend/entities"

type ChannelRepository interface {
	FetchRegisteredChannels(username string) ([]entities.Channel, error)
	RegisterChannel(channelName string, username string) error
	Subscribe(subscribeHandler func(entities.Message), channels []entities.Channel) error
	Publish(entities.Message) error
}
