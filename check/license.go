package check

import (
	"io/ioutil"
	"path/filepath"
	"strings"
)

// License is the check for the existence of a license file
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
	return .05
}

// thank you https://github.com/ryanuber/go-license and client9
var licenses = []string{
	"license",
	"copying",
	"copyright",
	"licence",
	"unlicense",
	"copyleft",
}

// Percentage returns 0 if no LICENSE, 1 if LICENSE
func (g License) Percentage() (float64, []FileSummary, error) {
	files, err := ioutil.ReadDir(g.Dir)
	if err != nil {
		return 0.0, []FileSummary{}, err
	}

	for _, file := range files {
		name := strings.ToLower(file.Name())

		if filepath.Ext(name) == "go" {
			continue
		}

		for i := range licenses {
			if strings.HasPrefix(name, licenses[i]) {
				return 1.0, []FileSummary{}, nil
			}
		}
	}

	return 0.0, []FileSummary{{"", "http://choosealicense.com/", []Error{}}}, nil
}

// Description returns the description of License
func (g License) Description() string {
	return "Checks whether your project has a LICENSE file."
}
