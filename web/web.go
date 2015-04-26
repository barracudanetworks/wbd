package web

import (
	"log"
	"net/http"
	"strings"
	"text/template"

	"github.com/gorilla/mux"
)

// TODO: Load from database
func getUrls(id string) (urls []string) {
	urls = append(urls,
		"http://wallboard.bbs.cudaops.com/control/",
		"http://wallboard.bbs.cudaops.com/leapserv_count/",
		"https://www.dropcam.com/e/60493aca2b854ce892ad0b9a1c2511a2?autoplay=true",
		"http://wallboard.bbs.cudaops.com/versions/",
	)
	return
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.FormValue("client"))
	if id != "" {
		log.Printf("Client %s loaded index from %s", id, r.RemoteAddr)
	} else {
		log.Printf("User loaded index from %s", r.RemoteAddr)
	}

	// Get URLs to show for this client
	urls := getUrls(id)

	// Load template, parse vars, write to client
	t, _ := template.ParseFiles("templates/index.html")
	t.Execute(w, struct {
		Title string
		URLs  []string
	}{
		"Wallboard Control",
		urls,
	})
}

func Start(address string) {
	r := mux.NewRouter()

	r.HandleFunc("/", IndexHandler)

	// Register mux router to http /
	http.Handle("/", r)

	log.Print("Launching web server on http://", address)
	log.Fatal(http.ListenAndServe(address, nil))
}
