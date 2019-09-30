package handlers

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"sync"

	"github.com/dgraph-io/badger"
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
			db, err := badger.Open(badger.DefaultOptions("/tmp/badger"))
			if err != nil {
				log.Println("Failed to open badger database: ", err)
				return
			}
			defer db.Close()

			recent := &[]recentItem{}
			err = db.View(func(txn *badger.Txn) error {
				item, err := txn.Get([]byte("recent"))
				if err != nil && err != badger.ErrKeyNotFound {
					return err
				}

				if item != nil {
					err = item.Value(func(val []byte) error {
						return json.Unmarshal(val, recent)
					})

					return err
				}

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
