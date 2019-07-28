package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gojp/goreportcard/database"

	"github.com/dustin/go-humanize"
)

type scoreItem struct {
	Repo  string  `json:"repo"`
	Score float64 `json:"score"`
	Files int     `json:"files"`
}

func add(x, y int) int {
	return x + y
}

func formatScore(x float64) string {
	return fmt.Sprintf("%.2f", x)
}

// HighScoresHandler handles the stats page
type HighScoresHandler struct {
	DB database.Database
}

// Handle handles the stats page
func (h *HighScoresHandler) Handle(w http.ResponseWriter, r *http.Request) {

	funcs := template.FuncMap{"add": add, "formatScore": formatScore}
	t := template.Must(template.New("high_scores.html").Delims("[[", "]]").Funcs(funcs).ParseFiles("templates/high_scores.html", "templates/footer.html"))

	repos, err := h.DB.GetHighScores(100)
	if err != nil {
		log.Print("error loading high scores:", err)
		return
	}

	scores := make([]scoreItem, len(repos))
	for i, repo := range repos {
		cached, err := getFromCache(h.DB, repo)
		if err != nil {
			log.Print("error loading cached repo:", err)
			return
		}
		scores[i] = scoreItem{
			Repo:  repo,
			Score: cached.Average * 100,
			Files: cached.Files,
		}
	}

	count, err := h.DB.GetRepoCount()
	if err != nil {
		log.Print("error loading repo count:", err)
		return
	}

	t.Execute(w, map[string]interface{}{
		"HighScores":           scores,
		"Count":                humanize.Comma(int64(count)),
		"google_analytics_key": googleAnalyticsKey,
	})
}
