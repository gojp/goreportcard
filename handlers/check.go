package handlers

import (
	"container/heap"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/dgraph-io/badger"
	"github.com/gojp/goreportcard/download"
)

const (
	// RepoPrefix is the badger prefix for repos
	RepoPrefix string = "repos-"
)

// CheckHandler handles the request for checking a repo
func CheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	repo, err := download.Clean(r.FormValue("repo"))
	if err != nil {
		log.Println("ERROR: from download.Clean:", err)
		http.Error(w, "Could not download the repository: "+err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Checking repo %q...", repo)

	forceRefresh := r.Method != "GET" // if this is a GET request, try to fetch from cached version in boltdb first
	_, err = newChecksResp(repo, forceRefresh)
	if err != nil {
		log.Println("ERROR: from newChecksResp:", err)
		http.Error(w, "Could not analyze the repository: "+err.Error(), http.StatusBadRequest)
		return
	}

	b, err := json.Marshal(map[string]string{"redirect": "/report/" + repo})
	if err != nil {
		log.Println("JSON marshal error:", err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func updateHighScores(txn *badger.Txn, resp checksResp, repo string) error {
	// check if we need to update the high score list
	if resp.Files < 100 {
		// only repos with >= 100 files are considered for the high score list
		return nil
	}

	var scoreBytes []byte
	// start updating high score list
	item, err := txn.Get([]byte("scores"))
	if err != nil && err != badger.ErrKeyNotFound {
		return err
	}

	if item == nil {
		scoreBytes, _ = json.Marshal([]ScoreHeap{})
	}

	if item != nil {
		err = item.Value(func(val []byte) error {
			err = json.Unmarshal(val, &scoreBytes)
			if err != nil {
				return fmt.Errorf("could not unmarshal high scores: %v", err)
			}

			return nil
		})
	}

	scores := &ScoreHeap{}
	json.Unmarshal(scoreBytes, scores)

	heap.Init(scores)
	if len(*scores) > 0 && (*scores)[0].Score > resp.Average*100.0 && len(*scores) == 50 {
		// lowest score on list is higher than this repo's score, so no need to add, unless
		// we do not have 50 high scores yet
		return nil
	}

	// if this repo is already in the list, remove the original entry:
	for i := range *scores {
		if strings.ToLower((*scores)[i].Repo) == strings.ToLower(repo) {
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

	scoreBytes, err = json.Marshal(&scores)
	if err != nil {
		return err
	}

	return txn.Set([]byte("scores"), scoreBytes)
}

func updateReposCount(txn *badger.Txn, repo string) error {
	log.Printf("New repo %q, adding to repo count...", repo)
	totalInt := 0
	item, err := txn.Get([]byte("total_repos"))
	if err != nil && err != badger.ErrKeyNotFound {
		return err
	}

	if item != nil {
		err = item.Value(func(val []byte) error {
			err = json.Unmarshal(val, &totalInt)
			if err != nil {
				return fmt.Errorf("could not unmarshal total repos count: %v", err)
			}

			return nil
		})
	}

	totalInt++ // increase repo count
	total, err := json.Marshal(totalInt)
	if err != nil {
		return fmt.Errorf("could not marshal total repos count: %v", err)
	}

	err = txn.Set([]byte("total_repos"), total)
	if err != nil {
		return err
	}
	log.Println("Repo count is now", totalInt)

	return nil
}

type recentItem struct {
	Repo string
}

func updateRecentlyViewed(txn *badger.Txn, repo string) error {
	var recent []recentItem
	item, err := txn.Get([]byte("recent"))
	if item != nil {
		item.Value(func(val []byte) error {
			return json.Unmarshal(val, &recent)
		})
	}

	// add it to the slice, if it is not in there already
	for i := range recent {
		if recent[i].Repo == repo {
			return nil
		}
	}

	recent = append(recent, recentItem{Repo: repo})
	if len(recent) > 5 {
		// trim recent if it's grown to over 5
		recent = (recent)[1:6]
	}
	b, err := json.Marshal(&recent)
	if err != nil {
		return err
	}

	return txn.Set([]byte("recent"), b)
}

func updateMetadata(txn *badger.Txn, resp checksResp, repo string, isNewRepo bool) error {
	// update total repos count
	if isNewRepo {
		err := updateReposCount(txn, repo)
		if err != nil {
			return err
		}
	}

	return updateHighScores(txn, resp, repo)
}
