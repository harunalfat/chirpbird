package entities

import (
	"time"
)

type Entities interface {
	GetLength() int
	GetID(position int) string
}

type Base struct {
	ID        string    `json:"id" bson:"id"`
	CreatedAt time.Time `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}
