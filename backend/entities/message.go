package entities

import "time"

type Message struct {
	Data           interface{} `json:"data"`
	SenderUsername string      `json:"senderUsername"`
	EventName      string      `json:"eventName"`
	Timestamp      time.Time   `json:"timestamp"`
	ChannelID      string      `json:"channelId"`
}
