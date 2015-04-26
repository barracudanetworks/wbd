package web

import (
	"html/template"
	"log"
	"net/http"

	"github.com/johnmaguire/wbc/database"
)

type IndexHandler struct {
	address  string
	database string
}

func (ih *IndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := getClient("index", r)

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

	// Load template, parse vars, write to client
	t, _ := template.New("index").Parse(indexTemplate)
	t.Execute(w, struct {
		Client string
		URLs   []string
	}{
		id,
		urls,
	})
}

type WelcomeHandler struct {
	address  string
	database string
}

func (ih *WelcomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := getClient("welcome", r)

	// Load template, parse vars, write to client
	t, _ := template.New("welcome").Parse(welcomeTemplate)
	t.Execute(w, struct {
		Client     string
		RemoteAddr string
	}{
		id,
		r.RemoteAddr,
	})
}
