package entities

import "github.com/google/uuid"

type InvitePayload struct {
	UserIDs   []uuid.UUID `json:"userIds"`
	ChannelID uuid.UUID   `json:"channelID"`
}
