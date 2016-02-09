package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/alecthomas/template"
)

// ReportHandler handles the report page
func ReportHandler(w http.ResponseWriter, r *http.Request, repo string) {
	log.Println("report", repo)
	t := template.Must(template.New("report.html").Delims("[[", "]]").ParseFiles("templates/report.html"))
	resp, err := getFromCache(repo)
	needToLoad := false
	if err != nil {
		log.Println("ERROR:", err) // log error, but continue
		needToLoad = true
	}

	respBytes, err := json.Marshal(resp)
	if err != nil {
		log.Println("ERROR: marshaling json: ", err)
		http.Error(w, "Failed to load cache object", 500)
		return
	}

	t.Execute(w, map[string]interface{}{"repo": repo, "response": string(respBytes), "loading": needToLoad})
}
