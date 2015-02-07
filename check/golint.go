package check

// GoLint is the check for the go lint command
type GoLint struct {
	Dir       string
	Filenames []string
}

// Name returns the name of the display name of the command
func (g GoLint) Name() string {
	return "golint"
}

// Percentage returns the percentage of .go files that pass golint
func (g GoLint) Percentage() (float64, []FileSummary, error) {
	return GoTool(g.Dir, g.Filenames, []string{"golint"})
}

// Description returns the description of go lint
func (g GoLint) Description() string {
	return `Golint is a linter for Go source code.`
}
