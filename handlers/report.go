package handlers

import (
	"encoding/json"
	"flag"
	"html/template"
	"log"
	"net/http"

	"github.com/dgraph-io/badger"
)

var domain = flag.String("domain", "goreportcard.com", "Domain used for your goreportcard installation")
var googleAnalyticsKey = flag.String("google_analytics_key", "UA-58936835-1", "Google Analytics Account Id")

// ReportHandler handles the report page
func ReportHandler(w http.ResponseWriter, r *http.Request, db *badger.DB, repo string) {
	log.Printf("Displaying report: %q", repo)
	t := template.Must(template.New("report.html").Delims("[[", "]]").ParseFiles("templates/report.html", "templates/footer.html"))
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

	_ = t.Execute(w, map[string]interface{}{
		"repo":                 repo,
		"response":             string(respBytes),
		"loading":              needToLoad,
		"domain":               domain,
		"google_analytics_key": googleAnalyticsKey,
	})
}
