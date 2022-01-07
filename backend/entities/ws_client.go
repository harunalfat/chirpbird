package entities

import "net"

type WSClient struct {
	Username string
	Conn     net.Conn
}
