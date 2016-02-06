package handlers

import (
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
			return fmt.Errorf("repo bucket not found")
		}
		return b.Put([]byte(repo), respBytes)
	})

	if err != nil {
		log.Println("Bolt writing error:", err)
	}
	return
}
