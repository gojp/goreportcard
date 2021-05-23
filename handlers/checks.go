package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/dgraph-io/badger/v2"
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
	return fmt.Sprintf("data/_repos/src/%s", repo)
}

func getKeyForCache(repo, branch string) []byte {
	return []byte(RepoPrefix + repo + "|branch-" + branch)
}

func getFromCache(db *badger.DB, repo, branch string) (checksResp, error) {
	// try and fetch from badger
	resp := checksResp{}
	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(getKeyForCache(repo, branch))
		if err != nil && err != badger.ErrKeyNotFound {
			return err
		}

		if item != nil {
			err = item.Value(func(val []byte) error {
				err = json.Unmarshal(val, &resp)
				if err != nil {
					return fmt.Errorf("failed to parse JSON for %q in cache", repo)
				}

				return nil
			})
		}

		if item == nil {
			log.Printf("Repo %q not found in badger cache", repo)
			return notFoundError{repo}
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
	Branch               string        `json:"branch"`
	LastRefresh          time.Time     `json:"last_refresh"`
	LastRefreshFormatted string        `json:"formatted_last_refresh"`
	LastRefreshHumanized string        `json:"humanized_last_refresh"`
}

func newChecksResp(db *badger.DB, repo, branch string, forceRefresh bool) (checksResp, error) {
	if !forceRefresh {
		resp, err := getFromCache(db, repo, branch)
		if err != nil {
			// just log the error and continue
			log.Println(err)
		} else {
			resp.Grade = check.GradeFromPercentage(resp.Average * 100) // grade is not stored for some repos, yet
			return resp, nil
		}
	}

	// fetch the repo and grade it
	repoRoot, err := download.Download(repo, branch, "data/_repos/src")
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
		Branch:               branch,
		LastRefresh:          t,
		LastRefreshFormatted: t.Format(time.UnixDate),
		LastRefreshHumanized: humanize.Time(t),
	}

	respBytes, err := json.Marshal(resp)
	if err != nil {
		return checksResp{}, fmt.Errorf("could not marshal json: %v", err)
	}

	// is this a new repo? if so, increase the count in the high scores bucket later
	isNewRepo := false
	var oldRepoBytes []byte
	err = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(getKeyForCache(repo, branch))
		if err != nil {
			return err
		}

		err = item.Value(func(val []byte) error {
			oldRepoBytes = val

			return nil
		})

		return err
	})

	if err != nil && err != badger.ErrKeyNotFound {
		log.Println("ERROR getting repo badger:", err)
	}

	isNewRepo = oldRepoBytes == nil

	// if this is a new repo, or the user force-refreshed, update the cache
	if isNewRepo || forceRefresh {
		err = db.Update(func(txn *badger.Txn) error {
			log.Printf("Saving repo %q to cache...", repo)

			// save repo to cache
			err = txn.Set(getKeyForCache(repo, branch), respBytes)
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

func saveChecksResp(db *badger.DB, checkResult *check.ChecksResult, repo, branch string) error {
	// fetch the repo and grade it
	repoRoot, errClean := download.GetRepoRoot(repo)
	if errClean != nil {
		return errClean
	}

	t := time.Now().UTC()
	resp := checksResp{
		Checks:               checkResult.Checks,
		Average:              checkResult.Average,
		Grade:                checkResult.Grade,
		Files:                checkResult.Files,
		Issues:               checkResult.Issues,
		Repo:                 repo,
		Branch:               branch,
		ResolvedRepo:         repoRoot.Repo,
		LastRefresh:          t,
		LastRefreshFormatted: t.Format(time.UnixDate),
		LastRefreshHumanized: humanize.Time(t),
	}

	respBytes, err := json.Marshal(resp)
	if err != nil {
		return fmt.Errorf("could not marshal json: %v", err)
	}

	// is this a new repo? if so, increase the count in the high scores bucket later
	isNewRepo := false
	var oldRepoBytes []byte
	err = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(getKeyForCache(repo, branch))
		if err != nil {
			return err
		}

		err = item.Value(func(val []byte) error {
			oldRepoBytes = val

			return nil
		})

		return err
	})

	if err != nil && err != badger.ErrKeyNotFound {
		log.Println("ERROR getting repo badger:", err)
	}

	isNewRepo = oldRepoBytes == nil

	// if this is a new repo, or the user force-refreshed, update the cache
	err = db.Update(func(txn *badger.Txn) error {
		log.Printf("Saving repo %q to cache...", repo)

		// save repo to cache
		err = txn.Set(getKeyForCache(repo, branch), respBytes)
		if err != nil {
			return err
		}

		return updateMetadata(txn, resp, repo, isNewRepo)
	})

	if err != nil {
		log.Println("Badger writing error:", err)
	}

	err = db.Update(func(txn *badger.Txn) error {
		return updateRecentlyViewed(txn, repo)
	})

	if err != nil {
		log.Printf("ERROR: could not update recently viewed: %v", err)
	}

	return nil
}

func (cs *checksResp) CalculateFileURLForFileSummaries() *checksResp {
	ResolvedBranch := check.GetBranchResolve(cs.Repo, cs.Branch)
	for indexCheck := range cs.Checks {
		for indexFS := range cs.Checks[indexCheck].FileSummaries {
			cs.Checks[indexCheck].FileSummaries[indexFS].FileURL =
				check.FileURL(cs.Repo, cs.Checks[indexCheck].FileSummaries[indexFS].Filename, ResolvedBranch)
		}
	}

	return cs
}
