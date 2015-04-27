package web

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/johnmaguire/wbc/database"
)

type App struct {
	Address  string
	Database *database.Database
}

func Start(address string, dbp string) {
	r := mux.NewRouter()

	// Start websocket hub
	go h.run()

	db, err := database.Connect(dbp)
	if err != nil {
		log.Fatal(err)
	}

	a := App{address, db}

	r.Handle("/", &indexHandler{a})
	r.Handle("/ws", &websocketHandler{a})
	r.Handle("/welcome", &welcomeHandler{a})

	// Register mux router
	http.Handle("/", r)

	log.Print("Launching web server on http://", address)
	log.Fatal(http.ListenAndServe(address, nil))
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
