package web

import (
	"fmt"
	"log"
	"net/http"

	"github.com/barracudanetworks/wbc/config"
	"github.com/barracudanetworks/wbc/database"

	"github.com/gorilla/mux"
)

type App struct {
	Address  string
	Database *database.Database
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

	// Start websocket hub
	go h.run(&a)

	r.Handle("/", &indexHandler{a})
	r.Handle("/ws", &websocketHandler{a})
	r.Handle("/welcome", &welcomeHandler{a})

	// Register mux router
	http.Handle("/", r)

	addr := fmt.Sprintf("%s:%d", c.ListenAddress, c.ListenPort)
	log.Printf("Web server listening on http://%s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func getClient(page string, r *http.Request) (id string) {
	id = r.FormValue("client")
	if id != "" {
		log.Printf("Client %s loaded %s from %s", id, page, r.RemoteAddr)
	} else {
		log.Printf("User loaded %s from %s", page, r.RemoteAddr)
	}
	return
}
