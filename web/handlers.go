package web

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"
)

type indexHandler struct{ App }

func (ih *indexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := getClient("index", r)

	// Get URLs from database
	urls, err := ih.App.Database.FetchUrls()
	if err != nil {
		log.Println(err)
		return
	}

	// Load template, parse vars, write to client
	t, _ := template.New("index").Parse(indexTemplate)
	t.Execute(w, struct {
		Address string
		Client  string
		URLs    []string
	}{
		ih.App.Address,
		id,
		urls,
	})
}

type welcomeHandler struct{ App }

func (wh *welcomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := getClient("welcome", r)

	// Load template, parse vars, write to client
	t, _ := template.New("welcome").Parse(welcomeTemplate)
	t.Execute(w, struct {
		Client     string
		RemoteAddr string
	}{
		id,
		r.RemoteAddr[:strings.Index(r.RemoteAddr, ":")],
	})
}

type websocketHandler struct{ App }

func (wh *websocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
	}
	id := getClient("websocket endpoint", r)
	if id == "" {
		id = fmt.Sprintf("User (%d)", time.Now().UnixNano())
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	c := &websocketClient{
		Id:   id,
		send: make(chan []byte, 256),
		ws:   ws,
	}

	h.register <- c
	go c.writePump()
	c.readPump(wh.App.Database)
}
