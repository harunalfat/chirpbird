package mongodb

import (
	"github.com/harunalfat/chirpbird/backend/entities"
	"github.com/harunalfat/chirpbird/backend/helpers"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type channelDTO struct {
	BaseDTO   `bson:",inline"`
	Name      string             `bson:"name,omitempty"`
	CreatorID primitive.ObjectID `bson:"creatorId,omitempty"`
}

func (dto *channelDTO) FromEntity(e entities.Channel) {
	dto.ID = helpers.ObjectIDFromHex(e.ID)
	dto.Name = e.Name
	dto.CreatorID = helpers.ObjectIDFromHex(e.CreatorID)
	dto.CreatedAt = e.CreatedAt
}

func (dto *channelDTO) ToEntity() entities.Channel {
	return entities.Channel{
		Base: entities.Base{
			ID:        dto.ID.Hex(),
			CreatedAt: dto.CreatedAt,
		},
		Name:      dto.Name,
		CreatorID: dto.CreatorID.Hex(),
	}
}

type channelDTOs []channelDTO

func (dtos channelDTOs) FromEntity(es []entities.Channel) {
	dtos = make(channelDTOs, len(es))
	for idx, e := range es {
		dtos[idx].FromEntity(e)
	}

	return
}

func (dtos channelDTOs) ToEntity() []entities.Channel {
	result := make([]entities.Channel, len(dtos))
	for idx, dto := range dtos {
		result[idx] = dto.ToEntity()
	}

	return result
}
