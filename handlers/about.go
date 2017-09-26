package handlers

import (
	"log"
	"net/http"
)

// AboutHandler handles the about page
func AboutHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving about page")

	t := parse(tmplAbout)
	t.Execute(w, map[string]interface{}{
		"google_analytics_key": googleAnalyticsKey,
	})
	return
}
