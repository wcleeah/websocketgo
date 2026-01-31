package websocket_test

import (
	"context"
	"net"
	"sync"
	"testing"
	"time"

	"com.lwc.message_center_server/internal/websocket"
	"github.com/stretchr/testify/assert"
)

func TestSetupSuccess(t *testing.T) {
	client, server := net.Pipe()
	defer client.Close()
	defer server.Close()

	ctx, cancelFunc := context.WithCancel(context.Background())
	ws := websocket.NewWebSocket(ctx, server, 5*time.Second, 0*time.Second)
	ws.Setup()

	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()
		assert.NotPanics(t, func() {
			ws.Read()
		})
	}()

	go func() {
		defer wg.Done()
		assert.NotPanics(t, func() {
			ws.Send(&websocket.TextFrame{Payload: []byte("haha")})
		})
	}()

	cancelFunc()
	wg.Wait()
}

func TestNoSetup(t *testing.T) {
	client, server := net.Pipe()
	defer client.Close()
	defer server.Close()

	ws := websocket.NewWebSocket(context.Background(), server, 5*time.Second, 0*time.Second)

	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()
		assert.Panics(t, func() {
			ws.Read()
		})
	}()

	go func() {
		defer wg.Done()
		assert.Panics(t, func() {
			ws.Send(&websocket.TextFrame{Payload: []byte("haha")})
		})
	}()

	wg.Wait()
}

