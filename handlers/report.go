package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"flag"

	"github.com/dgraph-io/badger/v2"
)

var domain = flag.String("domain", "goreportcard.com", "Domain used for your goreportcard installation")
var googleAnalyticsKey = flag.String("google_analytics_key", "G-TFTF5Y92QD", "Google Analytics Account ID (GA4)")

// ReportHandler handles the report page
func (gh *GRCHandler) ReportHandler(w http.ResponseWriter, r *http.Request, db *badger.DB, repo string) {
	log.Printf("Displaying report: %q", repo)
	t, err := gh.loadTemplate("/templates/report.html")
	if err != nil {
		log.Println("ERROR: could not get report template: ", err)
		http.Error(w, err.Error(), 500)
		return
	}

	resp, err := getFromCache(db, repo)
	needToLoad := false
	if err != nil {
		switch err.(type) {
		case notFoundError:
			// don't bother logging - we already log in getFromCache. continue
		default:
			log.Println("ERROR ReportHandler:", err) // log error, but continue
		}
		needToLoad = true
	}

	respBytes, err := json.Marshal(resp)
	if err != nil {
		log.Println("ERROR ReportHandler: could not marshal JSON: ", err)
		http.Error(w, "Failed to load cache object", 500)
		return
	}

	if err := t.Execute(w, map[string]interface{}{
		"repo":                 repo,
		"response":             string(respBytes),
		"loading":              needToLoad,
		"domain":               domain,
		"google_analytics_key": googleAnalyticsKey,
	}); err != nil {
		log.Println("ERROR:", err)
	}
}
