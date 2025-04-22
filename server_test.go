package main

import (
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func startTestServer() *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(handleConnections))
	go handleMessages()
	return server
}

func TestWebSocketEcho(t *testing.T) {
	ts := startTestServer()
	defer ts.Close()

	weUrl := "ws" + ts.URL[len("http"):] + "/ws"

	ws, _, err := websocket.DefaultDialer.Dial(weUrl, nil)
	if err != nil {
		t.Errorf("WebSocket Dial Error: %v", err)
	}
	defer func(ws *websocket.Conn) {
		err := ws.Close()
		if err != nil {

		}
	}(ws)

	testMessage := "hello world"
	err = ws.WriteMessage(websocket.TextMessage, []byte(testMessage))
	assert.NoError(t, err)

	err = ws.SetReadDeadline(time.Now().Add(2 * time.Second))
	if err != nil {
		return
	}
	_, message, err := ws.ReadMessage()
	assert.NoError(t, err)
	assert.Equal(t, testMessage, string(message))
}
