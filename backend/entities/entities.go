package entities

type User struct {
	Base
	Username   string   `json:"username"`
	ChannelIDs []string `json:"channelIds"`
}

type Channel struct {
	Base
	Name      string `json:"name"`
	CreatorID string `json:"creatorId"`
}

type Message struct {
	Base
	SenderID  string `json:"senderId"`
	ChannelID string `json:"channelId"`
	Text      string `json:"text"`
}
