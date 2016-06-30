package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

// SuggestionHandler handles the autocomplete JSON request
func SuggestionHandler(w http.ResponseWriter, r *http.Request, ctx Context) {
	r.ParseForm()
	q := r.FormValue("q")
	log.Println("Serving suggestions for", q)

	suggestions := ctx.Suggest(q)

	js, err := json.Marshal(suggestions)
	if err != nil {
		log.Println("ERROR:", err)
	}
	w.Write(js)
}
