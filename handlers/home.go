package handlers

import (
	"log"
	"net/http"
)

// HomeHandler handles the homepage
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving home page")
	if r.URL.Path[1:] == "" {
		http.ServeFile(w, r, "templates/home.html")
	} else {
		http.NotFound(w, r)
	}
}
