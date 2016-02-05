package main

import (
	"time"

	"github.com/gojp/goreportcard/check"
)

const (
	MongoConnection string = "localhost:10027"

	// DBPath is the relative (or absolute) path to the bolt database file
	DBPath string = "goreportcard.db"

	// RepoBucket is the bucket in which repos will be cached in the bolt DB
	RepoBucket string = "repos"

	// HighScoreBucket is the bucket containing the names of the projects with the
	// top 100 high scores
	HighScoreBucket string = "high_scores"
)

type Grade string

type score struct {
	Name          string              `json:"name"`
	Description   string              `json:"description"`
	FileSummaries []check.FileSummary `json:"file_summaries"`
	Weight        float64             `json:"weight"`
	Percentage    float64             `json:"percentage"`
}

type checksResp struct {
	Checks      []score   `json:"checks"`
	Average     float64   `json:"average"`
	Grade       Grade     `json:"grade"`
	Files       int       `json:"files"`
	Issues      int       `json:"issues"`
	Repo        string    `json:"repo"`
	LastRefresh time.Time `json:"last_refresh"`
}

func main() {

}
