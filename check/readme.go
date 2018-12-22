package check

import (
	"io/ioutil"
	"path/filepath"
	"strings"
)

// Readme is the check for the existence of a readme file
type Readme struct {
	Dir       string
	Filenames []string
}

// Name returns the name of the display name of the command
func (g Readme) Name() string {
	return "readme"
}

// Weight returns the weight this check has in the overall average
func (g Readme) Weight() float64 {
	return .05
}

// Percentage returns 0 if no README, 1 if README
func (g Readme) Percentage() (float64, []FileSummary, error) {
	files, err := ioutil.ReadDir(g.Dir)
	if err != nil {
		return 0.0, []FileSummary{}, err
	}

	for _, file := range files {
		name := strings.ToLower(file.Name())

		if filepath.Ext(name) == "go" {
			continue
		}

		if strings.HasPrefix(name, "readme") {
			return 1.0, []FileSummary{}, nil
		}
	}

	return 0.0, []FileSummary{{"", "Add a readme file", []Error{}}}, nil
}

// Description returns the description of Readme
func (g Readme) Description() string {
	return "Checks whether your project has a README file."
}
