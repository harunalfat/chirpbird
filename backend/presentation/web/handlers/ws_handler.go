package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/centrifugal/centrifuge"
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
	user, err := handler.userUseCase.Fetch(c.Context(), userID)
	if err != nil {
		return err
	}

	for _, channel := range user.Channels {
		if err = handler.userUseCase.SubsribeUserConnectionToChannel(c.Context(), userID, channel.ID); err != nil {
			log.Printf("Cannot subscribe client [%s] connection to channel [%s]", userID, channel.ID)
		}
	}

	for _, channel := range user.Channels {
		log.Printf("PUBLISH AH %s", channel.ID)
		handler.node.Publish(channel.ID, []byte("dataa"))
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
		userID := c.UserID()
		channelID := pe.Channel

		var message entities.Message
		err := json.Unmarshal(pe.Data, &message)
		if err != nil {
			log.Printf("Invalid message format!\n%s", err)
			return
		}

		err = handler.channelUseCase.UpdateChannelWithMessage(c.Context(), userID, channelID, message.Data.(string))
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
