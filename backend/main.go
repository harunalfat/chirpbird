package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
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

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowCredentials: true,
		AllowHeaders:     []string{"Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With"},
	}))

	// serve static assets
	//router.StaticFS("/statics", http.Dir("../frontend/dist/client"))

	router.GET("/connection/websocket", gin.WrapF(app.wsHandler.Serve))

	router.GET("/ping", gin.WrapF(func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("pong"))
	}))
	router.POST("/users", gin.WrapF(app.restHandler.RegisterUser))
	router.POST("/channels", gin.WrapF(app.restHandler.CreateChannel))

	router.Run(fmt.Sprintf(":%s", os.Getenv(env.PORT)))
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Cannot read environment from file\n%s", err)
	}

	mongoClient := persistence.MongoDBInit()

	app, err := NewApp(mongoClient)
	if err != nil {
		log.Fatalf("Failed to prepare application\n%s", err)
	}

	app.run()
}
