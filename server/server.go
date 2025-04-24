package server

import (
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"sync"
)

var upgrader = websocket.Upgrader{}

type Client struct {
	conn *websocket.Conn
	send chan []byte
}

var clients = make(map[*Client]bool)
var broadcast = make(chan []byte)
var mu sync.Mutex

func HandleConnections(c echo.Context) error {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer ws.Close()

	client := &Client{conn: ws, send: make(chan []byte, 256)}
	mu.Lock()
	clients[client] = true
	mu.Unlock()

	defer func() {
		mu.Lock()
		delete(clients, client)
		mu.Unlock()
		close(client.send)
	}()

	go func(c *Client) {
		for msg := range c.send {
			if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Println("write error:", err)
				break
			}
		}
	}(client)

	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			break
		}
		broadcast <- message
	}
	return nil
}

func HandleMessages() {
	for {
		message := <-broadcast
		mu.Lock()
		for client := range clients {
			select {
			case client.send <- message:
			default:
				close(client.send)
				delete(clients, client)
			}
		}
		mu.Unlock()
	}
}
