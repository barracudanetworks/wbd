package web

import (
	"log"
	"net/http"
)

type indexHandler struct{}

func (ih *indexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("It works!"))
}

func Start(address string) {
	ih := &indexHandler{}
	http.Handle("/", ih)

	log.Print("Starting web server")

	log.Fatal(http.ListenAndServe(address, nil))
}
