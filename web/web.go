package web

import (
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func Start(address string, database string) {
	r := mux.NewRouter()

	ih := &IndexHandler{address, database}
	wh := &WelcomeHandler{address, database}
	r.Handle("/", ih)
	r.Handle("/welcome", wh)

	// Register mux router to http /
	http.Handle("/", r)

	log.Print("Launching web server on http://", address)
	log.Fatal(http.ListenAndServe(address, nil))
}

func getClient(page string, r *http.Request) (id string) {
	id = strings.TrimSpace(r.FormValue("client"))
	if id != "" {
		log.Printf("Client %s loaded %s from %s", id, page, r.RemoteAddr)
	} else {
		log.Printf("User loaded %s from %s", page, r.RemoteAddr)
	}
	return
}
