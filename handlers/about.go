package handlers

import (
	"log"
	"net/http"
)

// AboutHandler handles the about page
func (gh *GRCHandler) AboutHandler(w http.ResponseWriter, r *http.Request) {
	t, err := gh.loadTemplate("templates/about.html")
	if err != nil {
		log.Println("ERROR: could not get about template: ", err)
		http.Error(w, err.Error(), 500)
		return
	}

	t.Execute(w, map[string]interface{}{
		"google_analytics_key": googleAnalyticsKey,
	})
}
