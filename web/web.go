package web

import (
	"log"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	// Create the template variable struct
	d := struct {
		Title string
		URLs  []string
	}{
		"Wallboard Control",
		make([]string, 4),
	}
	d.URLs[0] = "http://wallboard.bbs.cudaops.com/control/"
	d.URLs[1] = "http://wallboard.bbs.cudaops.com/leapserv_count/"
	d.URLs[2] = "https://www.dropcam.com/e/60493aca2b854ce892ad0b9a1c2511a2?autoplay=true"
	d.URLs[3] = "http://wallboard.bbs.cudaops.com/versions/"

	t, _ := template.ParseFiles("templates/index.html")
	t.Execute(w, d)
}

func Start(address string) {
	r := mux.NewRouter()

	r.HandleFunc("/", IndexHandler)

	// Register mux router to http /
	http.Handle("/", r)

	log.Print("Launching web server on http://", address)
	log.Fatal(http.ListenAndServe(address, nil))
}
