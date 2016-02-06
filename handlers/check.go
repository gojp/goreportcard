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

	// MetaBucket is the bucket containing the names of the projects with the
	// top 100 high scores, and other meta information
	MetaBucket string = "meta"
)

// CheckHandler handles the request for checking a repo
func CheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	repo := r.FormValue("repo")
	log.Printf("Checking repo %s...", repo)

	// repoRoot, err := vcs.RepoRootForImportPath(repo, true)
	// if err != nil {
	// 	log.Println("ERROR:", err)
	// 	http.Error(w, err.Error(), 404)
	// 	return
	// }
	// repo = repoRoot.Root
	// log.Println("RepoRootForImportPath:", repo)

	if strings.ToLower(repo) == "golang/go" {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("We've decided to omit results for the Go repository because it has lots of test files that (purposely) don't pass our checks. Go gets an A+ in our books though!"))
		return
	}
	forceRefresh := r.Method != "GET" // if this is a GET request, try to fetch from cached version in boltdb first
	resp, err := newChecksResp(repo, forceRefresh)
	if err != nil {
		log.Println("ERROR: from newChecksResp: ", err)
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

	// is this a new repo? if so, increase the count in the high scores bucket later
	isNewRepo := false
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(RepoBucket))
		if b == nil {
			return fmt.Errorf("repo bucket not found")
		}
		isNewRepo = b.Get([]byte(repo)) == nil
		return nil
	})
	if err != nil {
		log.Println(err)
	}

	// if this is a new repo, or the user force-refreshed, update the cache
	if isNewRepo || forceRefresh {
		err = db.Update(func(tx *bolt.Tx) error {
			log.Printf("Saving repo %q to cache...", repo)

			b := tx.Bucket([]byte(RepoBucket))
			if b == nil {
				return fmt.Errorf("repo bucket not found")
			}

			// save repo to cache
			err = b.Put([]byte(repo), respBytes)
			if err != nil {
				return err
			}

			// fetch meta-bucket
			mb := tx.Bucket([]byte(MetaBucket))
			if mb == nil {
				return fmt.Errorf("high score bucket not found")
			}

			// update total repos count
			if isNewRepo {
				err = updateReposCount(mb, resp, repo)
				if err != nil {
					return err
				}
			}

			return updateHighScores(mb, resp, repo)
		})

		if err != nil {
			log.Println("Bolt writing error:", err)
		}
	}

	return
}

func updateHighScores(mb *bolt.Bucket, resp checksResp, repo string) error {
	// check if we need to update the high score list
	if resp.Files < 100 {
		// only repos with >= 100 files are considered for the high score list
		return nil
	}

	// start updating high score list
	scoreBytes := mb.Get([]byte("scores"))
	if scoreBytes == nil {
		scoreBytes, _ = json.Marshal([]scoreHeap{})
	}
	scores := &scoreHeap{}
	json.Unmarshal(scoreBytes, scores)

	heap.Init(scores)
	if len(*scores) > 0 && (*scores)[0].Score > resp.Average*100.0 && len(*scores) == 50 {
		// lowest score on list is higher than this repo's score, so no need to add, unless
		// we do not have 50 high scores yet
		return nil
	}
	// if this repo is already in the list, remove the original entry:
	for i := range *scores {
		if (*scores)[i].Repo == repo {
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
	if len(*scores) > 50 {
		// trim heap if it's grown to over 50
		*scores = (*scores)[1:51]
	}
	scoreBytes, err := json.Marshal(&scores)
	if err != nil {
		return err
	}
	err = mb.Put([]byte("scores"), scoreBytes)
	if err != nil {
		return err
	}

	return nil
}

func updateReposCount(mb *bolt.Bucket, resp checksResp, repo string) (err error) {
	log.Printf("New repo %q, adding to repo count...", repo)
	totalInt := 0
	total := mb.Get([]byte("total_repos"))
	if total != nil {
		err = json.Unmarshal(total, &totalInt)
		if err != nil {
			return fmt.Errorf("could not unmarshal total repos count: %v", err)
		}
	}
	totalInt++ // increase repo count
	total, err = json.Marshal(totalInt)
	if err != nil {
		return fmt.Errorf("could not marshal total repos count: %v", err)
	}
	mb.Put([]byte("total_repos"), total)
	log.Println("Repo count is now", totalInt)
	return nil
}
