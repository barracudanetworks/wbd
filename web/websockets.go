package web

import (
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"github.com/barracudanetworks/wbd/database"
	"github.com/gorilla/websocket"
)

const (
	writeWait = 10 * time.Second

	pongWait = 60 * time.Second
	pingWait = (pongWait * 9) / 10

	urlPollWait = pongWait / 2

	maxMessageSize = 512
)

type websocketMessage struct {
	Action string      `json:"action"`
	Data   interface{} `json:"data"`
}

// The map is a bit weird in that it's pointer->string, but it's kind of a cheap
// hack around issues with overwriting a connection that has the same client name
type websocketHub struct {
	broadcast   chan *websocketMessage
	register    chan *websocketClient
	unregister  chan *websocketClient
	connections map[*websocketClient]string
}

var hub = websocketHub{
	broadcast:   make(chan *websocketMessage),
	register:    make(chan *websocketClient),
	unregister:  make(chan *websocketClient),
	connections: make(map[*websocketClient]string),
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (h *websocketHub) run(a *App) {
	db := a.Database

	ticker := time.NewTicker(urlPollWait)
	defer func() {
		ticker.Stop()
	}()

	for {
		select {
		// Save connection to hub
		case c := <-h.register:
			log.Printf("Added client '%s' to the hub", c.Id)

			if !c.Generic {
				// get client info and touch last ping
				client, err := db.GetClient(c.Id)
				switch {
				case err == sql.ErrNoRows:
					log.Printf("Unknown client, creating record")
					db.InsertClient(c.Id, c.IpAddress)
				case err != nil:
					log.Fatal(err)
				default:
					log.Printf("Client last seen %s from %s", client.LastPing, client.IpAddress)
					db.TouchClient(client.Identifier)
				}
			} else {
				log.Printf("Not attempting to track generic client")
			}

			h.connections[c] = c.Id

		// Remove connection from hub
		case c := <-h.unregister:
			if _, ok := h.connections[c]; ok {
				log.Printf("Removing client '%s'", c.Id)
				h.CloseConnection(c)
			}

		// Broadcast messages to clients
		case m := <-h.broadcast:
			for c := range h.connections {
				select {
				case c.send <- m:
					log.Printf("Sent message to client '%s', type '%s'", c.Id, m.Action)
				default:
					h.CloseConnection(c)
				}
			}

		// Send out URL updates every so often
		case <-ticker.C:
			log.Print("Polling for URL changes in database")

			urls, err := db.FetchUrls()
			if err != nil {
				log.Fatal(err)
			}

			urlWm, err := urlUpdateMessage(urls)
			if err != nil {
				log.Fatal(err)
			}

			// Send the JSON to all connected clients
			for c := range h.connections {
				select {
				case c.send <- urlWm:
					log.Printf("Sent updated URL list to client '%s'", c.Id)
				default:
					h.CloseConnection(c)
				}
			}
		}
	}
}

func (h *websocketHub) CloseConnection(c *websocketClient) {
	log.Printf("Closing connection to client '%s'", c.Id)
	close(c.send)
	delete(h.connections, c)
}

type websocketClient struct {
	Id         string
	IpAddress  string
	Controller bool
	Generic    bool

	ws   *websocket.Conn
	send chan *websocketMessage
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

			json, err := json.Marshal(message)
			if err != nil {
				log.Print(err)
				return
			}

			if err := c.write(websocket.TextMessage, json); err != nil {
				log.Print(err)
				return
			}

		// Send ping on timer
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				log.Print(err)
				return
			}
		}
	}
}

func (c *websocketClient) readPump(db *database.Database) {
	defer func() {
		hub.unregister <- c
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

		var wm websocketMessage
		if err := json.Unmarshal(message, &wm); err != nil {
			log.Print(err)
			break
		}

		// Respond to a requested action by the client
		switch wm.Action {
		case "flagController":
			log.Printf("Client '%s' flagged as a controller", c.Id)
			c.Controller = true

			urls, err := db.FetchUrls()
			if err != nil {
				log.Fatal(err)
			}

			urlWm, err := urlUpdateMessage(urls)
			if err != nil {
				log.Fatal(err)
			}

			clientWm := hub.clientUpdateMessage()

			c.send <- urlWm
			c.send <- clientWm
		case "sendUrls":
			log.Printf("Client '%s' requested URLs", c.Id)

			var fetch_global bool = true
			var urls []string

			if !c.Generic {
				info, err := db.GetClient(c.Id)

				if err != nil {
					log.Printf("Couldn't find client info -- fetching all URLs")
				} else if info.UrlListId != 0 {
					log.Printf("Fetching URLs from assigned list")

					urls, err = db.FetchListUrlsById(info.UrlListId)
					if err != nil {
						log.Printf("Failed to fetch client's list URLs -- fetching all URLs")
					} else {
						// found urls already, no need to fetch global url list
						log.Printf("Fetched URLs for client '%s' from list ID %d", c.Id, info.UrlListId)
						fetch_global = false
					}
				}
			}

			if fetch_global {
				urls, err = db.FetchUrls()
				if err != nil {
					log.Fatal("Unable to fetch global URLs")
				}
			}

			urlWm, err := urlUpdateMessage(urls)
			if err != nil {
				log.Fatal(err)
			}

			c.send <- urlWm
		case "sendClients":
			log.Printf("Client '%s' requested clients", c.Id)

			clientWm := hub.clientUpdateMessage()

			c.send <- clientWm
		default:
			log.Printf("Unknown action %s from client '%s'", wm.Action, c.Id)
		}
	}
}

func urlUpdateMessage(urls []string) (wm *websocketMessage, err error) {
	wm = &websocketMessage{
		Action: "updateUrls",
		Data: struct {
			URLs []string `json:"urls"`
		}{
			urls,
		},
	}

	return
}

func (h *websocketHub) GetClients() (clients []string) {
	for c := range h.connections {
		if c.Id != "" {
			clients = append(clients, c.Id)
		}
	}

	return
}

func (h *websocketHub) clientUpdateMessage() (wm *websocketMessage) {
	clients := h.GetClients()

	wm = &websocketMessage{
		Action: "updateClients",
		Data: struct {
			Clients []string `json:"clients"`
		}{
			clients,
		},
	}

	return
}
