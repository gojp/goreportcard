package handlers

import (
	"log"
	"net/http"
)

// AboutHandler handles the about page
func AboutHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving about page")
	http.ServeFile(w, r, "templates/about.html")
}
