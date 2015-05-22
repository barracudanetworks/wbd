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
	id := ih.App.GetClient(r)
	if id != "" {
		log.Printf("Client '%s' loaded index from %s", id, r.RemoteAddr)
	} else {
		log.Printf("User loaded index from %s", r.RemoteAddr)
	}

	// Web address to use in template
	addr := fmt.Sprintf("%s%s", r.Host, ih.App.Address)

	// Show welcome page by default
	defaultUrl := fmt.Sprintf("http://%s/welcome?client=%s", addr, id)

	// Get URLs from database
	urls, err := ih.App.Database.FetchUrls()
	if err != nil {
		log.Println(err)
		return
	}

	// Load template, parse vars, write to client
	t, err := template.New("index").Parse(indexTemplate)
	if err != nil {
		log.Println(err)
		return
	}
	t.Execute(w, struct {
		Address    template.URL
		DefaultUrl template.URL
		Client     string
		URLs       []string
	}{
		template.URL(addr),
		template.URL(defaultUrl),
		id,
		urls,
	})
}

type welcomeHandler struct{ App }

func (wh *welcomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := wh.App.GetClient(r)

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

type consoleHandler struct{ App }

func (ah *consoleHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := ah.App.GetClient(r)

	// Web address to use in template
	addr := fmt.Sprintf("%s%s", r.Host, ah.App.Address)

	// Load template, parse vars, write to client
	t, _ := template.New("console").Parse(consoleTemplate)
	t.Execute(w, struct {
		Address    template.URL
		Client     string
		RemoteAddr string
	}{
		template.URL(addr),
		id,
		r.RemoteAddr[:strings.Index(r.RemoteAddr, ":")],
	})
}

type websocketHandler struct{ App }

func (wh *websocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
	}

	// track whether this is a "registered" client or not
	var generic bool

	id := wh.App.GetClient(r)
	if id == "" {
		generic = true

		id = fmt.Sprintf("User (%d)", time.Now().UnixNano())
		log.Printf("%s connected to websocket from %s", id, r.RemoteAddr)
	} else {
		generic = false

		log.Printf("Client '%s' connected to websocket from %s", id, r.RemoteAddr)
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	c := &websocketClient{
		Id:        id,
		Generic:   generic,
		IpAddress: r.RemoteAddr,
		send:      make(chan *websocketMessage),
		ws:        ws,
	}

	hub.register <- c
	go c.writePump()
	c.readPump(wh.App.Database)
}
