package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dgraph-io/badger/v2"
	"github.com/gojp/goreportcard/check"
)

// BadgeHandler handles fetching the badge images
func BadgeHandler(w http.ResponseWriter, r *http.Request, db *badger.DB, repo string) {
	resp, err := newChecksResp(db, repo, false)

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

	http.Redirect(w, r, badgeURL(resp.Grade, style), http.StatusTemporaryRedirect)
}

func badgeURL(grade check.Grade, style string) string {
	var color string
	switch grade {
	case check.GradeAPlus:
		color = "brightgreen"
	case check.GradeA:
		color = "green"
	case check.GradeB:
		color = "yellowgreen"
	case check.GradeC:
		color = "yellow"
	case check.GradeD:
		color = "orange"
	case check.GradeE:
		fallthrough
	case check.GradeF:
		color = "red"
	}
	return fmt.Sprintf("https://img.shields.io/badge/go%%20report-%s-%s.svg?style=%s", grade, color, style)
}
