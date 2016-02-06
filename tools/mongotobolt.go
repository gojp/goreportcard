package main

import (
	"container/heap"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/boltdb/bolt"
	"github.com/gojp/goreportcard/check"
	"github.com/gojp/goreportcard/handlers"

	"gopkg.in/mgo.v2"
)

const (
	dbPath     string = "goreportcard.db"
	repoBucket string = "repos"
	metaBucket string = "meta"

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
	db, err := bolt.Open(handlers.DBPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(repoBucket))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte(metaBucket))
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

			mb := tx.Bucket([]byte(metaBucket))
			if mb == nil {
				return fmt.Errorf("repo bucket not found")
			}
			updateHighScores(mb, repo, repo.Repo)

			return bkt.Put([]byte(repo.Repo), b)
		})
		if err != nil {
			log.Println("Bolt writing error:", err)
		}
	}

	err = db.Update(func(tx *bolt.Tx) error {
		mb := tx.Bucket([]byte(metaBucket))
		if mb == nil {
			return fmt.Errorf("repo bucket not found")
		}
		totalInt := len(repos)
		total, err := json.Marshal(totalInt)
		if err != nil {
			return fmt.Errorf("could not marshal total repos count: %v", err)
		}
		return mb.Put([]byte("total_repos"), total)
	})
	if err != nil {
		log.Fatal(err)
	}
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
		scoreBytes, _ = json.Marshal([]scoreHeap{})
	}
	scores := &scoreHeap{}
	json.Unmarshal(scoreBytes, scores)

	heap.Init(scores)
	if len(*scores) > 0 && (*scores)[0].Score > resp.Average*100.0 && len(*scores) == 50 {
		// lowest score on list is higher than this repo's score, so no need to add, unless
		// we do not have 50 high scores yet
		return nil
	}
	// if this repo is already in the list, remove the original entry:
	for i := range *scores {
		if (*scores)[i].Repo == repo {
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
