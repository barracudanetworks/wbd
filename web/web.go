package web

import (
	"fmt"
	"log"
	"net/http"

	"github.com/barracudanetworks/wbd/config"
	"github.com/barracudanetworks/wbd/database"

	"github.com/gorilla/mux"
)

type App struct {
	Address  string
	Database *database.Database
}

func (a *App) GetClient(r *http.Request) (id string) {
	id = r.FormValue("client")
	return
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

	r.Handle("/", &indexHandler{a})
	r.Handle("/ws", &websocketHandler{a})
	r.Handle("/welcome", &welcomeHandler{a})
	r.Handle("/console", &consoleHandler{a})

	// Register mux router
	http.Handle("/", r)

	addr := fmt.Sprintf("%s:%d", c.ListenAddress, c.ListenPort)
	log.Printf("Web server listening on http://%s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
