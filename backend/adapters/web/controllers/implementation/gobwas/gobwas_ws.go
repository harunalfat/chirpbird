package gobwas

import (
	"encoding/json"
	"io"
	"net"
	"net/http"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/harunalfat/chirpbird/backend/adapters/web/controllers"
	"github.com/harunalfat/chirpbird/backend/entities"
)

type Handler interface {
	HandleEvent(entities.Message, entities.WSClient) ([]byte, byte, error)
	SubscribeToRegisteredChannels(entities.WSClient) error
}

type GobwasWSService struct {
	handler Handler
}

func NewGobwasWSService(handler Handler) controllers.WSService {
	return &GobwasWSService{
		handler: handler,
	}
}

func (gws *GobwasWSService) UpgradeHTTP(r *http.Request, rw http.ResponseWriter) (net.Conn, error) {
	conn, _, _, err := ws.UpgradeHTTP(r, rw)
	if err != nil {
		return nil, err
	}

	return conn, err
}

func (gws *GobwasWSService) ReadClientData(rw io.ReadWriter) (data []byte, opCode byte, err error) {
	data, code, err := wsutil.ReadClientData(rw)
	return data, byte(code), err
}

func (gws *GobwasWSService) WriteServerData(writer io.Writer, opCode byte, data []byte) error {
	return wsutil.WriteServerMessage(writer, ws.OpCode(opCode), data)
}

func (gws *GobwasWSService) ProcessData(client entities.WSClient, input []byte, opCode byte) (err error) {
	var message entities.Message
	if err = json.Unmarshal(input, &message); err != nil {
		return
	}

	gws.handler.HandleEvent(message, client)
	return
}

func (gws *GobwasWSService) SubscribeToRegisteredChannels(client entities.WSClient) (err error) {
	err = gws.handler.SubscribeToRegisteredChannels(client)
	return
}
