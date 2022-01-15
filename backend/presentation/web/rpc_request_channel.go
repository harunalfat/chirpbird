package web

import "github.com/harunalfat/chirpbird/backend/entities"

type RPCRequestChannel struct {
	Data entities.Channel `json:"data"`
}
