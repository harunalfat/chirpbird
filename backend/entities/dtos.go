package entities

type InvitePayload struct {
	UserIDs   []string `json:"userIds"`
	ChannelID string   `json:"channelID"`
}
