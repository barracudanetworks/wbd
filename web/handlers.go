package web

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type indexHandler struct{ App }

func (ih *indexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := ih.App.GetClient(r)

	// Web address to use in template
	addr := fmt.Sprintf("%s%s", r.Host, ih.App.Address)

	// Show welcome page by default
	defaultUrl := fmt.Sprintf("http://%s/welcome?client=%s", addr, c.Id)

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
		c.Id,
		urls,
	})
}

type welcomeHandler struct{ App }

func (wh *welcomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := wh.App.GetClient(r)

	// Load template, parse vars, write to client
	t, _ := template.New("welcome").Parse(welcomeTemplate)
	t.Execute(w, struct {
		Client     string
		RemoteAddr string
	}{
		c.Id,
		c.RemoteAddr,
	})
}

type consoleHandler struct{ App }

func (ah *consoleHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := ah.App.GetClient(r)

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
		c.Id,
		c.RemoteAddr,
	})
}

type websocketHandler struct{ App }

func (wh *websocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
	}

	c := wh.App.GetClient(r)

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	id := c.Id

	client := NewWebsocketClient(wh.Database, ws, id, r.RemoteAddr)

	hub.register <- client
	go client.writePump()
	client.readPump(wh.App.Database)
}
