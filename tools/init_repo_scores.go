package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/boltdb/bolt"
	"github.com/gojp/goreportcard/handlers"
)

const (
	dbPath     string = "goreportcard.db"
	repoBucket string = "repos"
	metaBucket string = "meta"
)

type checksResp struct {
	Repo    string  `json:"repo"`
	Average float64 `json:"average"`
}

func main() {
	db, err := bolt.Open(handlers.DBPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	stats := make([]int, 101, 101)
	err = db.Update(func(tx *bolt.Tx) error {
		rb := tx.Bucket([]byte(repoBucket))
		if rb == nil {
			return fmt.Errorf("repo bucket not found")
		}
		rb.ForEach(func(k, v []byte) error {
			resp := checksResp{}
			err := json.Unmarshal(v, &resp)
			if err != nil {
				return err
			}
			log.Printf("%s - %.2f", resp.Repo, resp.Average*100)
			stats[int(resp.Average*100)]++
			return nil
		})

		// finally update the stats
		mb := tx.Bucket([]byte(metaBucket))
		if mb == nil {
			return fmt.Errorf("meta bucket not found")
		}
		var statsBytes []byte
		statsBytes, err = json.Marshal(stats)
		if err != nil {
			return err
		}
		log.Printf("Stats: %v", stats)
		err = mb.Put([]byte("stats"), statsBytes)
		if err != nil {
			return err
		}

		tResp := mb.Get([]byte("stats"))
		tStats := []int{}
		err = json.Unmarshal(tResp, &tStats)
		if err != nil {
			return err
		}
		log.Printf("Stats Confirmation: %v", tStats)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}
