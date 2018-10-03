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
	"github.com/tokopedia/goreportcard/handlers"
)

var (
	repo      = flag.String("remove", "", "repo to remove from high scores")
	listDupes = flag.Bool("list-duplicates", false, "list duplicate repos in cache")
)

func deleteRepo(repo string) error {
	db, err := bolt.Open("goreportcard.db", 0755, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return fmt.Errorf("could not open bolt db: %v", err)
	}
	defer db.Close()
	return db.Update(func(tx *bolt.Tx) error {
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
			if strings.ToLower((*scores)[i].Repo) == strings.ToLower(repo) {
				log.Printf("repo %q found in high scores. Removing...", repo)
				heap.Remove(scores, i)
				found = true
				break
			}
		}

		if !found {
			log.Printf("repo %q not found in high scores. Exiting...", repo)
			return nil
		}

		scoreBytes, err := json.Marshal(&scores)
		if err != nil {
			return err
		}

		return mb.Put([]byte("scores"), scoreBytes)
	})

}

func listDuplicates() error {
	db, err := bolt.Open("goreportcard.db", 0755, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return fmt.Errorf("could not open bolt db: %v", err)
	}
	defer db.Close()
	return db.View(func(tx *bolt.Tx) error {
		repos := map[string][]string{}

		rb := tx.Bucket([]byte("repos"))
		if rb == nil {
			return fmt.Errorf("repos bucket not found")
		}

		err = rb.ForEach(func(k, v []byte) error {
			lower := strings.ToLower(string(k))
			if _, ok := repos[lower]; ok {
				repos[lower] = append(repos[lower], string(k))
				return nil
			}
			repos[lower] = []string{string(k)}

			return nil
		})

		if err != nil {
			return err
		}

		for _, v := range repos {
			if len(v) > 1 {
				for _, repo := range v {
					fmt.Println(repo)
				}
			}
		}

		return nil
	})

}

func main() {
	flag.Parse()
	if *repo == "" && !*listDupes {
		log.Println("Usage: manage_db.go [-list-duplicates] [-remove repo]")
		return
	}

	var err error
	if *repo != "" {
		err = deleteRepo(*repo)
	}

	if *listDupes {
		err = listDuplicates()
	}

	if err != nil {
		log.Fatal(err)
	}
}
