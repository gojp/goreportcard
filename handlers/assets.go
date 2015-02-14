package handlers

import (
	"log"
	"net/http"
)

// AssetsHandler handles serving static files
func AssetsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving " + r.URL.Path[1:])
	http.ServeFile(w, r, r.URL.Path[1:])
}
