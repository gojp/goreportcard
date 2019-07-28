package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gojp/goreportcard/database"

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

func getFromCache(db database.Database, repo string) (checksResp, error) {
	resp := checksResp{}
	v, err := db.GetRepo(repo)
	if err != nil {
		return resp, err
	}
	if v == "" {
		return resp, errors.New("repo not found in cache")
	}

	err = json.Unmarshal([]byte(v), &resp)
	if err != nil {
		return resp, fmt.Errorf("failed to parse JSON for %q in cache", repo)
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

func newChecksResp(db database.Database, repo string, forceRefresh bool) (checksResp, error) {
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

	respBytes, err := json.Marshal(resp)
	if err != nil {
		return checksResp{}, fmt.Errorf("could not marshal json: %v", err)
	}

	cached, err := db.GetRepo(repo)
	if err != nil {
		return checksResp{}, fmt.Errorf("could not load from database: %v", err)
	}
	isNewRepo := cached == ""

	// if this is a new repo, or the user force-refreshed, update the cache
	if isNewRepo || forceRefresh {
		db.SetRepo(repo, string(respBytes))
		if err != nil {
			log.Println("db writing error when calling SetRepo:", err)
		}
		err = updateMetadata(db, resp, repo)
		if err != nil {
			log.Println("db writing error when calling updateMetadata:", err)
		}
	}

	err = updateRecentlyViewed(db, repo)
	if err != nil {
		log.Println("db writing error when calling updateRecentlyViewed:", err)
	}

	return resp, nil
}
