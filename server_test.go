package emitter

import (
	"fmt"
	"testing"
)

func TestServer_On(t *testing.T) {
	server := NewServer(nil)
	server.On("hello", func(conn *Conn, handler ...any) {
		fmt.Println("hello")
	})
}
