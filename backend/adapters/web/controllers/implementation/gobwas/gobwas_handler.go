package gobwas

import (
	"encoding/json"
	"log"

	"github.com/gobwas/ws/wsutil"
	"github.com/harunalfat/chirpbird/backend/adapters/persistence"
	"github.com/harunalfat/chirpbird/backend/entities"
)

const (
	EVENT_CHAT               = "event_chat"
	EVENT_INITIAL_CONNECTION = "event_initial_connection"
)

type GobwasHandler struct {
	channelRepository persistence.ChannelRepository
}

func NewGobwasHandler(channelRepository persistence.ChannelRepository) Handler {
	return &GobwasHandler{
		channelRepository: channelRepository,
	}
}

func (gw *GobwasHandler) SubscribeToRegisteredChannels(client entities.WSClient) (err error) {
	channels, err := gw.channelRepository.FetchRegisteredChannels(client.Username)
	if err != nil {
		return
	}

	err = gw.channelRepository.Subscribe(func(message entities.Message) {
		chatPayload, err := json.Marshal(message)
		if err != nil {
			log.Printf("Cannot marshal subscribed chat payload\n%s", err)
			return
		}

		err = wsutil.WriteServerMessage(client.Conn, 0x1, chatPayload)
		if err != nil {
			log.Printf("Failed to send subscribed message\n%s", err)
		}

	}, channels)
	return
}

func (gw *GobwasHandler) HandleEvent(message entities.Message, client entities.WSClient) (response []byte, respOpCode byte, err error) {
	switch message.EventName {
	case EVENT_INITIAL_CONNECTION:
		gw.handleInitialConnection()
	case EVENT_CHAT:
		log.Println("Handling income chat")
		gw.handleChat(message, client)
	}
	return
}

func (gw *GobwasHandler) handleInitialConnection() {}
func (gw *GobwasHandler) handleChat(message entities.Message, client entities.WSClient) {
	gw.channelRepository.Publish(message)
}
