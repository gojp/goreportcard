package handlers

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gojp/goreportcard/check"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	mongoURL        = "mongodb://localhost:27017"
	mongoDatabase   = "goreportcard"
	mongoCollection = "reports"
)

func getFromCache(repo string) (checksResp, error) {
	// try and fetch from mongo
	session, err := mgo.Dial(mongoURL)
	if err != nil {
		return checksResp{}, fmt.Errorf("Failed to get mongo collection during GET: %v", err)
	}
	defer session.Close()
	coll := session.DB(mongoDatabase).C(mongoCollection)
	resp := checksResp{}
	err = coll.Find(bson.M{"repo": repo}).One(&resp)
	if err != nil {
		return checksResp{}, fmt.Errorf("Failed to fetch %q from mongo: %v", repo, err)
	}

	resp.LastRefresh = resp.LastRefresh.UTC()

	return resp, nil
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
		cmd := exec.Command("git", "clone", "--depth", "1", "--single-branch", url, dir)
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
			resp.Grade = grade(resp.Average * 100) // grade is not stored for some repos, yet
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
	if len(filenames) == 0 {
		return checksResp{}, fmt.Errorf("No .go files found")
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
	resp.Grade = grade(resp.Average * 100)

	return resp, nil
}
