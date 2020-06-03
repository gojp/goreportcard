package httpapi

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/gojp/goreportcard/internal/model"

	"github.com/yeqown/log"
)

// AssetsHandler handles serving static files
func AssetsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "max-age=86400")

	http.ServeFile(w, r, r.URL.Path[1:])
}

// FaviconHandler handles serving the favicon.ico
func FaviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "assets/favicon.ico")
}

// BadgeHandler handles fetching the badge images
func BadgeHandler(w http.ResponseWriter, r *http.Request, repo string) {
	// See: http://shields.io/#styles
	style := r.URL.Query().Get("style")
	if style == "" {
		style = "flat"
	}

	var grade model.Grade
	g, ok := badgeCache.Load(repo)
	if ok {
		log.Infof("Fetching badge for %q from cache...", repo)
		grade = g.(model.Grade)
	} else {
		resp, err := dolint(repo, false)
		if err != nil {
			log.Errorf("fetching badge for %s: %v", repo, err)
			url := "https://img.shields.io/badge/go%20report-error-lightgrey.svg?style=" + style
			http.Redirect(w, r, url, http.StatusTemporaryRedirect)
			return
		}

		log.Infof("Adding badge for %q to cache...", repo)
		badgeCache.Store(repo, resp.Grade)
		grade = resp.Grade
	}

	w.Header().Set("Cache-control", "no-store, no-badgeCache, must-revalidate")
	http.ServeFile(w, r, badgePath(grade, style))
}

var badgeCache = sync.Map{}

func badgePath(grade model.Grade, style string) string {
	return fmt.Sprintf("assets/badges/%s_%s.svg", strings.ToLower(string(grade)), strings.ToLower(style))
}
