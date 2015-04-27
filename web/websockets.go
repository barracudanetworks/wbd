package web

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

const (
	writeWait = 10 * time.Second

	pongWait = 60 * time.Second
	pingWait = (pongWait * 9) / 10

	maxMessageSize = 512
)

// The map is a bit weird in that it's pointer->string, but it's kind of a cheap
// hack around issues with overwriting a connection that has the same client name
type websocketHub struct {
	broadcast   chan []byte
	register    chan *websocketClient
	unregister  chan *websocketClient
	connections map[*websocketClient]string
}

var h = websocketHub{
	broadcast:   make(chan []byte),
	register:    make(chan *websocketClient),
	unregister:  make(chan *websocketClient),
	connections: make(map[*websocketClient]string),
}

func (h *websocketHub) run() {
	for {
		select {
		// Save connection to hub
		case c := <-h.register:
			log.Printf("Added client '%s' to the hub", c.Id)
			h.connections[c] = c.Id

		// Remove connection from hub
		case c := <-h.unregister:
			if _, ok := h.connections[c]; ok {
				log.Printf("Removing client '%s'", c.Id)

				delete(h.connections, c)
				close(c.send)
			}

		// Broadcast messages to clients
		case m := <-h.broadcast:
			for c := range h.connections {
				select {
				case c.send <- m:
					log.Printf("Sent to client '%s': %s", c.Id, string(m))
				default:
					log.Printf("Client '%s' looks to have gone away", c.Id)

					close(c.send)
					delete(h.connections, c)
				}
			}
		}
	}
}

type websocketClient struct {
	Id   string
	ws   *websocket.Conn
	send chan []byte
}

func (c *websocketClient) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

func (c *websocketClient) writePump() {
	ticker := time.NewTicker(pingWait)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()

	for {
		select {
		// If there's a message in the send queue, write to the socket
		case message, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.write(websocket.TextMessage, message); err != nil {
				return
			}
		// Send ping on timer
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func (c *websocketClient) readPump() {
	defer func() {
		h.unregister <- c
		c.ws.Close()
	}()

	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			log.Printf("Connection closed to client '%s'", c.Id)
			break
		}

		log.Printf("Received from client '%s': %s", c.Id, string(message))
		h.broadcast <- message
	}
}
