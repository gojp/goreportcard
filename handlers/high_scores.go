package handlers

import (
	"container/heap"
	"encoding/json"
	"log"
	"net/http"

	"github.com/dgraph-io/badger/v2"
	"github.com/dustin/go-humanize"
)

// HighScoresHandler handles the stats page
func (gh *GRCHandler) HighScoresHandler(w http.ResponseWriter, r *http.Request, db *badger.DB) {
	count, scores := 0, &ScoreHeap{}
	err := db.View(func(txn *badger.Txn) error {
		var scoreBytes = []byte("[]")
		item, err := txn.Get([]byte("scores"))

		if item != nil {
			err = item.Value(func(val []byte) error {
				scoreBytes = val
				return nil
			})

			if err != nil {
				log.Println("ERROR:", err)
			}
		}

		json.Unmarshal(scoreBytes, scores)

		heap.Init(scores)

		item, err = txn.Get([]byte("total_repos"))
		if item == nil {
			count = 0
			return nil
		}

		if item != nil {
			err = item.Value(func(val []byte) error {
				return json.Unmarshal(val, &count)
			})
		}

		return err
	})

	if err != nil {
		log.Println("ERROR: Failed to load high scores from bolt database: ", err)
		http.Error(w, err.Error(), 500)
		return
	}

	t, err := gh.loadTemplate("/templates/high_scores.html")
	if err != nil {
		log.Println("ERROR: could not get high scores template: ", err)
		http.Error(w, err.Error(), 500)
		return
	}

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
