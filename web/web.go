package web

import (
	"log"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
)

type Data struct {
	Title string
	URLs  []string
}

func globalData() Data {
	urls := make([]string, 4)
	urls[0] = "http://wallboard.bbs.cudaops.com/control/"
	urls[1] = "http://wallboard.bbs.cudaops.com/leapserv_count/"
	urls[2] = "https://www.dropcam.com/e/60493aca2b854ce892ad0b9a1c2511a2?autoplay=true"
	urls[3] = "http://wallboard.bbs.cudaops.com/versions/"

	return Data{
		Title: " Wallboard Control | Connecting... ",
		URLs:  urls,
	}
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	d := globalData()
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
