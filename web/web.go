package web

import (
	"log"
	"net/http"
	"strings"
	"text/template"

	"github.com/johnmaguire/wbc/database"

	"github.com/gorilla/mux"
)

type IndexHandler struct {
	address  string
	database string
}

func (ih *IndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.FormValue("client"))
	if id != "" {
		log.Printf("Client %s loaded index from %s", id, r.RemoteAddr)
	} else {
		log.Printf("User loaded index from %s", r.RemoteAddr)
	}

	// Connect to database
	db, err := database.Connect(ih.database)
	if err != nil {
		log.Fatal(err)
	}

	// Get URLs from database
	urls, err := db.FetchUrls()
	if err != nil {
		log.Fatal(err)
	}
	// Get URLs to show for this client
	// urls := getUrls(id)

	// Load template, parse vars, write to client
	t, _ := template.New("index").Parse(indexTemplate)
	t.Execute(w, struct {
		Title string
		URLs  []string
	}{
		"Wallboard Control",
		urls,
	})
}

func Start(address string, database string) {
	r := mux.NewRouter()

	ih := &IndexHandler{address, database}
	r.Handle("/", ih)

	// Register mux router to http /
	http.Handle("/", r)

	log.Print("Launching web server on http://", address)
	log.Fatal(http.ListenAndServe(address, nil))
}
