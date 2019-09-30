package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/dgraph-io/badger"
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

func getFromCache(db *badger.DB, repo string) (checksResp, error) {
	// try and fetch from badger
	fmt.Println("1")

	resp := checksResp{}
	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(RepoPrefix + repo))
		if err != nil && err != badger.ErrKeyNotFound {
			return err
		}

		if item == nil {
			return notFoundError{repo}
		}

		err = item.Value(func(val []byte) error {
			err = json.Unmarshal(val, &resp)
			if err != nil {
				return fmt.Errorf("failed to parse JSON for %q in cache", repo)
			}

			return nil
		})

		return err
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

func newChecksResp(db *badger.DB, repo string, forceRefresh bool) (checksResp, error) {
	if !forceRefresh {
		resp, err := getFromCache(db, repo)
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
	checkResult, err := check.Run(dirName(repo))
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

	badgeCache.Store(repo, checkResult.Grade)

	respBytes, err := json.Marshal(resp)
	if err != nil {
		return checksResp{}, fmt.Errorf("could not marshal json: %v", err)
	}

	// is this a new repo? if so, increase the count in the high scores bucket later
	isNewRepo := false
	var oldRepoBytes []byte
	err = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(RepoPrefix + repo))
		if err != nil {
			return err
		}

		err = item.Value(func(val []byte) error {
			oldRepoBytes = val

			return nil
		})

		return err
	})
	if err != nil {
		log.Println("ERROR getting repo from repo bucket:", err)
	}

	isNewRepo = oldRepoBytes == nil

	// if this is a new repo, or the user force-refreshed, update the cache
	if isNewRepo || forceRefresh {
		err = db.Update(func(txn *badger.Txn) error {
			log.Printf("Saving repo %q to cache...", repo)

			// save repo to cache
			err = txn.Set([]byte(RepoPrefix+repo), respBytes)
			if err != nil {
				return err
			}

			return updateMetadata(txn, resp, repo, isNewRepo)
		})

		if err != nil {
			log.Println("Badger writing error:", err)
		}

	}

	err = db.Update(func(txn *badger.Txn) error {
		return updateRecentlyViewed(txn, repo)
	})

	if err != nil {
		log.Printf("ERROR: could not update recently viewed: %v", err)
	}

	return resp, nil
}
