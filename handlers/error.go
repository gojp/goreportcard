package handlers

import (
	"log"
	"net/http"
)

func (gh *GRCHandler) errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		t, err := gh.loadTemplate("/templates/404.html")
		if err != nil {
			log.Println("ERROR: could not get 404 template: ", err)
			http.Error(w, err.Error(), 500)
			return
		}

		if err := t.ExecuteTemplate(w, "base", nil); err != nil {
			log.Println("ERROR:", err)
		}
	}
}
