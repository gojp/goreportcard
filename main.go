package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/YotpoLtd/goreportcard/handlers"

	"github.com/boltdb/bolt"
)

var addr = flag.String("http", ":8000", "HTTP listen address")

func makeHandler(name string, fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		validPath := regexp.MustCompile(fmt.Sprintf(`^/%s/([a-zA-Z0-9\-_\/\.]+)$`, name))

		m := validPath.FindStringSubmatch(r.URL.Path)

		if m == nil {
			http.NotFound(w, r)
			return
		}
		if len(m) < 1 || m[1] == "" {
			http.Error(w, "Please enter a repository", http.StatusBadRequest)
			return
		}

		repo := m[1]

		// for backwards-compatibility, we must support URLs formatted as
		//   /report/[org]/[repo]
		// and they will be assumed to be github.com URLs. This is because
		// at first Go Report Card only supported github.com URLs, and assumed
		// took only the org name and repo name as parameters. This is no longer the
		// case, but we do not want external links to break.
		oldFormat := regexp.MustCompile(fmt.Sprintf(`^/%s/([a-zA-Z0-9\-_]+)/([a-zA-Z0-9\-_]+)$`, name))
		m2 := oldFormat.FindStringSubmatch(r.URL.Path)
		if m2 != nil {
			// old format is being used
			repo = "github.com/" + repo
			log.Printf("Assuming intended repo is %q, redirecting", repo)
			http.Redirect(w, r, fmt.Sprintf("/%s/%s", name, repo), http.StatusMovedPermanently)
			return
		}

		fn(w, r, repo)
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
	flag.Parse()
	if err := os.MkdirAll("repos/src/github.com", 0755); err != nil && !os.IsExist(err) {
		log.Fatal("ERROR: could not create repos dir: ", err)
	}

	// initialize database
	if err := initDB(); err != nil {
		log.Fatal("ERROR: could not open bolt db: ", err)
	}

	http.HandleFunc("/assets/", handlers.AssetsHandler)
	http.HandleFunc("/favicon.ico", handlers.FaviconHandler)
	http.HandleFunc("/checks", handlers.CheckHandler)
	http.HandleFunc("/report/", makeHandler("report", handlers.ReportHandler))
	http.HandleFunc("/badge/", makeHandler("badge", handlers.BadgeHandler))
	http.HandleFunc("/high_scores/", handlers.HighScoresHandler)
	http.HandleFunc("/about/", handlers.AboutHandler)
	http.HandleFunc("/", handlers.HomeHandler)

	log.Printf("Running on %s ...", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
