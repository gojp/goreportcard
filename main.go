package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/gojp/goreportcard/handlers"

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

// context handles startup initialization steps, and implements the handlers.Context
// interface so that it can be passed as an extra argument to certain endpoints.
type context struct{}

func newContext() (ctxt context, err error) {
	// can add initialization here at a later stage
	return ctxt, nil
}

// Suggest returns autocomplete suggestions for the given string.
func (c context) Suggest(s string) (words []string, err error) {
	// initialize database
	var db *bolt.DB
	if db, err = initDB(); err != nil {
		return words, fmt.Errorf("ERROR: could not open bolt db: %v", err)
	}
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(handlers.RepoBucket)).Cursor()

		prefix := []byte(s)
		for k, _ := c.Seek(prefix); bytes.HasPrefix(k, prefix); k, _ = c.Next() {
			words = append(words, string(k))
		}
		return nil
	})

	return words, err
}

func (c context) wrap(fn func(http.ResponseWriter, *http.Request, handlers.Context)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, c)
	}
}

// initDB opens the bolt database file (or creates it if it does not exist), and creates
// a bucket for saving the repos, also only if it does not exist.
func initDB() (*bolt.DB, error) {
	db, err := bolt.Open(handlers.DBPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(handlers.RepoBucket))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte(handlers.MetaBucket))
		return err
	})
	return db, err
}

func main() {
	flag.Parse()
	if err := os.MkdirAll("repos/src/github.com", 0755); err != nil && !os.IsExist(err) {
		log.Fatal("ERROR: could not create repos dir: ", err)
	}

	ctxt, err := newContext()
	if err != nil {
		log.Fatal("ERROR: could not create context: ", err)
	}

	http.HandleFunc("/assets/", handlers.AssetsHandler)
	http.HandleFunc("/favicon.ico", handlers.FaviconHandler)
	http.HandleFunc("/checks", handlers.CheckHandler)
	http.HandleFunc("/report/", makeHandler("report", handlers.ReportHandler))
	http.HandleFunc("/badge/", makeHandler("badge", handlers.BadgeHandler))
	http.HandleFunc("/high_scores/", handlers.HighScoresHandler)
	http.HandleFunc("/about/", handlers.AboutHandler)
	http.HandleFunc("/suggest/", ctxt.wrap(handlers.SuggestionHandler))
	http.HandleFunc("/", handlers.HomeHandler)

	log.Printf("Running on %s ...", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
