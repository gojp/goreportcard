package handlers

import "net/http"

// ReportHandler handles the report page
func ReportHandler(w http.ResponseWriter, r *http.Request, org, repo string) {
	http.ServeFile(w, r, "templates/home.html")
}
