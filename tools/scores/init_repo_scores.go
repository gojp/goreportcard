package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/boltdb/bolt"
	"github.com/gojp/goreportcard/handlers"
)

const repoBucket = "repos"

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
			return nil
		})

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}
