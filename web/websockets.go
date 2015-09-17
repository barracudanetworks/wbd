package web

import (
	"database/sql"
	"encoding/json"
	"log"
	"math/rand"
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
					if err := db.InsertClient(c.Id, c.IpAddress); err != nil {
						log.Fatal(err)
					}
				case err != nil:
					log.Fatal(err)
				default:
					log.Printf("Client last seen %s from %s", client.LastPing, client.IpAddress)

					c.Database = db
					c.Touch()
					c.UpdateIpAddress()
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

			// Send the JSON to all connected clients
			for c := range h.connections {
				urls, err := db.FetchUrlsByClientId(c.Id)
				if err != nil && err != sql.ErrNoRows {
					log.Fatal(err)
				}

				urlWm, err := urlUpdateMessage(urls)
				if err != nil {
					log.Fatal(err)
				}

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
	Database   *database.Database

	ws   *websocket.Conn
	send chan *websocketMessage
}

func NewWebsocketClient(db *database.Database, ws *websocket.Conn, id string, ipAddress string) (wc *websocketClient) {
	rand.Seed(time.Now().Unix())

	generic := false

	if id == "" {
		generic = true

		chars := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F"}
		id = "Anonymous "
		for i := 0; i < 3; i++ {
			id += chars[rand.Intn(len(chars))]
		}
	}

	wc = &websocketClient{
		Id:         id,
		IpAddress:  ipAddress,
		Controller: false,
		Generic:    generic,
		Database:   db,

		send: make(chan *websocketMessage),
		ws:   ws,
	}

	return
}

func (c *websocketClient) Touch() error {
	log.Printf("Updating last active timestamp for client '%s'", c.Id)
	return c.Database.TouchClient(c.Id)
}

func (c *websocketClient) UpdateIpAddress() error {
	log.Printf("Updating IP address in DB for client '%s'", c.Id)
	return c.Database.SetClientIpAddress(c.Id, c.IpAddress)
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
	c.ws.SetPongHandler(
		func(string) error {
			c.ws.SetReadDeadline(time.Now().Add(pongWait))
			c.Touch()
			return nil
		})

	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			log.Printf("Connection closed to client '%s'", c.Id)
			break
		}

		c.Touch()

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

			urls, err := db.FetchUrlsByClientId(c.Id)
			if err != nil && err != sql.ErrNoRows {
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

			urls, err := db.FetchUrlsByClientId(c.Id)
			if err != nil && err != sql.ErrNoRows {
				log.Fatal(err)
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
	// Prevent empty slices from being marshalled into null
	if urls == nil {
		urls = make([]string, 0)
	}

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
