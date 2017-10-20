package handlers

import (
	"net/http"
)

// AssetsHandler handles serving static files
func AssetsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "max-age=86400")

	http.ServeFile(w, r, r.URL.Path[1:])
}

// FaviconHandler handles serving the favicon.ico
func FaviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "assets/favicon.ico")
}
