package check

import (
	"bytes"
	"os/exec"
)

// License is the check for the existance of a license file
type License struct {
	Dir       string
	Filenames []string
}

// Name returns the name of the display name of the command
func (g License) Name() string {
	return "license"
}

// Weight returns the weight this check has in the overall average
func (g License) Weight() float64 {
	return .10
}

// Percentage returns 0 if no LICENSE, 1 if LICENSE
func (g License) Percentage() (float64, []FileSummary, error) {
	cmd := exec.Command("find", g.Dir, "-maxdepth", "1", "-type", "f", "-name", "LICENSE*")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return 0.0, []FileSummary{}, err
	}
	if out.String() == "" {
		return 0.0, []FileSummary{{"", "http://choosealicense.com/", []Error{}}}, nil
	}

	return 1.0, []FileSummary{}, nil
}

// Description returns the description of License
func (g License) Description() string {
	return "Checks whether your project has a LICENSE file."
}
