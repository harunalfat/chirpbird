package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/harunalfat/chirpbird/backend/adapters/persistence"
	"github.com/harunalfat/chirpbird/backend/entities"
	"github.com/harunalfat/chirpbird/backend/env"
)

type RedisChannelRepository struct {
	client *redis.ClusterClient
}

func NewRedisClusterClient() (*redis.ClusterClient, error) {
	addressEnv := os.Getenv(env.EnvVarKey.RedisAddresses)
	log.Println(addressEnv)
	addresses := strings.Split(addressEnv, ",")
	c := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: addresses,
	})

	if err := c.Ping(context.TODO()).Err(); err != nil {
		return nil, err
	}

	return c, nil
}

func getRegisteredChannelsKey(username string) string {
	return fmt.Sprintf("%s:registeredChannels", username)
}

func NewRedisChannelRepository(client *redis.ClusterClient) persistence.ChannelRepository {
	return &RedisChannelRepository{
		client: client,
	}
}

func (repo *RedisChannelRepository) FetchRegisteredChannels(username string) (channels []entities.Channel, err error) {
	channelNames, err := repo.client.SMembers(context.TODO(), getRegisteredChannelsKey(username)).Result()
	if err != nil {
		return
	}

	for _, channelName := range channelNames {
		channels = append(channels, entities.Channel{
			Name: channelName,
		})
	}

	return
}

func (repo *RedisChannelRepository) RegisterChannel(channelName string, username string) (err error) {
	err = repo.client.SAdd(context.TODO(), getRegisteredChannelsKey(username)).Err()
	return
}

func (repo *RedisChannelRepository) Subscribe(subscribeHandler func(entities.Message), channels []entities.Channel) (err error) {
	channelNames := make([]string, len(channels))
	for idx, v := range channels {
		channelNames[idx] = v.Name
	}

	subscriber := repo.client.Subscribe(context.TODO(), channelNames...)
	for {
		msg, err := subscriber.ReceiveMessage(context.TODO())
		if err != nil {
			break
		}

		var message entities.Message
		if err := json.Unmarshal([]byte(msg.Payload), &message); err != nil {
			log.Printf("message cannot be unmarshall\n%s", err)
		}

		subscribeHandler(message)
	}

	return
}

func (repo *RedisChannelRepository) Publish(message entities.Message) (err error) {
	payload, err := json.Marshal(message)
	if err != nil {
		return
	}

	err = repo.client.Publish(context.TODO(), message.ChannelID, payload).Err()
	return
}
