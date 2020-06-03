package linter

import (
	"fmt"
	"sort"

	"github.com/gojp/goreportcard/internal/model"
	"github.com/yeqown/log"
)

// ILinter describes what methods various checks (gofmt, go lint, etc.)
// should implement
type ILinter interface {
	// Name of ILinter
	Name() string

	// Description of ILinter
	Description() string

	// Weight of ILinter to calc score
	Weight() float64

	// Percentage returns the passing percentage of the check,
	// as well as a map of filename to output
	Percentage() (float64, []model.FileSummary, error)
}

// getLinters . load all linters to run
func getLinters(dir string, filenames []string) []ILinter {
	return []ILinter{
		GoFmt{Dir: dir, Filenames: filenames},       //
		GoVet{Dir: dir, Filenames: filenames},       //
		GoLint{Dir: dir, Filenames: filenames},      //
		GoCyclo{Dir: dir, Filenames: filenames},     //
		License{Dir: dir, Filenames: []string{}},    //
		Misspell{Dir: dir, Filenames: filenames},    //
		IneffAssign{Dir: dir, Filenames: filenames}, //
		ErrCheck{Dir: dir, Filenames: filenames},    // disable errcheck for now, too slow and not finalized
	}
}

// Lint executes all checks on the given directory
// TODO: support linter options and optimise this function logic
func Lint(dir string) (model.ChecksResult, error) {
	log.Debugf("Lint recv params @dir=%s", dir)

	filenames, skipped, err := visitGoFiles(dir)
	if err != nil {
		return model.ChecksResult{}, fmt.Errorf("could not get filenames: %v", err)
	}
	if len(filenames) == 0 {
		return model.ChecksResult{}, fmt.Errorf("no .go files found")
	}

	err = RenameFiles(skipped)
	if err != nil {
		log.Errorf("Could not remove files, err=%v", err)
	}
	defer RevertFiles(skipped)

	var (
		linters   = getLinters(dir, filenames)
		n         = len(linters)
		chanScore = make(chan model.Score)
	)

	for _, linter := range linters {
		go execLinter(linter, chanScore)
	}

	var (
		total, totalWeight float64
		r                  = model.ChecksResult{
			Files: len(filenames),
		}
		issuesCnt int
		// issues             = make(map[string]bool)
	)

	// calc grade and score, then save into `model.CheckResult`
	for i := 0; i < n; i++ {
		score := <-chanScore
		r.Scores = append(r.Scores, score)
		total += score.Percentage * score.Weight
		totalWeight += score.Weight
		// for _, fs := range s.FileSummaries {
		// 	issues[fs.Filename] = true
		// }
		issuesCnt++
	}
	close(chanScore)

	total /= totalWeight
	sort.Sort(model.ByWeight(r.Scores))
	r.Average = total
	// r.Issues = len(issues)
	r.Issues = issuesCnt
	r.Grade = model.GradeFromPercentage(total * 100)

	return r, nil
}

// execLinter exec linter.Percentage and send model.Score by `chanScore`
func execLinter(linter ILinter, chanScore chan<- model.Score) {
	var errMsg string
	p, summaries, err := linter.Percentage()
	if err != nil {
		log.Errorf("Lint run linter=%s failed, err=%v", linter.Name(), err)
		errMsg = err.Error()
	}

	score := model.Score{
		Name:          linter.Name(),
		Description:   linter.Description(),
		FileSummaries: summaries,
		Weight:        linter.Weight(),
		Percentage:    p,
		Error:         errMsg,
	}
	chanScore <- score
}
