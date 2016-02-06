package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/gojp/goreportcard/Godeps/_workspace/src/github.com/boltdb/bolt"
	"github.com/gojp/goreportcard/handlers"
)

func makeHandler(name string, fn func(http.ResponseWriter, *http.Request, string, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		validPath := regexp.MustCompile(fmt.Sprintf(`^/%s/([a-zA-Z0-9\-_]+)/([a-zA-Z0-9\-_.]+)$`, name))

		m := validPath.FindStringSubmatch(r.URL.Path)

		if m == nil {
			http.NotFound(w, r)
			return
		}

		// catch the special period cases that github does not allow for repos
		if m[2] == "." || m[2] == ".." {
			http.NotFound(w, r)
			return
		}

		fn(w, r, m[1], m[2])
	}
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
		_, err := tx.CreateBucketIfNotExists([]byte(handlers.RepoBucket))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte(handlers.MetaBucket))
		return err
	})
	return err
}

func main() {
	if err := os.MkdirAll("repos/src/github.com", 0755); err != nil && !os.IsExist(err) {
		log.Fatal("ERROR: could not create repos dir: ", err)
	}

	// initialize database
	if err := initDB(); err != nil {
		log.Fatal("ERROR: could not open bolt db: ", err)
	}

	http.HandleFunc("/assets/", handlers.AssetsHandler)
	http.HandleFunc("/checks", handlers.CheckHandler)
	http.HandleFunc("/report/", makeHandler("report", handlers.ReportHandler))
	http.HandleFunc("/badge/", makeHandler("badge", handlers.BadgeHandler))
	http.HandleFunc("/high_scores/", handlers.HighScoresHandler)
	http.HandleFunc("/about/", handlers.AboutHandler)
	http.HandleFunc("/", handlers.HomeHandler)

	fmt.Println("Running on 127.0.0.1:8080...")
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
}
