package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/centrifugal/centrifuge"
	"github.com/harunalfat/chirpbird/backend/entities"
	"github.com/harunalfat/chirpbird/backend/presentation/web"
	usecases "github.com/harunalfat/chirpbird/backend/use_cases"
)

const (
	CHANNEL_GENERAL = "channel_general"

	CREATE_CHANNEL = "create_channel"
	FETCH_MESSAGE  = "fetch_message"
	SEARCH_USERS   = "search_users"
)

var wsHandler http.Handler

type WSHandler struct {
	channelUseCase *usecases.ChannelUseCase
	messageUseCase *usecases.MessageUseCase
	node           *centrifuge.Node
	userUseCase    *usecases.UserUseCase
}

func NewCentrifugeNode() (*centrifuge.Node, error) {
	return centrifuge.New(centrifuge.DefaultConfig)
}

func NewWSHandler(channelUseCase *usecases.ChannelUseCase, messageUseCase *usecases.MessageUseCase, node *centrifuge.Node, userUseCase *usecases.UserUseCase) *WSHandler {
	return &WSHandler{
		channelUseCase,
		messageUseCase,
		node,
		userUseCase,
	}
}

func (handler *WSHandler) auth(rw http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userId")
	cred := &centrifuge.Credentials{
		UserID: userID,
	}

	// request's context need to be set with Centrifuge credentials
	// or else client handshake will be failed
	newCtx := centrifuge.SetCredentials(r.Context(), cred)
	wsHandler.ServeHTTP(rw, r.WithContext(newCtx))
}

func (handler *WSHandler) Serve(rw http.ResponseWriter, r *http.Request) {
	handler.auth(rw, r)
}

func (handler *WSHandler) newConnectionProcedure(c *centrifuge.Client) error {
	userID := c.UserID()

	log.Printf("Successfully open connection for client [%s]", userID)
	_, err := handler.userUseCase.Fetch(c.Context(), userID)
	if err != nil {
		return err
	}

	return nil
}

func (handler *WSHandler) handleClientCallbacks(c *centrifuge.Client) {
	c.OnSubscribe(func(se centrifuge.SubscribeEvent, sc centrifuge.SubscribeCallback) {
		log.Printf("Subscribe from user: %s, data: %s, channel: %s", c.UserID(), string(se.Data), se.Channel)
		sc(centrifuge.SubscribeReply{
			Options: centrifuge.SubscribeOptions{
				Presence:  true,
				Recover:   true,
				JoinLeave: true,
				Data:      []byte(`{"msg": "welcome"}`),
			},
		}, nil)
	})

	c.OnMessage(func(me centrifuge.MessageEvent) {
		log.Printf("Received echo message [%s] from [%s]", me.Data, c.UserID())
	})

	c.OnPublish(func(pe centrifuge.PublishEvent, pc centrifuge.PublishCallback) {
		log.Printf("Publish '%s' received from [%s] to channel [%s]", pe.Data, c.UserID(), pe.Channel)

		var message entities.Message
		err := json.Unmarshal(pe.Data, &message)
		if err != nil {
			log.Printf("Invalid message format!\n%s", err)
			return
		}

		err = handler.channelUseCase.UpdateChannelWithMessage(c.Context(), message)
		if err != nil {
			log.Printf("Failed to process message!\n%s", err)
		}

		pc(centrifuge.PublishReply{
			Options: centrifuge.PublishOptions{
				ClientInfo: &centrifuge.ClientInfo{
					ClientID: c.ID(),
					UserID:   c.UserID(),
				},
			},
		}, nil)
	})

	c.OnPresence(func(pe centrifuge.PresenceEvent, pc centrifuge.PresenceCallback) {
		log.Printf("Presence check from [%s], channel [%s]", c.UserID(), pe.Channel)
		pc(centrifuge.PresenceReply{}, nil)
	})

	c.OnRPC(func(ev centrifuge.RPCEvent, cb centrifuge.RPCCallback) {
		log.Printf("RPC from user: %s, data: %s, method: %s", c.UserID(), string(ev.Data), ev.Method)
		var result []byte
		var err error
		switch ev.Method {
		case FETCH_MESSAGE:
			result, err = handler.FetchMessage(c.Context(), ev.Data)
		case SEARCH_USERS:
			result, err = handler.SearchUsersByName(c.Context(), ev.Data)
		case CREATE_CHANNEL:
			result, err = handler.CreateChannelIfNotExist(c.Context(), ev.Data, c.UserID())
		}

		if err != nil {
			log.Printf("RPC error\n%s", err)
		}

		cb(centrifuge.RPCReply{
			Data: result,
		}, err)
	})

	c.OnDisconnect(func(de centrifuge.DisconnectEvent) {
		log.Printf("Client [%s] disconnected", c.UserID())
	})
}

func (handler *WSHandler) FetchMessage(ctx context.Context, input []byte) (result []byte, err error) {
	var payload web.Response
	err = json.Unmarshal(input, &payload)
	if err != nil {
		return
	}

	messages, err := handler.messageUseCase.FetchAllMessagesByChannel(ctx, payload.Data.(string))
	if err != nil {
		return
	}

	result, err = json.Marshal(messages)
	return
}

func (handler *WSHandler) SearchUsersByName(ctx context.Context, input []byte) (result []byte, err error) {
	var payload web.Response
	err = json.Unmarshal(input, &payload)
	if err != nil {
		return
	}

	users, err := handler.userUseCase.SearchByUsername(ctx, payload.Data.(string))
	if err != nil {
		return
	}

	result, err = json.Marshal(users)
	return
}

func (handler *WSHandler) CreateChannelIfNotExist(ctx context.Context, input []byte, creatorID string) (result []byte, err error) {
	var payload web.Response
	err = json.Unmarshal(input, &payload)
	if err != nil {
		return
	}

	creator, err := handler.userUseCase.Fetch(ctx, creatorID)
	if err != nil {
		return
	}

	channel, err := handler.channelUseCase.CreateIfNameNotExist(ctx, payload.Data.(entities.Channel), creator)
	if err != nil {
		return
	}

	result, err = json.Marshal(channel)
	return
}

func (handler *WSHandler) Init() (err error) {
	handler.node.OnConnecting(func(c context.Context, ce centrifuge.ConnectEvent) (centrifuge.ConnectReply, error) {
		fmt.Println(ce.Token)
		fmt.Println("AMAAN")

		return centrifuge.ConnectReply{}, nil
	})
	handler.node.OnConnect(func(c *centrifuge.Client) {
		if err = handler.newConnectionProcedure(c); err != nil {
			log.Println(err)
			c.Disconnect(centrifuge.DisconnectBadRequest)
		}

		err = c.Subscribe("NEW_PRIVATE_CHANNEL")
		if err != nil {
			log.Printf("User cannot subscribe to notification channel\n%s", err)
		}

		handler.handleClientCallbacks(c)
	})

	if err := handler.node.Run(); err != nil {
		log.Fatalf("Could not start centrifuge node\n%s", err)
	}

	wsHandler = centrifuge.NewWebsocketHandler(handler.node, centrifuge.WebsocketConfig{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	})
	return nil
}
