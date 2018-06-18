package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/boltdb/bolt"
	"github.com/dustin/go-humanize"
	"github.com/gojp/goreportcard/check"
	"github.com/gojp/goreportcard/download"
)

type notFoundError struct {
	repo string
}

func (n notFoundError) Error() string {
	return fmt.Sprintf("%q not found in cache", n.repo)
}

func dirName(repo string) string {
	return fmt.Sprintf("_repos/src/%s", repo)
}

func getFromCache(repo string) (checksResp, error) {
	// try and fetch from boltdb
	db, err := bolt.Open(DBPath, 0600, &bolt.Options{Timeout: 3 * time.Second})
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
			return notFoundError{repo}
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
	resp.LastRefreshFormatted = resp.LastRefresh.Format(time.UnixDate)
	resp.LastRefreshHumanized = humanize.Time(resp.LastRefresh.UTC())

	return resp, nil
}

type score struct {
	Name          string              `json:"name"`
	Description   string              `json:"description"`
	FileSummaries []check.FileSummary `json:"file_summaries"`
	Weight        float64             `json:"weight"`
	Percentage    float64             `json:"percentage"`
	Error         string              `json:"error"`
}

type checksResp struct {
	Checks               []score   `json:"checks"`
	Average              float64   `json:"average"`
	Grade                Grade     `json:"grade"`
	Files                int       `json:"files"`
	Issues               int       `json:"issues"`
	Repo                 string    `json:"repo"`
	ResolvedRepo         string    `json:"resolvedRepo"`
	LastRefresh          time.Time `json:"last_refresh"`
	LastRefreshFormatted string    `json:"formatted_last_refresh"`
	LastRefreshHumanized string    `json:"humanized_last_refresh"`
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
	repoRoot, err := download.Download(repo, "_repos/src")
	if err != nil {
		return checksResp{}, fmt.Errorf("could not clone repo: %v", err)
	}

	repo = repoRoot.Root

	dir := dirName(repo)
	filenames, skipped, err := check.GoFiles(dir)
	if err != nil {
		return checksResp{}, fmt.Errorf("could not get filenames: %v", err)
	}
	if len(filenames) == 0 {
		return checksResp{}, fmt.Errorf("no .go files found")
	}

	err = check.RenameFiles(skipped)
	if err != nil {
		log.Println("Could not remove files:", err)
	}
	defer check.RevertFiles(skipped)

	checks := []check.Check{
		check.GoFmt{Dir: dir, Filenames: filenames},
		check.GoVet{Dir: dir, Filenames: filenames},
		check.GoLint{Dir: dir, Filenames: filenames},
		check.GoCyclo{Dir: dir, Filenames: filenames},
		check.License{Dir: dir, Filenames: []string{}},
		check.Misspell{Dir: dir, Filenames: filenames},
		check.IneffAssign{Dir: dir, Filenames: filenames},
		// check.ErrCheck{Dir: dir, Filenames: filenames}, // disable errcheck for now, too slow and not finalized
	}

	ch := make(chan score)
	for _, c := range checks {
		go func(c check.Check) {
			p, summaries, err := c.Percentage()
			errMsg := ""
			if err != nil {
				log.Printf("ERROR: (%s) %v", c.Name(), err)
				errMsg = err.Error()
			}
			s := score{
				Name:          c.Name(),
				Description:   c.Description(),
				FileSummaries: summaries,
				Weight:        c.Weight(),
				Percentage:    p,
				Error:         errMsg,
			}
			ch <- s
		}(c)
	}

	t := time.Now().UTC()
	resp := checksResp{
		Repo:                 repo,
		ResolvedRepo:         repoRoot.Repo,
		Files:                len(filenames),
		LastRefresh:          t,
		LastRefreshFormatted: t.Format(time.UnixDate),
		LastRefreshHumanized: humanize.Time(t),
	}

	var total, totalWeight float64
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

	err = updateCache(repo, resp, forceRefresh)
	if err != nil {
		return checksResp{}, err
	}

	return resp, nil
}

func updateCache(repo string, resp checksResp, forceRefresh bool) error {
	// write to boltdb
	db, err := bolt.Open(DBPath, 0755, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return fmt.Errorf("could not open bolt db: %v", err)
	}
	defer db.Close()

	// is this a new repo? if so, increase the count in the high scores bucket later
	isNewRepo := isNewRepo(db, repo)

	// if this is a new repo, or the user force-refreshed, update the cache
	if isNewRepo || forceRefresh {
		err := db.Update(func(tx *bolt.Tx) error {
			log.Printf("Saving repo %q to cache...", repo)

			b := tx.Bucket([]byte(RepoBucket))
			if b == nil {
				return fmt.Errorf("repo bucket not found")
			}

			respBytes, err := json.Marshal(resp)
			if err != nil {
				return fmt.Errorf("could not marshal json: %v", err)
			}

			// save repo to cache
			err = b.Put([]byte(repo), respBytes)
			if err != nil {
				return err
			}

			return updateMetadata(tx, resp, repo, isNewRepo)
		})
		if err != nil {
			log.Println("Bolt writing error:", err)
		}
	}

	db.Update(func(tx *bolt.Tx) error {
		// fetch meta-bucket
		mb := tx.Bucket([]byte(MetaBucket))
		return updateRecentlyViewed(mb, repo)
	})
	return nil
}

func isNewRepo(db *bolt.DB, repo string) bool {
	var oldRepoBytes []byte
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(RepoBucket))
		if b == nil {
			return fmt.Errorf("repo bucket not found")
		}
		oldRepoBytes = b.Get([]byte(repo))
		return nil
	})
	if err != nil {
		log.Println("ERROR getting repo from repo bucket:", err)
		return false
	}

	return oldRepoBytes == nil
}

// ByWeight implements sorting for checks by weight descending
type ByWeight []score

func (a ByWeight) Len() int           { return len(a) }
func (a ByWeight) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByWeight) Less(i, j int) bool { return a[i].Weight > a[j].Weight }
