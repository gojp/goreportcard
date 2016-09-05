package main

import (
	"fmt"
	"log"

	"github.com/gojp/goreportcard/check"
	"github.com/gojp/goreportcard/handlers"
)

type score struct {
	Name          string              `json:"name"`
	Description   string              `json:"description"`
	FileSummaries []check.FileSummary `json:"file_summaries"`
	Weight        float64             `json:"weight"`
	Percentage    float64             `json:"percentage"`
	Error         string              `json:"error"`
}

var allScores []score

func main() {
	dir := "."
	filenames, skipped, err := check.GoFiles(dir)
	if err != nil {
		log.Fatalf("could not get filenames: %v", err)
	}
	if len(filenames) == 0 {
		log.Fatalf("no .go files found")
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

	var (
		total       float64
		totalWeight float64
	)
	for i := 0; i < len(checks); i++ {
		s := <-ch
		allScores = append(allScores, s)
		total += s.Percentage * s.Weight
		totalWeight += s.Weight
	}
	total /= totalWeight

	grade := handlers.PercentToGrade(total * 100)

	for _, score := range allScores {
		fmt.Printf("%s: %.2f%%\n", score.Name, score.Percentage*100)
	}
	fmt.Println("Grade:", grade)
}
