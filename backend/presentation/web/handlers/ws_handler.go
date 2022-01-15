package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/centrifugal/centrifuge"
	"github.com/harunalfat/chirpbird/backend/entities"
	"github.com/harunalfat/chirpbird/backend/env"
	"github.com/harunalfat/chirpbird/backend/presentation/web"
	usecases "github.com/harunalfat/chirpbird/backend/use_cases"
)

const (
	NEW_PRIVATE_CHANNEL = "NEW_PRIVATE_CHANNEL"

	CREATE_CHANNEL = "create_channel"
	FETCH_MESSAGE  = "fetch_message"
	SEARCH_USERS   = "search_users"

	SERVER_NOTIFICATION = "SERVER_NOTIFICATION"
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
	err := c.Send([]byte(`{"hello": "world"}`))
	if err != nil {
		log.Println(err)
	}
	err = c.Subscribe(SERVER_NOTIFICATION)
	if err != nil {
		log.Printf("Failed to add client for server side notif\n%s", err)
	}
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
		log.Printf("Received echo message [%s] from [%s]", string(me.Data), c.UserID())

		var payload entities.Channel
		err = json.Unmarshal(me.Data, &payload)
		if err != nil {
			return
		}
		log.Println(payload)

		if !payload.IsPrivate {
			return
		}

		_, err := handler.node.Publish(SERVER_NOTIFICATION, me.Data)
		if err != nil {
			log.Printf("Failed to publish server notification\n%s", err)
		}
	})

	c.OnPublish(func(pe centrifuge.PublishEvent, pc centrifuge.PublishCallback) {
		log.Printf("Publish '%s' received from [%s] to channel [%s]", pe.Data, c.UserID(), pe.Channel)

		var message entities.Message
		err := json.Unmarshal(pe.Data, &message)
		if err != nil {
			log.Printf("Invalid message format!\n%s", err)
			return
		}

		if pe.Channel != NEW_PRIVATE_CHANNEL {
			err = handler.channelUseCase.UpdateChannelWithMessage(c.Context(), message)
			if err != nil {
				log.Printf("Failed to process message!\n%s", err)
			}
		}

		pc(centrifuge.PublishReply{
			Options: centrifuge.PublishOptions{
				ClientInfo: &centrifuge.ClientInfo{
					ClientID: c.ID(),
					UserID:   c.UserID(),
				},
			},
		}, err)
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
	var payload web.RPCRequestString
	log.Println(string(input))
	err = json.Unmarshal(input, &payload)
	if err != nil {
		return
	}

	messages, err := handler.messageUseCase.FetchAllMessagesByChannel(ctx, payload.Data)
	if err != nil {
		return
	}

	result, err = json.Marshal(messages)
	return
}

func (handler *WSHandler) SearchUsersByName(ctx context.Context, input []byte) (result []byte, err error) {
	var payload web.RPCRequestString
	err = json.Unmarshal(input, &payload)
	if err != nil {
		return
	}

	users, err := handler.userUseCase.SearchByUsername(ctx, payload.Data)
	if err != nil {
		return
	}

	result, err = json.Marshal(users)
	return
}

func (handler *WSHandler) CreateChannelIfNotExist(ctx context.Context, input []byte, creatorID string) (result []byte, err error) {
	var payload web.RPCRequestChannel
	log.Println(payload)
	log.Println(string(input))
	err = json.Unmarshal(input, &payload)
	if err != nil {
		return
	}

	creator, err := handler.userUseCase.Fetch(ctx, creatorID)
	if err != nil {
		return
	}

	channel, err := handler.channelUseCase.CreateIfNameNotExist(ctx, payload.Data, creator)
	if err != nil {
		return
	}

	_, err = handler.userUseCase.EmbedChannelIfNotExist(ctx, creator, channel)
	if err != nil {
		return
	}

	for _, p := range channel.Participants {
		var user entities.User
		user, err = handler.userUseCase.Fetch(ctx, p.ID)
		if err != nil {
			return
		}

		user, err = handler.userUseCase.EmbedChannelIfNotExist(ctx, user, channel)
		if err != nil {
			return
		}
	}

	result, err = json.Marshal(channel)
	return
}

func (handler *WSHandler) setupRedisAdapter() (err error) {
	sentinels := strings.Split(os.Getenv(env.REDIS_SENTINEL_ADDRESSES), ",")
	log.Println(sentinels[0])
	redisShard, err := centrifuge.NewRedisShard(handler.node, centrifuge.RedisShardConfig{
		Address:            os.Getenv(env.REDIS_ADDRESS),
		SentinelAddresses:  sentinels,
		Password:           os.Getenv(env.REDIS_PASSWORD),
		SentinelMasterName: os.Getenv(env.REDIS_MASTER_NAME),
	})

	if err != nil {
		log.Fatalf("Cannot create redis shard instance\n%s", err)
		return
	}

	redisBroker, err := centrifuge.NewRedisBroker(handler.node, centrifuge.RedisBrokerConfig{
		Shards: []*centrifuge.RedisShard{redisShard},
	})

	if err != nil {
		log.Fatalf("Cannot create redis broker instance\n%s", err)
		return
	}

	redisPresenceManager, err := centrifuge.NewRedisPresenceManager(handler.node, centrifuge.RedisPresenceManagerConfig{
		Shards: []*centrifuge.RedisShard{redisShard},
	})

	if err != nil {
		log.Fatalf("Cannot create redis presence manager instance\n%s", err)
		return
	}

	handler.node.SetBroker(redisBroker)
	handler.node.SetPresenceManager(redisPresenceManager)
	return nil
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

	if err = handler.setupRedisAdapter(); err != nil {
		log.Fatalf("Cannot use redis as adapter\n%s", err)
	}

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
