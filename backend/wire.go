//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/harunalfat/chirpbird/backend/adapters/persistence/redis"
	"github.com/harunalfat/chirpbird/backend/adapters/web"
	"github.com/harunalfat/chirpbird/backend/adapters/web/controllers"
	"github.com/harunalfat/chirpbird/backend/adapters/web/controllers/implementation/gobwas"
)

func Initialize() (*web.WebServer, error) {
	wire.Build(web.NewWebServer, controllers.NewWSServer, gobwas.NewGobwasWSService, gobwas.NewGobwasHandler, redis.NewRedisChannelRepository, redis.NewRedisClusterClient)
	return &web.WebServer{}, nil
}
