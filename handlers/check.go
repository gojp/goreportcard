package handlers

import (
	"container/heap"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/boltdb/bolt"
)

const (
	// DBPath is the relative (or absolute) path to the bolt database file
	DBPath string = "goreportcard.db"

	// RepoBucket is the bucket in which repos will be cached in the bolt DB
	RepoBucket string = "repos"

	// HighScoreBucket is the bucket containing the names of the projects with the
	// top 100 high scores
	HighScoreBucket string = "high_scores"
)

// CheckHandler handles the request for checking a repo
func CheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	repo := r.FormValue("repo")
	log.Printf("Checking repo %s...", repo)
	if strings.ToLower(repo) == "golang/go" {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("We've decided to omit results for the Go repository because it has lots of test files that (purposely) don't pass our checks. Go gets an A+ in our books though!"))
		return
	}
	forceRefresh := r.Method != "GET" // if this is a GET request, try fetch from cached version in mongo first
	resp, err := newChecksResp(repo, forceRefresh)
	if err != nil {
		log.Println("ERROR: ", err)
		b, _ := json.Marshal(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(b)
		return
	}

	respBytes, err := json.Marshal(resp)
	if err != nil {
		log.Println("ERROR: could not marshal json:", err)
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(respBytes)

	// write to boltdb
	db, err := bolt.Open(DBPath, 0755, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Println("Failed to open bolt database: ", err)
		return
	}
	defer db.Close()

	log.Printf("Saving repo %q to cache...", repo)

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(RepoBucket))
		if b == nil {
			return fmt.Errorf("Repo bucket not found")
		}

		// is this a new repo? if so, increase the count in the high scores bucket later
		isNewRepo := b.Get([]byte(repo)) == nil

		err := b.Put([]byte(repo), respBytes)
		if err != nil {
			return err
		}

		// check if we might need to update the high score list
		if resp.Files < 100 {
			// only repos with >= 100 files are considered for the high score list
			return nil
		}

		hsb := tx.Bucket([]byte(HighScoreBucket))
		if hsb == nil {
			return fmt.Errorf("High score bucket not found")
		}
		// update total repos count
		if isNewRepo {
			totalInt := 0
			total := hsb.Get([]byte("total_repos"))
			if total != nil {
				err = json.Unmarshal(total, totalInt)
				if err != nil {
					return fmt.Errorf("Could not unmarshal total repos count: %v", err)
				}
			}
			total, err = json.Marshal(totalInt + 1)
			if err != nil {
				return fmt.Errorf("Could not marshal total repos count: %v", err)
			}
			hsb.Put([]byte("total_repos"), total)
		}

		scoreBytes := hsb.Get([]byte("scores"))
		if scoreBytes == nil {
			scoreBytes, _ = json.Marshal([]scoreHeap{})
		}
		scores := &scoreHeap{}
		json.Unmarshal(scoreBytes, scores)

		heap.Init(scores)
		if len(*scores) > 0 && (*scores)[0].Score > resp.Average*100.0 {
			// lowest score on list is higher than this repo's score, so no need to add
			return nil
		}
		// if this repo is already in the list, remove the original entry:
		for i := range *scores {
			if strings.Compare((*scores)[i].Repo, repo) == 0 {
				heap.Remove(scores, i)
				break
			}
		}
		// now we can safely push it onto the heap
		heap.Push(scores, scoreItem{
			Repo:  repo,
			Score: resp.Average * 100.0,
			Files: resp.Files,
		})
		if len(*scores) > 100 {
			// trim heap if it's grown to over 100
			*scores = (*scores)[:100]
		}
		scoreBytes, err = json.Marshal(&scores)
		if err != nil {
			return err
		}
		err = hsb.Put([]byte("scores"), scoreBytes)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		log.Println("Bolt writing error:", err)
	}
	return
}
