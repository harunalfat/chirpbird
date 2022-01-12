package mongodb

import (
	"github.com/harunalfat/chirpbird/backend/entities"
	"github.com/harunalfat/chirpbird/backend/helpers"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type messageDTO struct {
	BaseDTO   `bson:",inline"`
	SenderID  primitive.ObjectID `bson:"senderId,omitempty"`
	ChannelID primitive.ObjectID `bson:"channelId,omitempty"`
	Text      string             `bson:"text,omitempty"`
}

func (dto *messageDTO) FromEntity(e entities.Message) {
	dto.ID = helpers.ObjectIDFromHex(e.ID)
	dto.CreatedAt = e.CreatedAt
	dto.SenderID = helpers.ObjectIDFromHex(e.SenderID)
	dto.ChannelID = helpers.ObjectIDFromHex(e.ChannelID)
	dto.Text = e.Text
}

func (dto *messageDTO) ToEntity() entities.Message {
	return entities.Message{
		Base: entities.Base{
			ID:        dto.ID.Hex(),
			CreatedAt: dto.CreatedAt,
		},
		SenderID:  dto.SenderID.Hex(),
		ChannelID: dto.ChannelID.Hex(),
		Text:      dto.Text,
	}
}

type messageDTOs []messageDTO

func (dtos messageDTOs) FromEntity(es []entities.Message) {
	dtos = make(messageDTOs, len(es))
	for idx, e := range es {
		dtos[idx].FromEntity(e)
	}

	return
}

func (dtos messageDTOs) ToEntity() []entities.Message {
	result := make([]entities.Message, len(dtos))
	for idx, dto := range dtos {
		result[idx] = dto.ToEntity()
	}

	return result
}
