package linter

import "github.com/gojp/goreportcard/internal/model"

// IneffAssign is the check for the ineffassign command
type IneffAssign struct {
	Dir       string
	Filenames []string
}

// Name returns the name of the display name of the command
func (g IneffAssign) Name() string {
	return "ineffassign"
}

// Weight returns the weight this check has in the overall average
func (g IneffAssign) Weight() float64 {
	return 0.05
}

// Percentage returns the percentage of .go files that pass gofmt
// golangci-lint run --deadline=180s --disable-all --enable=ineffassign
func (g IneffAssign) Percentage() (float64, []model.FileSummary, error) {
	return GoTool(g.Dir, g.Filenames, []string{"golangci-lint", "run", "--deadline=180s", "--disable-all", "--enable=ineffassign"})
}

// Description returns the description of IneffAssign
func (g IneffAssign) Description() string {
	return `<a href="https://github.com/gordonklaus/ineffassign">IneffAssign</a> detects ineffectual assignments in Go code.`
}
