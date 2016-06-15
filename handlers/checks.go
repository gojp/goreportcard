package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/boltdb/bolt"
	humanize "github.com/dustin/go-humanize"
	"github.com/gojp/goreportcard/check"
)

var reBadRepo = regexp.MustCompile(`package\s([\w\/\.]+)\: exit status \d+`)

func dirName(repo string) string {
	return fmt.Sprintf("repos/src/%s", repo)
}

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
	resp.HumanizedLastRefresh = humanize.Time(resp.LastRefresh.UTC())

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
	Checks               []score   `json:"checks"`
	Average              float64   `json:"average"`
	Grade                Grade     `json:"grade"`
	Files                int       `json:"files"`
	Issues               int       `json:"issues"`
	Repo                 string    `json:"repo"`
	LastRefresh          time.Time `json:"last_refresh"`
	HumanizedLastRefresh string    `json:"humanized_last_refresh"`
}

func goGet(repo string, prevError string) error {
	log.Printf("Go getting %q...", repo)
	if err := os.Mkdir("repos", 0755); err != nil && !os.IsExist(err) {
		return fmt.Errorf("could not create dir: %v", err)
	}
	d, err := filepath.Abs("repos")
	if err != nil {
		return fmt.Errorf("could not fetch absolute path: %v", err)
	}

	os.Setenv("GOPATH", d)
	cmd := exec.Command("go", "get", "-d", "-u", repo)
	cmd.Stdout = os.Stdout
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("could not get stderr pipe: %v", err)
	}

	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("could not start command: %v", err)
	}

	b, err := ioutil.ReadAll(stderr)
	if err != nil {
		return fmt.Errorf("could not read stderr: %v", err)
	}

	err = cmd.Wait()
	errStr := string(b)

	// we don't care if there are no buildable Go source files, we just need the source on disk
	hadError := err != nil && !strings.Contains(errStr, "no buildable Go source files")

	if hadError {
		log.Println("Go get error log:", string(b))
		if errStr != prevError {
			// try again, this time deleting the cached directory, and also the
			// package that caused the error in case our cache is stale
			// (remote repository or one of its dependencices was force-pushed,
			// replaced, etc)
			err = os.RemoveAll(filepath.Join(d, "src", repo))
			if err != nil {
				return fmt.Errorf("could not delete repo: %v", err)
			}

			packageNames := reBadRepo.FindStringSubmatch(errStr)
			if len(packageNames) >= 2 {
				pkg := packageNames[1]
				fp := filepath.Clean(filepath.Join(d, "src", pkg))
				if strings.HasPrefix(fp, filepath.Join(d, "src")) {
					// if the path is prefixed with the path to our
					// cached repos, then it's safe to delete it.
					// These precautions are here so that someone can't
					// craft a malicious package name with .. in it
					// and cause us to delete our server's root directory.
					log.Println("Cleaning out rebased dependency:", fp)
					err = os.RemoveAll(fp)
					if err != nil {
						return err
					}
				}
			}
			return goGet(repo, errStr)
		}

		return fmt.Errorf("could not run go get: %v", err)
	}
	return nil
}

func newChecksResp(repo string, forceRefresh bool) (checksResp, error) {
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
	err := goGet(repo, "")
	if err != nil {
		return checksResp{}, fmt.Errorf("could not clone repo: %v", err)
	}

	dir := dirName(repo)
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
		check.IneffAssign{Dir: dir, Filenames: filenames},
		check.ErrCheck{Dir: dir, Filenames: filenames},
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

	resp := checksResp{
		Repo:                 repo,
		Files:                len(filenames),
		LastRefresh:          time.Now().UTC(),
		HumanizedLastRefresh: humanize.Time(time.Now().UTC()),
	}

	var total float64
	var totalWeight float64
	var issues = make(map[string]bool)
	for i := 0; i < len(checks); i++ {
		s := <-ch
		resp.Checks = append(resp.Checks, s)
		total += s.Percentage * s.Weight
		totalWeight += s.Weight
		for _, fs := range s.FileSummaries {
			issues[fs.Filename] = true
		}
	}
	total /= totalWeight

	sort.Sort(ByWeight(resp.Checks))
	resp.Average = total
	resp.Issues = len(issues)
	resp.Grade = grade(total * 100)

	return resp, nil
}

// ByWeight implements sorting for checks by weight descending
type ByWeight []score

func (a ByWeight) Len() int           { return len(a) }
func (a ByWeight) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByWeight) Less(i, j int) bool { return a[i].Weight > a[j].Weight }
