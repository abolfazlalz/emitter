package emitter

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
	"time"
)

type ServerOption struct {
	ReadBufferSize    int
	WriteBufferSize   int
	HandshakeTimeout  time.Duration
	CheckOrigin       func(r *http.Request) bool
	EnableCompression bool
}

type ListenerEvent func(conn *Conn, handler ...any)

type Server struct {
	mu        sync.RWMutex
	conn      []*Conn
	listeners map[string][]ListenerEvent
	upgrader  websocket.Upgrader
	msgChan   chan []byte
	quitChan  chan *Conn
}

func NewServer(option *ServerOption) *Server {
	var opt ServerOption
	if option != nil {
		opt = *option
	}

	return &Server{
		listeners: make(map[string][]ListenerEvent),
		upgrader: websocket.Upgrader{
			HandshakeTimeout:  opt.HandshakeTimeout,
			ReadBufferSize:    opt.ReadBufferSize,
			WriteBufferSize:   opt.WriteBufferSize,
			WriteBufferPool:   nil,
			Subprotocols:      nil,
			Error:             nil,
			CheckOrigin:       opt.CheckOrigin,
			EnableCompression: opt.EnableCompression,
		},
	}
}

func (s *Server) handleReceiveMessage(conn *Conn) {
	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			s.quitChan <- conn
			return
		}

		if msgType == websocket.BinaryMessage {
			s.msgChan <- msg
		}
	}
}

func (s *Server) Handler(w http.ResponseWriter, r *http.Request, header http.Header) {
	client, err := s.upgrader.Upgrade(w, r, header)
	if err != nil {
		panic(err)
	}

	conn := newConn(client)

	go s.handleReceiveMessage(conn)

	s.conn = append(s.conn, conn)
}

func (s *Server) StartListen() {
	for {
		msg := <-s.msgChan
		fmt.Println(msg)
	}
}

func (s *Server) On(key string, event ListenerEvent) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.listeners[key]; ok {
		s.listeners[key] = make([]ListenerEvent, 0)
	}
	s.listeners[key] = append(s.listeners[key], event)
}
