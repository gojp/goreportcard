package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gojp/goreportcard/database"

	"github.com/boltdb/bolt"
	"github.com/gojp/goreportcard/download"
)

type CheckHandler struct {
	DB database.Database
}

// Handle handles the request for checking a repo
func (c *CheckHandler) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	repo, err := download.Clean(r.FormValue("repo"))
	if err != nil {
		log.Println("ERROR: from download.Clean:", err)
		http.Error(w, "Could not download the repository: "+err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Checking repo %q...", repo)

	forceRefresh := r.Method != "GET" // if this is a GET request, try to fetch from cached version in boltdb first
	_, err = newChecksResp(c.DB, repo, forceRefresh)
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

func updateHighScores(db database.Database, resp checksResp, repo string) error {
	// check if we need to update the high score list
	if resp.Files < 100 {
		// only repos with >= 100 files are considered for the high score list
		return nil
	}
	return db.SetScore(repo, int(resp.Average*100))
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

func updateRecentlyViewed(db database.Database, repo string) error {
	return db.SetRecentlyViewed(repo)
}

func updateMetadata(db database.Database, resp checksResp, repo string) error {
	return updateHighScores(db, resp, repo)
}
