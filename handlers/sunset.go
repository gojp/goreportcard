package handlers

import (
	"log"
	"net/http"
)

// SunsetHandler handles the about page
func (gh *GRCHandler) SunsetHandler(w http.ResponseWriter, r *http.Request) {
	t, err := gh.loadTemplate("templates/sunset.html")
	if err != nil {
		log.Println("ERROR: could not get sunset template: ", err)
		http.Error(w, err.Error(), 500)
		return
	}

	if err := t.ExecuteTemplate(w, "base", map[string]interface{}{
		"google_analytics_key": googleAnalyticsKey,
	}); err != nil {
		log.Println("ERROR:", err)
	}
}
