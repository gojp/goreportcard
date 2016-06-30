package handlers

import (
	"fmt"
	"log"
	"net/http"
)

// Grade represents a grade returned by the server, which is normally
// somewhere between A+ (highest) and F (lowest).
type Grade string

// The Grade constants below indicate the current available
// grades.
const (
	GradeAPlus Grade = "A+"
	GradeA           = "A"
	GradeB           = "B"
	GradeC           = "C"
	GradeD           = "D"
	GradeE           = "E"
	GradeF           = "F"
)

// grade is a helper for getting the grade for a percentage
func grade(percentage float64) Grade {
	switch {
	case percentage > 90:
		return GradeAPlus
	case percentage > 80:
		return GradeA
	case percentage > 70:
		return GradeB
	case percentage > 60:
		return GradeC
	case percentage > 50:
		return GradeD
	case percentage > 40:
		return GradeE
	default:
		return GradeF
	}
}

func badgeURL(grade Grade, style string) string {
	colorMap := map[Grade]string{
		GradeAPlus: "brightgreen",
		GradeA:     "brightgreen",
		GradeB:     "yellowgreen",
		GradeC:     "yellow",
		GradeD:     "orange",
		GradeE:     "red",
		GradeF:     "red",
	}
	url := fmt.Sprintf("https://img.shields.io/badge/go_report-%s-%s.svg?style=%s", grade, colorMap[grade], style)
	return url
}

// BadgeHandler handles fetching the badge images
func BadgeHandler(w http.ResponseWriter, r *http.Request, repo string) {
	name := fmt.Sprintf("%s", repo)
	resp, err := newChecksResp(name, false)

	// See: http://shields.io/#styles
	style := r.URL.Query().Get("style")
	if style == "" {
		style = "flat"
	}

	if err != nil {
		log.Printf("ERROR: fetching badge for %s: %v", name, err)
		url := "https://img.shields.io/badge/go%20report-error-lightgrey.svg?style=" + style
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		return
	}

	http.Redirect(w, r, badgeURL(resp.Grade, style), http.StatusTemporaryRedirect)
}
