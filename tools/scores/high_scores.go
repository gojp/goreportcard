package main

import (
	"container/heap"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/boltdb/bolt"
	"github.com/gojp/goreportcard/handlers"
)

var repo = flag.String("remove", "", "repo to remove from high scores")

func main() {
	flag.Parse()
	if *repo == "" {
		log.Println("No repo provided. Usage: high_scores.go -remove [repo]")
		return
	}
	db, err := bolt.Open("goreportcard.db", 0755, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Println("Failed to open bolt database: ", err)
		return
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		mb := tx.Bucket([]byte("meta"))
		if mb == nil {
			return fmt.Errorf("high score bucket not found")
		}
		scoreBytes := mb.Get([]byte("scores"))

		scores := &handlers.ScoreHeap{}
		json.Unmarshal(scoreBytes, scores)

		heap.Init(scores)

		var found bool
		for i := range *scores {
			if strings.ToLower((*scores)[i].Repo) == strings.ToLower(*repo) {
				log.Printf("repo %q found in high scores. Removing...", *repo)
				heap.Remove(scores, i)
				found = true
				break
			}
		}

		if !found {
			log.Printf("repo %q not found in high scores. Exiting...", *repo)
			return nil
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
	})

	if err != nil {
		log.Fatal(err)
	}
}
