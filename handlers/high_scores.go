package handlers

import (
	"container/heap"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/boltdb/bolt"
	"github.com/dustin/go-humanize"
)

func add(x, y int) int {
	return x + y
}

func formatScore(x float64) string {
	return fmt.Sprintf("%.2f", x)
}

// HighScoresHandler handles the stats page
func HighScoresHandler(w http.ResponseWriter, r *http.Request) {
	// write to boltdb
	db, err := bolt.Open(DBPath, 0755, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Println("Failed to open bolt database: ", err)
		return
	}
	defer db.Close()

	count, scores := 0, &ScoreHeap{}
	err = db.View(func(tx *bolt.Tx) error {
		hsb := tx.Bucket([]byte(MetaBucket))
		if hsb == nil {
			return fmt.Errorf("high score bucket not found")
		}
		scoreBytes := hsb.Get([]byte("scores"))
		if scoreBytes == nil {
			scoreBytes, err = json.Marshal([]ScoreHeap{})
			if err != nil {
				return err
			}
		}
		json.Unmarshal(scoreBytes, scores)

		heap.Init(scores)

		total := hsb.Get([]byte("total_repos"))
		if total == nil {
			count = 0
			return nil
		}
		return json.Unmarshal(total, &count)
	})

	if err != nil {
		log.Println("ERROR: Failed to load high scores from bolt database: ", err)
		http.Error(w, err.Error(), 500)
		return
	}

	funcs := template.FuncMap{"add": add, "formatScore": formatScore}
	t := template.Must(template.New("high_scores.html").Delims("[[", "]]").Funcs(funcs).ParseFiles("templates/high_scores.html", "templates/footer.html"))

	sortedScores := make([]scoreItem, len(*scores))
	for i := range sortedScores {
		sortedScores[len(sortedScores)-i-1] = heap.Pop(scores).(scoreItem)
	}

	t.Execute(w, map[string]interface{}{
		"HighScores":           sortedScores,
		"Count":                humanize.Comma(int64(count)),
		"google_analytics_key": googleAnalyticsKey,
	})
}
