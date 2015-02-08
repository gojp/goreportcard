package handlers

import "net/http"

func ReportHandler(w http.ResponseWriter, r *http.Request, org, repo string) {
	http.ServeFile(w, r, "templates/home.html")
}
