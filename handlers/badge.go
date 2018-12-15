package handlers

import (
	"fmt"
	"github.com/gojp/goreportcard/check"
	"log"
	"net/http"
	"strings"
)

func badgePath(grade check.Grade, style string, dev bool) string {
	return fmt.Sprintf("assets/badges/%s_%s.svg", strings.ToLower(string(grade)), strings.ToLower(style))
}

// BadgeHandler handles fetching the badge images
func BadgeHandler(w http.ResponseWriter, r *http.Request, repo string, dev bool) {
	resp, err := newChecksResp(repo, false)

	// See: http://shields.io/#styles
	style := r.URL.Query().Get("style")
	if style == "" {
		style = "flat"
	}

	if err != nil {
		log.Printf("ERROR: fetching badge for %s: %v", repo, err)
		url := "https://img.shields.io/badge/go%20report-error-lightgrey.svg?style=" + style
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		return
	}

	w.Header().Set("Cache-control", "no-store, no-cache, must-revalidate")
	http.ServeFile(w, r, badgePath(resp.Grade, style, dev))
}
