package emitter

import "github.com/gorilla/websocket"

type Conn struct {
	*websocket.Conn
}

func newConn(conn *websocket.Conn) *Conn {
	return &Conn{
		conn,
	}
}
