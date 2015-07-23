package handlers

import (
	"net/http"
	"text/template"
)

// ShutdownHandler handles the shutdown page
func ShutdownHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.New("shutdown.html").ParseFiles("templates/shutdown.html"))

	t.Execute(w, nil)
}
