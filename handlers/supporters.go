package handlers

import (
	"net/http"
	"text/template"
)

// SupportersHandler handles the supporters page
func SupportersHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.New("supporters.html").Delims("[[", "]]").ParseFiles("templates/supporters.html", "templates/footer.html"))
	t.Execute(w, map[string]interface{}{
		"google_analytics_key": googleAnalyticsKey,
	})
}
