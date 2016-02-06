package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/boltdb/bolt"
	"github.com/gojp/goreportcard/check"
)

func getFromCache(repo string) (checksResp, error) {
	// try and fetch from boltdb
	db, err := bolt.Open(DBPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return checksResp{}, fmt.Errorf("failed to open bolt database during GET: %v", err)
	}
	defer db.Close()

	resp := checksResp{}
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(RepoBucket))
		if b == nil {
			return errors.New("No repo bucket")
		}
		cached := b.Get([]byte(repo))
		if cached == nil {
			return fmt.Errorf("%q not found in cache", repo)
		}

		err = json.Unmarshal(cached, &resp)
		if err != nil {
			return fmt.Errorf("failed to parse JSON for %q in cache", repo)
		}
		return nil
	})

	if err != nil {
		return resp, err
	}

	resp.LastRefresh = resp.LastRefresh.UTC()
	return resp, nil
}

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
	d, err := filepath.Abs("repos")
	if err != nil {
		return fmt.Errorf("could not fetch absolute path: %v", err)
	}
	os.Setenv("GOPATH", d)
	dir := dirName(url)
	_, err = os.Stat(dir)
	if os.IsNotExist(err) {
		cmd := exec.Command("go", "get", "-u", url)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("could not run go get: %v", err)
		}

		return nil
	}
	if err != nil {
		return fmt.Errorf("could not stat dir: %v", err)
	}

	// cmd := exec.Command("git", "-C", dir, "fetch", "origin", "master")
	// if err := cmd.Run(); err != nil {
	// 	return fmt.Errorf("could not fetch master branch: %v", err)
	// }
	// cmd = exec.Command("git", "-C", dir, "reset", "--hard", "origin/master")
	// if err = cmd.Run(); err != nil {
	// 	return fmt.Errorf("could not reset origin/master: %v", err)
	// }

	return nil
}

func newChecksResp(repo string, forceRefresh bool) (checksResp, error) {
	url := repo
	if strings.Count(repo, ".") == 0 {
		url = "github.com/" + url
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
		return checksResp{}, fmt.Errorf("could not clone repo: %v", err)
	}

	dir := dirName(url)
	filenames, err := check.GoFiles(dir)
	if err != nil {
		return checksResp{}, fmt.Errorf("could not get filenames: %v", err)
	}
	if len(filenames) == 0 {
		return checksResp{}, fmt.Errorf("no .go files found")
	}
	checks := []check.Check{
		check.GoFmt{Dir: dir, Filenames: filenames},
		check.GoVet{Dir: dir, Filenames: filenames},
		check.GoLint{Dir: dir, Filenames: filenames},
		check.GoCyclo{Dir: dir, Filenames: filenames},
		check.License{Dir: dir, Filenames: []string{}},
		check.Misspell{Dir: dir, Filenames: filenames},
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
				Weight:        c.Weight(),
				Percentage:    p,
			}
			ch <- s
		}(c)
	}

	resp := checksResp{Repo: repo,
		Files:       len(filenames),
		LastRefresh: time.Now().UTC()}
	var total float64
	var issues = make(map[string]bool)
	for i := 0; i < len(checks); i++ {
		s := <-ch
		resp.Checks = append(resp.Checks, s)
		total += s.Percentage * s.Weight
		for _, fs := range s.FileSummaries {
			issues[fs.Filename] = true
		}
	}

	resp.Average = total
	resp.Issues = len(issues)
	resp.Grade = grade(total * 100)

	return resp, nil
}
