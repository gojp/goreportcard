package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/gojp/goreportcard/check"
	"gopkg.in/mgo.v2"
	"labix.org/v2/mgo/bson"
)

var (
	mongoURL        = "mongodb://localhost:27017"
	mongoDatabase   = "goreportcard"
	mongoCollection = "reports"
)

func getMongoCollection() (*mgo.Collection, error) {
	session, err := mgo.Dial(mongoURL)
	if err != nil {
		return nil, err
	}
	c := session.DB(mongoDatabase).C(mongoCollection)
	return c, nil
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving home page")
	if r.URL.Path[1:] == "" {
		http.ServeFile(w, r, "templates/home.html")
	} else {
		http.NotFound(w, r)
	}
}

func assetsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving " + r.URL.Path[1:])
	http.ServeFile(w, r, r.URL.Path[1:])
}

func orgRepoNames(url string) (string, string) {
	dir := strings.TrimSuffix(url, ".git")
	split := strings.Split(dir, "/")
	org := split[len(split)-2]
	repoName := split[len(split)-1]

	return org, repoName
}

func dirName(url string) string {
	org, repoName := orgRepoNames(url)

	return fmt.Sprintf("repos/src/github.com/%s/%s", org, repoName)
}

func clone(url string) error {
	org, _ := orgRepoNames(url)
	if err := os.Mkdir(fmt.Sprintf("repos/src/github.com/%s", org), 0755); err != nil && !os.IsExist(err) {
		return fmt.Errorf("could not create dir: %v", err)
	}
	dir := dirName(url)
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		cmd := exec.Command("timeout", "120", "git", "clone", "--depth", "1", "--single-branch", url, dir)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("could not run git clone: %v", err)
		}
	} else if err != nil {
		return fmt.Errorf("could not stat dir: %v", err)
	} else {
		cmd := exec.Command("git", "-C", dir, "pull")
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("could not pull repo: %v", err)
		}
	}

	return nil
}

