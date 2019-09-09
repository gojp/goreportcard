package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/boltdb/bolt"
)

var cache struct {
	items []string
	mux   sync.Mutex
	count int
}

// HomeHandler handles the homepage
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path[1:] == "" {
		var recentRepos = []string{}

		cache.mux.Lock()
		cache.count++
		defer cache.mux.Unlock()

		if cache.count < 100 && len(cache.items) == 5 {
			recentRepos = cache.items
			log.Println("Fetching recent repos from cache...")
		} else {
			log.Println("Updating recent repos cache...")
			db, err := bolt.Open(DBPath, 0755, &bolt.Options{Timeout: 1 * time.Second})
			if err != nil {
				log.Println("Failed to open bolt database: ", err)
				return
			}
			defer db.Close()

			recent := &[]recentItem{}
			err = db.View(func(tx *bolt.Tx) error {
				rb := tx.Bucket([]byte(MetaBucket))
				if rb == nil {
					return fmt.Errorf("meta bucket not found")
				}
				b := rb.Get([]byte("recent"))
				if b == nil {
					b, err = json.Marshal([]recentItem{})
					if err != nil {
						return err
					}
				}
				json.Unmarshal(b, recent)

				return nil
			})

			recentRepos = make([]string, len(*recent))
			var j = len(*recent) - 1
			for _, r := range *recent {
				recentRepos[j] = r.Repo
				j--
			}

			cache.items = recentRepos
			cache.count = 0
		}

		t := template.Must(template.New("home.html").Delims("[[", "]]").ParseFiles("templates/home.html", "templates/footer.html"))
		t.Execute(w, map[string]interface{}{
			"Recent":               recentRepos,
			"google_analytics_key": googleAnalyticsKey,
		})

		return
	}

	errorHandler(w, r, http.StatusNotFound)
}
