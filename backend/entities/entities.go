package entities

type User struct {
	Base       `bson:",inline"`
	Username   string   `json:"username,omitempty" bson:"username,omitempty"`
	ChannelIDs []string `json:"channelIds,omitempty" bson:"channelIds,omitempty"`
}

type Channel struct {
	Base      `bson:",inline"`
	Name      string `json:"name,omitempty"`
	CreatorID string `json:"creatorId,omitempty"`
}

type Message struct {
	Base      `bson:",inline"`
	SenderID  string `json:"senderId,omitempty"`
	ChannelID string `json:"channelId,omitempty"`
	Data      string `json:"text,omitempty"`
}
