package web

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/harunalfat/chirpbird/backend/adapters/web/controllers"
)

type WebServer struct {
	wsServer *controllers.WSServer
}

func NewWebServer(wsServer *controllers.WSServer) *WebServer {
	return &WebServer{
		wsServer: wsServer,
	}
}

func setupRouter(webServer *WebServer) error {
	router := gin.Default()
	router.GET("/ws", func(c *gin.Context) {
		webServer.wsServer.ServeWS(c.Writer, c.Request)
	})

	return router.Run()
}

func Run(webServer *WebServer) (err error) {
	log.Printf("Starting web server")
	return setupRouter(webServer)
}
