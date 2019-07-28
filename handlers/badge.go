package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gojp/goreportcard/check"
	"github.com/gojp/goreportcard/database"
)

// BadgeHandler handles fetching badge images
type BadgeHandler struct {
	DB database.Database
}

func badgePath(grade check.Grade, style string) string {
	return fmt.Sprintf("assets/badges/%s_%s.svg", strings.ToLower(string(grade)), strings.ToLower(style))
}

// Handle handles fetching the badge images
func (b *BadgeHandler) Handle(w http.ResponseWriter, r *http.Request, repo string) {
	resp, err := newChecksResp(b.DB, repo, false)

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
	http.ServeFile(w, r, badgePath(resp.Grade, style))
}
