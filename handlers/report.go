package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"flag"
)

var domain = flag.String("domain", "goreportcard.com", "Domain used for your goreportcard installation")
var googleAnalyticsKey = flag.String("google_analytics_key", "UA-58936835-1", "Google Analytics Account Id")

// ReportHandler handles the report page
func ReportHandler(w http.ResponseWriter, r *http.Request, repo string, dev bool) {
	log.Printf("Displaying report: %q", repo)

	resp, err := getFromCache(repo)
	needToLoad := false
	if err != nil {
		log.Println("ERROR ReportHandler:", err) // log error, but continue
		needToLoad = true
	}

	respBytes, err := json.Marshal(resp)
	if err != nil {
		log.Println("ERROR ReportHandler: could not marshal JSON: ", err)
		http.Error(w, "Failed to load cache object", 500)
		return
	}

	t := parse(tmplReport)
	t.Execute(w, map[string]interface{}{
		"repo":                 repo,
		"response":             string(respBytes),
		"loading":              needToLoad,
		"domain":               domain,
		"google_analytics_key": googleAnalyticsKey,
	})
}
