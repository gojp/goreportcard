package check

import (
	"fmt"
	"log"
	"sort"
)

// Check describes what methods various checks (gofmt, go lint, etc.)
// should implement
type Check interface {
	Name() string
	Description() string
	Weight() float64
	// Percentage returns the passing percentage of the check,
	// as well as a map of filename to output
	Percentage() (float64, []FileSummary, error)
}

// Score represents the result of a single check
type Score struct {
	Name          string        `json:"name"`
	Description   string        `json:"description"`
	FileSummaries []FileSummary `json:"file_summaries"`
	Weight        float64       `json:"weight"`
	Percentage    float64       `json:"percentage"`
	Error         string        `json:"error"`
}

// ChecksResult represents the combined result of multiple checks
type ChecksResult struct {
	Checks  []Score `json:"checks"`
	Average float64 `json:"average"`
	Grade   Grade   `json:"GradeFromPercentage"`
	Files   int     `json:"files"`
	Issues  int     `json:"issues"`
}

// Run executes all checks on the given directory
func Run(dir string) (ChecksResult, error) {
	filenames, skipped, err := GoFiles(dir)
	if err != nil {
		return ChecksResult{}, fmt.Errorf("could not get filenames: %v", err)
	}
	if len(filenames) == 0 {
		return ChecksResult{}, fmt.Errorf("no .go files found")
	}

	err = RenameFiles(skipped)
	if err != nil {
		log.Println("Could not remove files:", err)
	}
	defer RevertFiles(skipped)

	checks := []Check{
		GoFmt{Dir: dir, Filenames: filenames},
		GoVet{Dir: dir, Filenames: filenames},
		GoLint{Dir: dir, Filenames: filenames},
		GoCyclo{Dir: dir, Filenames: filenames},
		License{Dir: dir, Filenames: []string{}},
		Misspell{Dir: dir, Filenames: filenames},
		IneffAssign{Dir: dir, Filenames: filenames},
		// ErrCheck{Dir: dir, Filenames: filenames}, // disable errcheck for now, too slow and not finalized
	}

	ch := make(chan Score)
	for _, c := range checks {
		go func(c Check) {
			p, summaries, err := c.Percentage()
			errMsg := ""
			if err != nil {
				log.Printf("ERROR: (%s) %v", c.Name(), err)
				errMsg = err.Error()
			}
			s := Score{
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

	resp := ChecksResult{
		Files: len(filenames),
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
	resp.Grade = GradeFromPercentage(total * 100)

	return resp, nil
}

// ByWeight implements sorting for checks by weight descending
type ByWeight []Score

func (a ByWeight) Len() int           { return len(a) }
func (a ByWeight) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByWeight) Less(i, j int) bool { return a[i].Weight > a[j].Weight }
