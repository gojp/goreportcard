package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/boltdb/bolt"
	"github.com/gojp/goreportcard/check"

	"gopkg.in/mgo.v2"
)

const (
	dbPath     string = "goreportcard.db"
	repoBucket string = "repos"

	mongoURL        = "mongodb://127.0.0.1:27017"
	mongoDatabase   = "goreportcard"
	mongoCollection = "reports"
)

type Grade string

type score struct {
	Name          string              `json:"name"`
	Description   string              `json:"description"`
	FileSummaries []check.FileSummary `json:"file_summaries"`
	Weight        float64             `json:"weight"`
	Percentage    float64             `json:"percentage"`
}

type checksResp struct {
	Checks      []score   `json:"checks"`
	Average     float64   `json:"average"`
	Grade       Grade     `json:"grade"`
	Files       int       `json:"files"`
	Issues      int       `json:"issues"`
	Repo        string    `json:"repo"`
	LastRefresh time.Time `json:"last_refresh"`
}

// initDB opens the bolt database file (or creates it if it does not exist), and creates
// a bucket for saving the repos, also only if it does not exist.
func initDB() error {
	db, err := bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(repoBucket))
		return err
	})
	return err
}

func main() {
	// initialize bolt database
	if err := initDB(); err != nil {
		log.Fatal("ERROR: could not open bolt db: ", err)
	}

	session, err := mgo.Dial(mongoURL)
	if err != nil {
		log.Fatal("ERROR: could not get collection:", err)
	}
	defer session.Close()
	coll := session.DB(mongoDatabase).C(mongoCollection)

	db, err := bolt.Open(dbPath, 0755, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Println("Failed to open bolt database: ", err)
		return
	}
	defer db.Close()

	var repos []checksResp
	coll.Find(nil).All(&repos)

	for _, repo := range repos {
		fmt.Printf("inserting %q into bolt...\n", repo.Repo)
		err = db.Update(func(tx *bolt.Tx) error {
			bkt := tx.Bucket([]byte(repoBucket))
			if bkt == nil {
				return fmt.Errorf("repo bucket not found")
			}
			b, err := json.Marshal(repo)
			if err != nil {
				return err
			}
			return bkt.Put([]byte(repo.Repo), b)
		})
		if err != nil {
			log.Println("Bolt writing error:", err)
		}
	}
}
