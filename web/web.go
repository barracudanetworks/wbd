package web

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/barracudanetworks/wbd/config"
	"github.com/barracudanetworks/wbd/database"

	"github.com/gorilla/mux"
)

type App struct {
	Address  string
	Database *database.Database
}

type Client struct {
	Id         string
	RemoteAddr string
}

func (a *App) GetClient(r *http.Request) (client Client) {
	client.Id = r.FormValue("client")

	// attempt to set via X-Forwarded-For header
	client.RemoteAddr = r.Header.Get("X-Forwarded-For")

	// if none was set, use the ip address from the request, minus the port
	if client.RemoteAddr == "" {
		client.RemoteAddr = r.RemoteAddr[:strings.LastIndex(r.RemoteAddr, ":")]
	}

	return
}

func (a *App) Route(route string) http.Handler {
	var handler http.Handler

	switch {
	case route == "index":
		handler = &indexHandler{*a}
	case route == "websocket":
		handler = &websocketHandler{*a}
	case route == "welcome":
		handler = &welcomeHandler{*a}
	case route == "console":
		handler = &consoleHandler{*a}
	}

	wrapper := func(w http.ResponseWriter, r *http.Request) {
		// fetch identifying info about client
		c := a.GetClient(r)

		// log header
		if c.Id == "" {
			log.Printf("Anonymous client loaded %s from %s", route, c.RemoteAddr)
		} else {
			log.Printf("Client %s loaded %s from %s", route, c.Id, c.RemoteAddr)
		}

		handler.ServeHTTP(w, r)
	}

	return http.HandlerFunc(wrapper)
}

func Start(c *config.Configuration) {
	r := mux.NewRouter()

	db, err := database.Connect(c.Database)
	if err != nil {
		log.Fatal(err)
	}

	a := App{
		Address:  c.WebAddress,
		Database: db,
	}

	// Goroutine the websocket loop
	go hub.run(&a)

	r.Handle("/", a.Route("index"))
	r.Handle("/ws", a.Route("websocket"))
	r.Handle("/welcome", a.Route("welcome"))
	r.Handle("/console", a.Route("console"))

	// Register mux router
	http.Handle("/", r)

	addr := fmt.Sprintf("%s:%d", c.ListenAddress, c.ListenPort)
	log.Printf("Web server listening on http://%s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
