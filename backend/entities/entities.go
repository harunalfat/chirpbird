package entities

import "github.com/google/uuid"

type User struct {
	Base       `bson:",inline"`
	Username   string      `json:"username,omitempty" bson:"username,omitempty"`
	ChannelIDs []uuid.UUID `json:"channelIds,omitempty" bson:"channelIds,omitempty"`
}

type Channel struct {
	Base           `bson:",inline"`
	Name           string    `json:"name,omitempty" bson:"name,omitempty"`
	CreatorID      uuid.UUID `json:"creatorId,omitempty" bson:"creatorId,omitempty"`
	IsPrivate      bool      `json:"isPrivate" bson:"isPrivate"`
	HashIdentifier string    `json:"hashIdentifier,omitempty" bson:"hashIdentifier,omitempty"`
}

type Message struct {
	Base      `bson:",inline"`
	Sender    User      `json:"senderId,omitempty"`
	ChannelID uuid.UUID `json:"channelId,omitempty"`
	Data      string    `json:"text,omitempty"`
}
