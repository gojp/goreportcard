package httpapi

import (
	"container/heap"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gojp/goreportcard/internal/linter"
	"github.com/gojp/goreportcard/internal/model"
	"github.com/gojp/goreportcard/internal/repository"
	vcshelper "github.com/gojp/goreportcard/internal/vcs-helper"

	"github.com/dustin/go-humanize"
	"github.com/pkg/errors"
	"github.com/yeqown/log"
)

var (
	downloader vcshelper.IDownloader
)

// Init linter package
func Init() {
	// home, _ := os.UserHomeDir()
	// downloader = vcshelper.NewGitDownloader([]*vcshelper.SSHPrikeyConfig{
	// 	{
	// 		Host:       "github.com",
	// 		PrikeyPath: filepath.Join(home, ".ssh", "id_rsa"),
	// 		Prefix:     "git",
	// 	},
	// })

	downloader = vcshelper.NewBuiltinToolVCS(map[string]string{
		"github.com":        "git",
		"git.medlinker.com": "medgit",
	})
}

// TODO: add comments, adjust logic here
func dolint(repoName string, forceRefresh bool) (model.LintResult, error) {
	log.Debugf("dolint called with repoName=%s, forceRefresh=%v", repoName, forceRefresh)
	if !forceRefresh {
		resp, err := loadLintResult(repoName)
		if err == nil {
			return *resp, nil
		}
		// just log the error and continue
		log.Warnf("dolint failed to loadLintResult, err=%v", err)
	}

	// fetch the repoName and grade it
	root, err := downloader.Download(repoName, "_repos/src")
	if err != nil {
		return model.LintResult{}, fmt.Errorf("could not clone repoName: %v", err)
	}
	log.Infof("dolint download repoName=%s finished", repoName)

	result, err := linter.Lint(root)
	if err != nil {
		return model.LintResult{}, err
	}

	t := time.Now().UTC()
	lintResult := model.LintResult{
		Checks:               result.Scores,
		Average:              result.Average,
		Grade:                result.Grade,
		Files:                result.Files,
		Issues:               result.Issues,
		Repo:                 repoName,
		ResolvedRepo:         repoName,
		LastRefresh:          t,
		LastRefreshFormatted: t.Format(time.UnixDate),
		LastRefreshHumanized: humanize.Time(t),
	}

	badgeCache.Store(repoName, result.Grade)

	var (
		isNewRepo bool // current repoName is first encounter with goreportcard
		key       = lintResultKey(repoName)
	)

	v, err := repository.GetRepo().Get(key)
	if err != nil {
		log.Debugf("dolint failed to getting lint result, key=%s, err=%v", key, err)
	}
	isNewRepo = (len(v) == 0 || errors.Cause(err) == repository.ErrKeyNotFound)

	// if this is a new repo, or the user force-refreshed, update the cache
	// TODO: make update as a function
	if isNewRepo || forceRefresh {
		if err = updateLintResult(repoName, lintResult); err != nil {
			log.Errorf("dolint update lintResult failed key=%s, err=%v", key, err)
		}
		log.Debugf("dolint updateLintResult success")
	}

	if err := updateMetadata(lintResult, repoName, isNewRepo); err != nil {
		log.Errorf("dolint.updateMetadata failed: err=%v", err)
	}

	return lintResult, nil
}

// lintResultKey . to generate db.Key of lint result
func lintResultKey(repoName string) []byte {
	return []byte("repos-" + repoName)
}

// loadLintResult .
// TODO: add comments
func loadLintResult(repoName string) (*model.LintResult, error) {
	key := lintResultKey(repoName)
	data, err := repository.GetRepo().Get(key)
	if err != nil {
		return nil, err
	}

	// TRUE: hit cache
	resp := new(model.LintResult)
	if err = json.Unmarshal(data, resp); err != nil {
		return nil, err
	}

	resp.LastRefresh = resp.LastRefresh.UTC()
	resp.LastRefreshFormatted = resp.LastRefresh.Format(time.UnixDate)
	resp.LastRefreshHumanized = humanize.Time(resp.LastRefresh.UTC())
	resp.Grade = model.GradeFromPercentage(resp.Average * 100) // grade is not stored for some repos, yet
	return resp, nil
}

// updateLintResult .
// TODO: add more comments
func updateLintResult(repoName string, result model.LintResult) error {
	data, err := json.Marshal(result)
	if err != nil {
		return errors.Wrap(err, "updateLintResult.jsonMarshal")
	}

	key := lintResultKey(repoName)
	if err = repository.GetRepo().Update(key, data); err != nil {
		return errors.Wrap(err, "updateLintResult.Update")
	}

	return nil
}

type recentItem struct {
	Repo string
}

var (
	_recentKey   = []byte("recent")
	_scoreKey    = []byte("scores")
	_reposCntKey = []byte("total_repos")
)

