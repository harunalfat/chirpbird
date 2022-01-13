package entities

type User struct {
	Base     `bson:",inline"`
	Username string   `json:"username,omitempty" bson:"username,omitempty"`
	Channels Channels `json:"channels,omitempty" bson:"channels,omitempty"`
}

type Channel struct {
	Base           `bson:",inline"`
	Name           string `json:"name,omitempty" bson:"name,omitempty"`
	CreatorID      string `json:"creatorId,omitempty" bson:"creatorId,omitempty"`
	IsPrivate      bool   `json:"isPrivate" bson:"isPrivate"`
	Participants   []User `json:"participants,omitempty" bson:"participants,omitempty"`
	HashIdentifier string `json:"hashIdentifier,omitempty" bson:"hashIdentifier,omitempty"`
}

type Channels []Channel

func (c Channels) GetLength() int {
	return len(c)
}

func (c Channels) GetID(pos int) string {
	return c[pos].ID
}

type Message struct {
	Base      `bson:",inline"`
	Sender    User        `json:"sender,omitempty" bson:"sender,omitempty"`
	ChannelID string      `json:"channelId,omitempty" bson:"channelId,omitempty"`
	Data      interface{} `json:"data,omitempty" bson:"data,omitempty"`
}
