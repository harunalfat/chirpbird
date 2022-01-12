package mongodb

import (
	"github.com/harunalfat/chirpbird/backend/entities"
	"github.com/harunalfat/chirpbird/backend/helpers"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type userDTO struct {
	BaseDTO    `bson:",inline"`
	Username   string               `bson:"username,omitempty"`
	ChannelIDs []primitive.ObjectID `bson:"channelids,omitempty"`
}

func (dto *userDTO) FromEntity(e entities.User) {
	dto.ID = helpers.ObjectIDFromHex(e.ID)
	dto.Username = e.Username
	dto.ChannelIDs = helpers.ObjectIDsFromHexes(e.ChannelIDs)
	dto.CreatedAt = e.CreatedAt
}

func (dto *userDTO) ToEntity() entities.User {
	return entities.User{
		Base: entities.Base{
			ID:        dto.ID.Hex(),
			CreatedAt: dto.CreatedAt,
		},
		Username:   dto.Username,
		ChannelIDs: helpers.HexesFromObjectIDs(dto.ChannelIDs),
	}
}

type userDTOs []userDTO

func (dtos userDTOs) FromEntity(es []entities.User) {
	dtos = make(userDTOs, len(es))
	for idx, e := range es {
		dtos[idx].FromEntity(e)
	}
}

func (dtos userDTOs) ToEntity() []entities.User {
	result := make([]entities.User, len(dtos))
	for idx, dto := range dtos {
		result[idx] = dto.ToEntity()
	}

	return result
}
