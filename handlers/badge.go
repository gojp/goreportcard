package handlers

import (
	"log"
	"net/http"
	"net/url"

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

	if err != nil || resp.DidError {
		log.Printf("ERROR: fetching badge for %s: %v", repo, err)
		url := "https://img.shields.io/badge/go%20report-error-lightgrey.svg?style=" + style
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		return
	}

	http.Redirect(w, r, badgeURL(resp.Grade, r.URL.Query()), http.StatusTemporaryRedirect)
}

func badgeURL(grade check.Grade, queryParams url.Values) string {
	badgeURL, _ := url.ParseRequestURI("https://img.shields.io/static/v1")

	style := queryParams.Get("style")
	color := queryParams.Get("color")
	labelColor := queryParams.Get("labelColor")
	logo := queryParams.Get("logo")
	logoColor := queryParams.Get("logoColor")
	logoWidth := queryParams.Get("logoWidth")
	label := queryParams.Get("label")

	if style == "" {
		style = "flat"
	}

	if color == "" {
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
	}

	if label == "" {
		label = "go report"
	}

	values := badgeURL.Query()
	values.Set("color", color)
	values.Set("style", style)
	values.Set("label", label)
	values.Set("message", string(grade))

	// optional parameters
	if labelColor != "" {
		values.Set("labelColor", labelColor)
	}

	if logo != "" {
		values.Set("logo", logo)
	}

	if logoColor != "" {
		values.Set("logoColor", logoColor)
	}

	if logoWidth != "" {
		values.Set("logoWidth", logoWidth)
	}

	badgeURL.RawQuery = values.Encode()

	return badgeURL.String()
}
