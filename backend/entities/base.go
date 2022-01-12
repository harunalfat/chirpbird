package entities

import (
	"time"

	"github.com/google/uuid"
)

type Base struct {
	ID        uuid.UUID `json:"id" bson:"id"`
	CreatedAt time.Time `json:"createdAt,omitempty" bson:"id,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty" bson:"id,omitempty"`
}
