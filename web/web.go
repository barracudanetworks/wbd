package web

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type IndexHandler struct{}

func (ih *IndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("It works!"))
}

func Start(address string) {
	r := mux.NewRouter()

	// Routes
	ih := &IndexHandler{}
	r.Handle("/", ih)

	// Register mux router to http /
	http.Handle("/", r)

	log.Print("Starting web server")
	log.Fatal(http.ListenAndServe(address, nil))
}
