package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/harunalfat/chirpbird/backend/env"
	"github.com/harunalfat/chirpbird/backend/presentation/web/handlers"
)

type App struct {
	restHandler *handlers.RestHandler
	wsHandler   *handlers.WSHandler
	httpSrv     *http.Server
}

func (app *App) Shutdown(ctx context.Context) error {
	err := app.wsHandler.Shutdown(ctx)
	if err != nil {
		log.Printf("Failed to shutdown centrifuge node\n%s")
		return err
	}

	err = app.httpSrv.Shutdown(ctx)
	if err != nil {
		log.Printf("Failed to shutdown http server\n%s", err)
	}

	return err
}

func (app *App) routes(router *gin.Engine) {
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowCredentials: true,
		AllowHeaders:     []string{"Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With"},
	}))

	router.GET("/connection/websocket", gin.WrapF(app.wsHandler.Serve))

	router.GET("/ping", gin.WrapF(func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("pong"))
	}))
	router.POST("/users", gin.WrapF(app.restHandler.RegisterUser))
}

func (app *App) Run() {
	app.wsHandler.Init()
	router := gin.Default()

	app.routes(router)

	app.httpSrv = &http.Server{
		Addr:    fmt.Sprintf(":%s", os.Getenv(env.PORT)),
		Handler: router,
	}

	go func() {
		log.Printf("Starting HTTP server on %s", app.httpSrv.Addr)
		err := app.httpSrv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start HTTP Server\n%s", err)
		}
	}()
}
