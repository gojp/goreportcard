package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

type singleSuggestion struct {
	Name string `json:"name"`
}

// SuggestionHandler handles the autocomplete JSON request
func SuggestionHandler(w http.ResponseWriter, r *http.Request, ctx Context) {
	r.ParseForm()
	q := r.FormValue("q")
	log.Println("Serving suggestions for", q)

	suggestions, err := ctx.Suggest(q)
	if err != nil {
		log.Println("ERROR:", err)
		w.Write([]byte("[]"))
	}
	suggestionResponse := make([]singleSuggestion, len(suggestions))
	for i := range suggestions {
		suggestionResponse[i].Name = suggestions[i]
	}

	js, err := json.Marshal(suggestionResponse)
	if err != nil {
		log.Println("ERROR:", err)
	}
	w.Write(js)
}
