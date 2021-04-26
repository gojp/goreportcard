package handlers

import (
	"testing"

	"github.com/gojp/goreportcard/check"
)

func TestBadgeURL(t *testing.T) {
	for grade, expectedURL := range map[check.Grade]string{
		check.GradeAPlus: "https://img.shields.io/badge/go%20report-A+-brightgreen.svg?style=for-the-badge",
		check.GradeA:     "https://img.shields.io/badge/go%20report-A-green.svg?style=for-the-badge",
		check.GradeB:     "https://img.shields.io/badge/go%20report-B-yellowgreen.svg?style=for-the-badge",
		check.GradeC:     "https://img.shields.io/badge/go%20report-C-yellow.svg?style=for-the-badge",
		check.GradeD:     "https://img.shields.io/badge/go%20report-D-orange.svg?style=for-the-badge",
		check.GradeE:     "https://img.shields.io/badge/go%20report-E-red.svg?style=for-the-badge",
		check.GradeF:     "https://img.shields.io/badge/go%20report-F-red.svg?style=for-the-badge",
	} {
		grade := grade
		expectedURL := expectedURL
		t.Run(string(grade), func(t *testing.T) {
			t.Parallel()
			got := badgeURL(grade, "for-the-badge")
			if got != expectedURL {
				t.Errorf("expected %s, got %s", expectedURL, got)
			}
		})
	}
}
