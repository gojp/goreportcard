package handlers

import (
	"container/heap"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/gojp/goreportcard/download"
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
		scoreBytes, _ = json.Marshal([]ScoreHeap{})
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

func updateReposCount(mb *bolt.Bucket, repo string) (err error) {
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

type recentItem struct {
	Repo string
}

func updateRecentlyViewed(mb *bolt.Bucket, repo string) error {
	if mb == nil {
		return fmt.Errorf("meta bucket not found")
	}
	b := mb.Get([]byte("recent"))
	if b == nil {
		b, _ = json.Marshal([]recentItem{})
	}
	recent := []recentItem{}
	json.Unmarshal(b, &recent)

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
	err = mb.Put([]byte("recent"), b)
	if err != nil {
		return err
	}

	return nil
}

//func updateMetadata(tx *bolt.Tx, resp checksResp, repo string, isNewRepo bool, oldScore *float64) error {
func updateMetadata(tx *bolt.Tx, resp checksResp, repo string, isNewRepo bool) error {
	// fetch meta-bucket
	mb := tx.Bucket([]byte(MetaBucket))
	if mb == nil {
		return fmt.Errorf("high score bucket not found")
	}
	// update total repos count
	if isNewRepo {
		err := updateReposCount(mb, repo)
		if err != nil {
			return err
		}
	}

	return updateHighScores(mb, resp, repo)
}
