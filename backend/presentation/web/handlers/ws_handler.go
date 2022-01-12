package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/centrifugal/centrifuge"
	"github.com/google/uuid"
	"github.com/harunalfat/chirpbird/backend/entities"
	usecases "github.com/harunalfat/chirpbird/backend/use_cases"
)

const (
	CHANNEL_GENERAL = "channel_general"
)

var wsHandler http.Handler

type WSHandler struct {
	channelUseCase *usecases.ChannelUseCase
	node           *centrifuge.Node
	userUseCase    *usecases.UserUseCase
}

func NewCentrifugeNode() (*centrifuge.Node, error) {
	return centrifuge.New(centrifuge.DefaultConfig)
}

func NewWSHandler(channelUseCase *usecases.ChannelUseCase, node *centrifuge.Node, userUseCase *usecases.UserUseCase) *WSHandler {
	return &WSHandler{
		channelUseCase,
		node,
		userUseCase,
	}
}

func (handler *WSHandler) auth(rw http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	user, err := handler.userUseCase.CreateIfUsernameNotExist(r.Context(), entities.User{
		Username: username,
	})
	if err != nil {
		log.Printf("Could not insert client\n%s", err)
		jsonError(rw, http.StatusBadRequest, err)
		return
	}

	cred := &centrifuge.Credentials{
		UserID: user.ID.String(),
	}

	// request's context need to be set with Centrifuge credentials
	// or else client handshake will be failed
	newCtx := centrifuge.SetCredentials(r.Context(), cred)

	r = r.WithContext(context.WithValue(newCtx, "username", username))
	wsHandler.ServeHTTP(rw, r)
}

func (handler *WSHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	handler.auth(rw, r)
}

func (handler *WSHandler) newConnectionProcedure(c *centrifuge.Client) error {
	userID := uuid.MustParse(c.UserID())
	log.Printf("Successfully open connection for client [%s]", userID)
	user, err := handler.userUseCase.Fetch(c.Context(), userID)
	if err != nil {
		return err
	}

	for _, channelID := range user.ChannelIDs {
		if err = handler.userUseCase.SubsribeUserConnectionToChannel(c.Context(), userID, channelID); err != nil {
			log.Printf("Cannot subscribe client [%s] connection to channel [%s]", userID, channelID)
		}
	}

	for _, channelID := range user.ChannelIDs {
		log.Printf("PUBLISH AH %s", channelID)
		handler.node.Publish(channelID.String(), []byte("dataa"))
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
		userID := uuid.MustParse(c.UserID())
		channelID := uuid.MustParse(pe.Channel)
		log.Printf("Publish '%s' received from [%s] to channel [%s]", pe.Data, userID, channelID)
		err := handler.channelUseCase.UpdateChannelWithMessage(c.Context(), userID, channelID, string(pe.Data))
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
		cb(centrifuge.RPCReply{}, nil)
		//handler.handleRPC(c, ev, cb)
	})

	c.OnDisconnect(func(de centrifuge.DisconnectEvent) {
		log.Printf("Client [%s] disconnected", c.UserID())
	})
}

func (handler *WSHandler) Init() (err error) {
	handler.node.OnConnect(func(c *centrifuge.Client) {
		if err = handler.newConnectionProcedure(c); err != nil {
			log.Println(err)
			c.Disconnect(centrifuge.DisconnectBadRequest)
		}

		handler.handleClientCallbacks(c)
	})

	if err := handler.node.Run(); err != nil {
		log.Fatalf("Could not start centrifuge node\n%s", err)
	}

	wsHandler = centrifuge.NewWebsocketHandler(handler.node, centrifuge.WebsocketConfig{})
	return nil
}

func (handler *WSHandler) handleRPC(client *centrifuge.Client, ev centrifuge.RPCEvent, cb centrifuge.RPCCallback) {
	switch ev.Method {
	case "channel/invite":
		handler.handleRPCChannelInvite(client, ev, cb)
	default:
		cb(centrifuge.RPCReply{}, centrifuge.ErrorBadRequest)
	}
}

func (handler *WSHandler) handleRPCChannelInvite(client *centrifuge.Client, ev centrifuge.RPCEvent, cb centrifuge.RPCCallback) {
	var payload entities.InvitePayload
	if err := json.Unmarshal(ev.Data, &payload); err != nil {
		cb(centrifuge.RPCReply{}, err)
	}

	//err := handler.userUseCase.EmbedChannelToMultipleUsers(client.Context(), payload.Usernames, payload.ChannelName)
	cb(centrifuge.RPCReply{}, nil)
}
