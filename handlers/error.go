package handlers

import (
	"net/http"
)

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		t := parse(tmplError)
		t.Execute(w, nil)
	}
}
