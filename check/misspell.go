package check

// Misspell is the check for the misspell command
type Misspell struct {
	Dir       string
	Filenames []string
}

// Name returns the name of the display name of the command
func (g Misspell) Name() string {
	return "misspell"
}

// Weight returns the weight this check has in the overall average
func (g Misspell) Weight() float64 {
	return 0.0
}

// Percentage returns the percentage of .go files that pass gofmt
func (g Misspell) Percentage() (float64, []FileSummary, error) {
	return GoTool(g.Dir, g.Filenames, []string{"golangci-lint", "run", "--deadline=180s", "--disable-all", "--enable=misspell"})
}

// Description returns the description of Misspell
func (g Misspell) Description() string {
	return `<a href="https://github.com/client9/misspell">Misspell</a> Finds commonly misspelled English words`
}
