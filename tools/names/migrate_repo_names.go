package main

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/boltdb/bolt"
	"github.com/gojp/goreportcard/handlers"
)

const (
	dbPath     string = "goreportcard.db"
	repoBucket string = "repos"
	metaBucket string = "meta"
)

func main() {
	oldFormat := regexp.MustCompile(`^([a-zA-Z0-9\-_]+)/([a-zA-Z0-9\-_]+)$`)

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
		toDelete := []string{}
		rb.ForEach(func(k, v []byte) error {
			sk := string(k)
			m := oldFormat.FindStringSubmatch(sk)
			if m != nil {
				err = rb.Put([]byte("github.com/"+sk), v)
				if err != nil {
					return err
				}
				toDelete = append(toDelete, string(v))
			}
			return nil
		})
		for i := range toDelete {
			err = rb.Delete([]byte(toDelete[i]))
			if err != nil {
				return err
			}
		}

		// finally update the high scores
		mb := tx.Bucket([]byte(metaBucket))
		if mb == nil {
			return fmt.Errorf("meta bucket not found")
		}

		scoreBytes := mb.Get([]byte("scores"))
		if scoreBytes == nil {
			scoreBytes, _ = json.Marshal([]scoreHeap{})
		}
		scores := &scoreHeap{}
		json.Unmarshal(scoreBytes, scores)
		for i := range *scores {
			m := oldFormat.FindStringSubmatch((*scores)[i].Repo)
			if m != nil {
				(*scores)[i] = scoreItem{
					Repo:  "github.com/" + (*scores)[i].Repo,
					Score: (*scores)[i].Score,
					Files: (*scores)[i].Files,
				}
			}
		}
		scoreBytes, err = json.Marshal(scores)
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

type scoreItem struct {
	Repo  string  `json:"repo"`
	Score float64 `json:"score"`
	Files int     `json:"files"`
}

// An scoreHeap is a min-heap of ints.
type scoreHeap []scoreItem

func (h scoreHeap) Len() int           { return len(h) }
func (h scoreHeap) Less(i, j int) bool { return h[i].Score < h[j].Score }
func (h scoreHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *scoreHeap) Push(x interface{}) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(scoreItem))
}

func (h *scoreHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
