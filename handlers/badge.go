package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/dgraph-io/badger/v2"
	"github.com/gojp/goreportcard/check"
)

func badgePath(grade check.Grade, style string) string {
	return fmt.Sprintf("assets/badges/%s_%s.svg", strings.ToLower(string(grade)), strings.ToLower(style))
}

var badgeCache = sync.Map{}

// BadgeHandler handles fetching the badge images
func BadgeHandler(w http.ResponseWriter, r *http.Request, db *badger.DB, repo string) {
	// See: http://shields.io/#styles
	style := r.URL.Query().Get("style")
	if style == "" {
		style = "flat"
	}

	var grade check.Grade
	g, ok := badgeCache.Load(repo)
	if ok {
		log.Printf("Fetching badge for %q from cache...", repo)
		grade = g.(check.Grade)
	} else {
		resp, err := newChecksResp(db, repo, false)
		if err != nil {
			log.Printf("ERROR: fetching badge for %s: %v", repo, err)
			url := "https://img.shields.io/badge/go%20report-error-lightgrey.svg?style=" + style
			http.Redirect(w, r, url, http.StatusTemporaryRedirect)
			return
		}

		log.Printf("Adding badge for %q to cache...", repo)
		badgeCache.Store(repo, resp.Grade)
		grade = resp.Grade
	}

	w.Header().Set("Cache-control", "no-store, no-badgeCache, must-revalidate")
	http.ServeFile(w, r, badgePath(grade, style))
}
