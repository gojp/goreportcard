package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
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

type checksResp struct {
	Checks               []check.Score `json:"checks"`
	Average              float64       `json:"average"`
	Grade                check.Grade   `json:"grade"`
	Files                int           `json:"files"`
	Issues               int           `json:"issues"`
	Repo                 string        `json:"repo"`
	ResolvedRepo         string        `json:"resolvedRepo"`
	LastRefresh          time.Time     `json:"last_refresh"`
	LastRefreshFormatted string        `json:"formatted_last_refresh"`
	LastRefreshHumanized string        `json:"humanized_last_refresh"`
}

func newChecksResp(repo string, forceRefresh bool) (checksResp, error) {
	if !forceRefresh {
		resp, err := getFromCache(repo)
		if err != nil {
			// just log the error and continue
			log.Println(err)
		} else {
			resp.Grade = check.GradeFromPercentage(resp.Average * 100) // grade is not stored for some repos, yet
			return resp, nil
		}
	}

	// fetch the repo and grade it
	repoRoot, err := download.Download(repo, "_repos/src")
	if err != nil {
		return checksResp{}, fmt.Errorf("could not clone repo: %v", err)
	}

	repo = repoRoot.Root
	checkResult, err := check.CheckDir(dirName(repo))
	if err != nil {
		return checksResp{}, err
	}

	t := time.Now().UTC()
	resp := checksResp{
		Checks:               checkResult.Checks,
		Average:              checkResult.Average,
		Grade:                checkResult.Grade,
		Files:                checkResult.Files,
		Issues:               checkResult.Issues,
		Repo:                 repo,
		ResolvedRepo:         repoRoot.Repo,
		LastRefresh:          t,
		LastRefreshFormatted: t.Format(time.UnixDate),
		LastRefreshHumanized: humanize.Time(t),
	}

	respBytes, err := json.Marshal(resp)
	if err != nil {
		return checksResp{}, fmt.Errorf("could not marshal json: %v", err)
	}

	// write to boltdb
	db, err := bolt.Open(DBPath, 0755, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return checksResp{}, fmt.Errorf("could not open bolt db: %v", err)
	}
	defer db.Close()

	// is this a new repo? if so, increase the count in the high scores bucket later
	isNewRepo := false
	var oldRepoBytes []byte
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(RepoBucket))
		if b == nil {
			return fmt.Errorf("repo bucket not found")
		}
		oldRepoBytes = b.Get([]byte(repo))
		return nil
	})
	if err != nil {
		log.Println("ERROR getting repo from repo bucket:", err)
	}

	isNewRepo = oldRepoBytes == nil

	// if this is a new repo, or the user force-refreshed, update the cache
	if isNewRepo || forceRefresh {
		err = db.Update(func(tx *bolt.Tx) error {
			log.Printf("Saving repo %q to cache...", repo)

			b := tx.Bucket([]byte(RepoBucket))
			if b == nil {
				return fmt.Errorf("repo bucket not found")
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

	return resp, nil
}
