package handlers

import (
	"log"
	"net/http"
	"strings"
)

// AssetsHandler handles serving static files
func AssetsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving " + r.URL.Path[1:])

	if strings.HasSuffix(r.URL.Path, ".svg") {
		// don't cache badges
		w.Header().Set("Cache-control", "no-store, no-cache, must-revalidate")
	}

	http.ServeFile(w, r, r.URL.Path[1:])
}

// FaviconHandler handles serving the favicon.ico
func FaviconHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving " + r.URL.Path[1:])
	http.ServeFile(w, r, "assets/favicon.ico")
}
