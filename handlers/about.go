package handlers

import (
	"net/http"
	"text/template"
)

// AboutHandler handles the about page
func AboutHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(
		template.
			New("about.html").
			Delims("[[", "]]").
			ParseFiles("templates/about.html", "templates/footer.html"),
	)
	_ = t.Execute(w, map[string]interface{}{ // error cannot be triggered here in production
		"google_analytics_key": googleAnalyticsKey,
	})
}
