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

// Weight returns the weight this check has in the overall average
func (g GoLint) Weight() float64 {
	return .10
}

// Percentage returns the percentage of .go files that pass golint
func (g GoLint) Percentage() (float64, []FileSummary, error) {
	return GoTool(g.Dir, g.Filenames, []string{"gometalinter", "--deadline=180s", "--disable-all", "--enable=golint", "--min-confidence=0.85", "--vendor"})
}

// Description returns the description of go lint
func (g GoLint) Description() string {
	return `Golint is a linter for Go source code.`
}
