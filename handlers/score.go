package handlers

import (
	"container/heap"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/boltdb/bolt"
	"github.com/gojp/goreportcard/download"
)

// ScoresHandler handles the stats page
func ScoresHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	repo, err := download.Clean(r.FormValue("repo"))
	if err != nil {
		log.Println("ERROR: from download.Clean:", err)
		http.Error(w, "Could not download the repository: "+err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Checking repo %q...", repo)
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

	sortedScores := make([]scoreItem, len(*scores))
	fmt.Println("", sortedScores)
	score := "0"
	for i := range sortedScores {
		sortedScores[len(sortedScores)-i-1] = heap.Pop(scores).(scoreItem)
		fmt.Println("", sortedScores[i].Repo)
		if repo == sortedScores[i].Repo {
			score = strconv.FormatFloat(sortedScores[i].Score, 'f', -1, 64)
			break
		}
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(score))
}