type score struct {
	Name          string              `json:"name"`
	Description   string              `json:"description"`
	FileSummaries []check.FileSummary `json:"file_summaries"`
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

func getFromCache(repo string) (checksResp, error) {
	// try and fetch from mongo
	coll, err := getMongoCollection()
	if err != nil {
		return checksResp{}, fmt.Errorf("Failed to get mongo collection during GET: %v", err)
	}
	resp := checksResp{}
	err = coll.Find(bson.M{"repo": repo}).One(&resp)
	if err != nil {
		return checksResp{}, fmt.Errorf("Failed to fetch %q from mongo: %v", repo, err)
	}

	resp.LastRefresh = resp.LastRefresh.UTC()

	return resp, nil
}

// Grade represents a grade returned by the server, which is normally
// somewhere between A+ (highest) and F (lowest).
type Grade string

// The Grade constants below indicate the current available
// grades.
const (
	GradeAPlus Grade = "A+"
	GradeA           = "A"
	GradeB           = "B"
	GradeC           = "C"
	GradeD           = "D"
	GradeE           = "E"
	GradeF           = "F"
)

// getGrade is a helper for getting the grade for a percentage
func getGrade(percentage float64) Grade {
	switch true {
	case percentage > 90:
		return GradeAPlus
	case percentage > 80:
		return GradeA
	case percentage > 70:
		return GradeB
	case percentage > 60:
		return GradeC
	case percentage > 50:
		return GradeD
	case percentage > 40:
		return GradeE
	default:
		return GradeF
	}
}

func newChecksResp(repo string, forceRefresh bool) (checksResp, error) {
	url := repo
	if !strings.HasPrefix(url, "https://gojp:gojp@github.com/") {
		url = "https://gojp:gojp@github.com/" + url
	}

	if !forceRefresh {
		resp, err := getFromCache(repo)
		if err != nil {
			// just log the error and continue
			log.Println(err)
		} else {
			resp.Grade = getGrade(resp.Average * 100) // grade is not stored for some repos, yet
			return resp, nil
		}
	}

	// fetch the repo and grade it
	err := clone(url)
	if err != nil {
		return checksResp{}, fmt.Errorf("Could not clone repo: %v", err)
	}

	dir := dirName(url)
	filenames, err := check.GoFiles(dir)
	if err != nil {
		return checksResp{}, fmt.Errorf("Could not get filenames: %v", err)
	}
	checks := []check.Check{check.GoFmt{Dir: dir, Filenames: filenames},
		check.GoVet{Dir: dir, Filenames: filenames},
		check.GoLint{Dir: dir, Filenames: filenames},
		check.GoCyclo{Dir: dir, Filenames: filenames},
	}

	ch := make(chan score)
	for _, c := range checks {
		go func(c check.Check) {
			p, summaries, err := c.Percentage()
			if err != nil {
				log.Printf("ERROR: (%s) %v", c.Name(), err)
			}
			s := score{
				Name:          c.Name(),
				Description:   c.Description(),
				FileSummaries: summaries,
				Percentage:    p,
			}
			ch <- s
		}(c)
	}

	resp := checksResp{Repo: repo,
		Files:       len(filenames),
		LastRefresh: time.Now().UTC()}
	var avg float64
	var issues = make(map[string]bool)
	for i := 0; i < len(checks); i++ {
		s := <-ch
		resp.Checks = append(resp.Checks, s)
		avg += s.Percentage
		for _, fs := range s.FileSummaries {
			issues[fs.Filename] = true
		}
	}

	resp.Average = avg / float64(len(checks))
	resp.Issues = len(issues)
	resp.Grade = getGrade(resp.Average * 100)

	return resp, nil
}

func checkHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	repo := r.FormValue("repo")
	forceRefresh := r.Method != "GET" // if this is a GET request, try fetch from cached version in mongo first
	resp, err := newChecksResp(repo, forceRefresh)
	if err != nil {
		b, _ := json.Marshal(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(b)
	}

	b, err := json.Marshal(resp)
	if err != nil {
		log.Println("ERROR: could not marshal json:", err)
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(b)

	// write to mongo
	coll, err := getMongoCollection()
	if err != nil {
		log.Println("Failed to get mongo collection: ", err)
	} else {
		log.Println("Writing to mongo...")
		_, err := coll.Upsert(bson.M{"Repo": repo}, resp)
		if err != nil {
			log.Println("Mongo writing error:", err)
		}
	}
}

func reportHandler(w http.ResponseWriter, r *http.Request, org, repo string) {
	http.ServeFile(w, r, "templates/home.html")
}

func badgeURL(grade Grade) string {
	colorMap := map[Grade]string{
		GradeAPlus: "brightgreen",
		GradeA:     "brightgreen",
		GradeB:     "yellowgreen",
		GradeC:     "yellow",
		GradeD:     "orange",
		GradeE:     "red",
		GradeF:     "red",
	}
	url := fmt.Sprintf("https://img.shields.io/badge/go_report-%s-%s.svg", grade, colorMap[grade])
	return url
}

func badgeHandler(w http.ResponseWriter, r *http.Request, org, repo string) {
	name := fmt.Sprintf("%s/%s", org, repo)
	resp, err := newChecksResp(name, false)
	if err != nil {
		log.Printf("ERROR: fetching badge for %s: %v", name, err)
		http.Redirect(w, r, "https://img.shields.io/badge/go%20report-error-lightgrey.svg", http.StatusTemporaryRedirect)
		return
	}

	http.Redirect(w, r, badgeURL(resp.Grade), http.StatusTemporaryRedirect)
}

func makeHandler(name string, fn func(http.ResponseWriter, *http.Request, string, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		validPath := regexp.MustCompile(fmt.Sprintf(`^/%s/([a-zA-Z0-9\-_]+)/([a-zA-Z0-9\-_]+)$`, name))

		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[1], m[2])
	}
}

func main() {
	if err := os.MkdirAll("repos/src/github.com", 0755); err != nil && !os.IsExist(err) {
		log.Fatal("ERROR: could not create repos dir: ", err)
	}

	http.HandleFunc("/assets/", assetsHandler)
	http.HandleFunc("/checks", checkHandler)
	http.HandleFunc("/report/", makeHandler("report", reportHandler))
	http.HandleFunc("/badge/", makeHandler("badge", badgeHandler))
	http.HandleFunc("/", homeHandler)

	fmt.Println("Running on 127.0.01:8080...")
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
}
