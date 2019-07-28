package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/gojp/goreportcard/database"

	"github.com/gojp/goreportcard/handlers"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	addr      string
	redisHost string
)

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
		// at first Go Report Card only supported github.com URLs, and
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

// metrics provides functionality for monitoring the application status
type metrics struct {
	responseTimes *prometheus.SummaryVec
}

// setupMetrics creates custom Prometheus metrics for monitoring
// application statistics.
func setupMetrics() *metrics {
	m := &metrics{}
	m.responseTimes = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "response_time_ms",
			Help: "Time to serve requests",
		},
		[]string{"endpoint"},
	)

	prometheus.MustRegister(m.responseTimes)
	return m
}

// recordDuration records the length of a request from start to now
func (m metrics) recordDuration(start time.Time, name string) {
	elapsed := time.Since(start)
	m.responseTimes.WithLabelValues(name).Observe(float64(elapsed.Nanoseconds()))
	log.Printf("Served %s in %s", name, elapsed)
}

// instrument adds metric instrumentation to the handler passed in as argument
func (m metrics) instrument(path string, h http.HandlerFunc) (string, http.HandlerFunc) {
	return path, func(w http.ResponseWriter, r *http.Request) {
		defer m.recordDuration(time.Now(), path)
		h.ServeHTTP(w, r)
	}
}

func main() {
	flag.StringVar(&addr, "http", ":8000", "HTTP listen address")
	flag.StringVar(&redisHost, "redis", "", "Address of Redis server")
	flag.Parse()
	if err := os.MkdirAll("_repos/src/github.com", 0755); err != nil && !os.IsExist(err) {
		log.Fatal("ERROR: could not create repos dir: ", err)
	}

	// initialize database
	db, err := database.GetConnection(redisHost)
	if err != nil {
		log.Fatal("ERROR: could not connect to db: ", err)
	}

	m := setupMetrics()

	homeHandler := handlers.HomeHandler{DB: db}
	checkHandler := handlers.CheckHandler{DB: db}
	reportHandler := handlers.ReportHandler{DB: db}
	badgeHandler := handlers.BadgeHandler{DB: db}
	highScoresHandler := handlers.HighScoresHandler{DB: db}

	http.HandleFunc(m.instrument("/assets/", handlers.AssetsHandler))
	http.HandleFunc(m.instrument("/favicon.ico", handlers.FaviconHandler))
	http.HandleFunc(m.instrument("/checks", checkHandler.Handle))
	http.HandleFunc(m.instrument("/report/", makeHandler("report", reportHandler.Handle)))
	http.HandleFunc(m.instrument("/badge/", makeHandler("badge", badgeHandler.Handle)))
	http.HandleFunc(m.instrument("/high_scores/", highScoresHandler.Handle))
	http.HandleFunc(m.instrument("/about/", handlers.AboutHandler))
	http.HandleFunc(m.instrument("/", homeHandler.Handle))

	http.Handle("/metrics", promhttp.Handler())

	log.Printf("Running on %s ...", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
