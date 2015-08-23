package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/gojp/goreportcard/handlers"
)

func makeHandler(name string, fn func(http.ResponseWriter, *http.Request, string, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		validPath := regexp.MustCompile(fmt.Sprintf(`^/%s/([a-zA-Z0-9\-_.]+)/([a-zA-Z0-9\-_.]+)$`, name))

		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[1], m[2])
	}
}

func main() {
	if err := os.MkdirAll("repos/src/github.com", 0755); err != nil && !os.IsExist(err) {
		log.Fatal("ERROR: could not create repos dir: ", err)
	}

	http.HandleFunc("/assets/", handlers.AssetsHandler)
	http.HandleFunc("/checks", handlers.CheckHandler)
	http.HandleFunc("/report/", makeHandler("report", handlers.ReportHandler))
	http.HandleFunc("/badge/", makeHandler("badge", handlers.BadgeHandler))
	http.HandleFunc("/high_scores/", handlers.HighScoresHandler)
	http.HandleFunc("/", handlers.HomeHandler)

	fmt.Println("Running on 127.0.01:8080...")
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
}