// updateRecentlyViewed .
func updateRecentlyViewed(repoName string) (err error) {
	var (
		items []recentItem
		_repo = repository.GetRepo()
		d     []byte
	)

	if d, err = _repo.Get(_recentKey); err != nil &&
		errors.Cause(err) != repository.ErrKeyNotFound {
		return errors.Wrap(err, "updateRecentlyViewed.repo.Get")
	}

	if len(d) != 0 {
		err = json.Unmarshal(d, &items)
		if err != nil {
			return errors.Wrap(err, "updateRecentlyViewed.jsonUnmarshal")
		}
	}

	// add it to the slice, if it is not in there already
	for _, v := range items {
		if v.Repo == repoName {
			log.Infof("updateRecentlyViewed has exists repoName=%s, so skipped", repoName)
			return
		}
	}

	items = append(items, recentItem{Repo: repoName})
	if len(items) > 5 {
		// trim recent if it's grown to over 5
		items = (items)[1:6]
	}
	d, err = json.Marshal(&items)
	if err != nil {
		return errors.Wrap(err, "updateRecentlyViewed.jsonMarshal")
	}

	log.Debugf("updateRecentlyViewed will save key=%s, v=%s", _recentKey, d)
	return _repo.Update(_recentKey, d)
}

// loadRecentlyViewed .
func loadRecentlyViewed() ([]recentItem, error) {
	var (
		items []recentItem
		_repo = repository.GetRepo()
	)

	d, err := _repo.Get(_recentKey)
	if err != nil {
		return nil, err
	}

	log.Debugf("loadRecentlyViewed got d=%s", d)
	if err = json.Unmarshal(d, &items); err != nil {
		return nil, err
	}

	return items, nil
}

// loadHighScores .
func loadHighScores() (scores ScoreHeap, err error) {
	var (
		_repo = repository.GetRepo()
		d     []byte
	)

	d, err = _repo.Get(_scoreKey)
	if err != nil {
		// if there's nothing to show, make it default empty
		if errors.Cause(err) == repository.ErrKeyNotFound {
			return scores, nil
		}
		err = errors.Wrap(err, "loadHighScores.repoGet")

		return
	}

	if err = json.Unmarshal(d, &scores); err != nil {
		err = errors.Wrap(err, "loadHighScores.jsonUnmarshal")
		return
	}
	return
}

// updateHighScores .
func updateHighScores(result model.LintResult, repoName string) (err error) {
	var (
		_repo = repository.GetRepo()
		d     []byte
	)

	// check if we need to update the high score list
	// only repos with 100+ files are considered for the high score list
	// TODO: make this as configable
	if result.Files < 10 {
		return nil
	}

	if d, err = _repo.Get(_scoreKey); err != nil &&
		errors.Cause(err) != repository.ErrKeyNotFound {

		return
	}

	var scores = new(ScoreHeap)
	if len(d) != 0 {
		if err = json.Unmarshal(d, scores); err != nil {
			return err
		}
	}

	if len(*scores) > 0 && (*scores)[0].Score > result.Average*100.0 && len(*scores) == 50 {
		// lowest score on list is higher than this repo's score, so no need to add, unless
		// we do not have 50 high scores yet
		return nil
	}

	// if this repo is already in the list, remove the original entry:
	for idx, v := range *scores {
		if strings.EqualFold(v.Repo, repoName) {
			heap.Remove(scores, idx)
			break
		}
	}

	// now we can safely push it onto the heap
	heap.Init(scores)
	heap.Push(scores, scoreItem{
		Repo:  repoName,
		Score: result.Average * 100.0,
		Files: result.Files,
	})

	if len(*scores) > 50 {
		// trim heap if it's grown to over 50
		*scores = (*scores)[1:51]
	}

	// save back
	if d, err = json.Marshal(&scores); err != nil {
		return err
	}

	return _repo.Update(_scoreKey, d)
}

// loadReposCount .
func loadReposCount() (cnt int, err error) {
	var (
		_repo = repository.GetRepo()
		d     []byte
	)

	d, err = _repo.Get(_reposCntKey)
	if err != nil && errors.Cause(err) != repository.ErrKeyNotFound {
		return
	}

	if err = json.Unmarshal(d, &cnt); err != nil {

		return
	}
	return
}

// only new can update
func incrReposCnt(repoName string) (err error) {
	log.Infof("New repo %q, adding to repo count...", repoName)
	var (
		_repo = repository.GetRepo()
		d     []byte
		cnt   int
	)

	// load and unmarshal
	if d, err = _repo.Get(_reposCntKey); err != nil &&
		errors.Cause(err) != repository.ErrKeyNotFound {

		return err
	}
	if len(d) != 0 {
		if err := json.Unmarshal(d, &cnt); err != nil {
			return err
		}
	}

	cnt++
	if d, err = json.Marshal(cnt); err != nil {
		return err
	}
	if err = _repo.Update(_reposCntKey, d); err != nil {
		return err
	}

	return nil
}

// updateMetadata to record some data of goreportcard
func updateMetadata(result model.LintResult, repoName string, isNewRepo bool) (err error) {
	// increase repos count
	if isNewRepo {
		if err = incrReposCnt(repoName); err != nil {
			log.Errorf("updateMetadata.incrReposCnt failed: err=%v", err)
		}
	}
	if err = updateRecentlyViewed(repoName); err != nil {
		log.Errorf("updateMetadata.updateRecentlyViewed failed: err=%v", err)
	}
	if err = updateHighScores(result, repoName); err != nil {
		log.Errorf("updateMetadata.updateHighScores failed: err=%v", err)
	}
	log.Debugf("updateMetadata success")

	return
}
