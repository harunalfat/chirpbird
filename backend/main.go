package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/harunalfat/chirpbird/backend/env"
	"github.com/harunalfat/chirpbird/backend/presentation/persistence"
	"github.com/harunalfat/chirpbird/backend/presentation/web/handlers"
	"github.com/joho/godotenv"
)

type App struct {
	restHandler *handlers.RestHandler
	wsHandler   *handlers.WSHandler
}

func (app *App) run() {
	app.wsHandler.Init()
	router := gin.Default()

	// serve static assets
	router.StaticFS("/statics", http.Dir("../frontend"))

	router.GET("/connection/websocket", gin.WrapH(app.wsHandler))

	router.POST("/users", gin.WrapF(app.restHandler.RegisterUser))
	router.POST("/channels", gin.WrapF(app.restHandler.CreateChannel))
	router.POST("/channels/invite", gin.WrapF(app.restHandler.InviteToChannel))

	router.Run(fmt.Sprintf(":%s", os.Getenv(env.PORT)))
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Failed to initialize environment variables\n%s", err)
	}

	mongoClient := persistence.MongoDBInit()

	app, err := NewApp(mongoClient)
	if err != nil {
		log.Fatalf("Failed to prepare application\n%s", err)
	}

	app.run()
}
