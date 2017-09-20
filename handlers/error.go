package handlers

import (
	"html/template"
	"net/http"
)

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		t := template.Must(template.New("404.html").ParseFiles("templates/404.html", "templates/footer.html"))
		t.Execute(w, nil)
	}
}
