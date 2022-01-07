package controllers

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"

	"github.com/harunalfat/chirpbird/backend/entities"
)

type WSService interface {
	UpgradeHTTP(*http.Request, http.ResponseWriter) (net.Conn, error)
	ReadClientData(io.ReadWriter) (data []byte, opCode byte, err error)
	WriteServerData(writer io.Writer, opCode byte, data []byte) error
	ProcessData(client entities.WSClient, data []byte, opCode byte) error
	SubscribeToRegisteredChannels(client entities.WSClient) error
}

type WSServer struct {
	service WSService
}

var clients = make(map[string]*entities.WSClient)

func NewWSServer(service WSService) *WSServer {
	return &WSServer{
		service: service,
	}
}

func (wss WSServer) handleClient(client entities.WSClient) {
	defer func() {
		client.Conn.Close()
		delete(clients, client.Username)
	}()

	for {
		input, opCode, err := wss.service.ReadClientData(client.Conn)
		if err != nil {
			fmt.Printf("Error reading client data, [%v]", err)
			break
		}

		if err = wss.service.ProcessData(client, input, opCode); err != nil {
			log.Printf("Error when processing client data\n%s", err)
		}
	}
	log.Println("Connection closed")
	log.Println(clients)
	log.Println(wss)
}

func (wss WSServer) ServeWS(w http.ResponseWriter, r *http.Request) {
	conn, err := wss.service.UpgradeHTTP(r, w)
	if err != nil {
		log.Println(err)
		return
	}

	username := r.FormValue("username")
	if v, exist := clients[username]; exist {
		v.Conn.Close()
		delete(clients, username)
		log.Printf("Delete old connection for username [%s]", username)
	}

	client := entities.WSClient{
		Username: username,
		Conn:     conn,
	}
	clients[username] = &client
	log.Printf("Successfully open WS connection")

	go wss.handleClient(client)
}
