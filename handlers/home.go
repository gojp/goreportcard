package handlers

import (
	"container/heap"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/boltdb/bolt"
)

// HomeHandler handles the homepage
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving home page")

	if r.URL.Path[1:] == "" {
		db, err := bolt.Open(DBPath, 0755, &bolt.Options{Timeout: 1 * time.Second})
		if err != nil {
			log.Println("Failed to open bolt database: ", err)
			return
		}
		defer db.Close()

		recent := &recentHeap{}
		err = db.View(func(tx *bolt.Tx) error {
			rb := tx.Bucket([]byte(MetaBucket))
			if rb == nil {
				return fmt.Errorf("meta bucket not found")
			}
			b := rb.Get([]byte("recent"))
			if b == nil {
				b, err = json.Marshal([]recentHeap{})
				if err != nil {
					return err
				}
			}
			json.Unmarshal(b, recent)

			heap.Init(recent)

			return nil
		})

		var recentRepos []string
		for _, r := range *recent {
			recentRepos = append(recentRepos, r.Repo)
		}

		t := template.Must(template.New("home.html").Delims("[[", "]]").ParseFiles("templates/home.html"))
		t.Execute(w, map[string]interface{}{"Recent": recentRepos})

		return
	} else {
		http.NotFound(w, r)
	}
}
