package check

// Staticcheck is the check for the staticcheck command
type Staticcheck struct {
	Dir       string
	Filenames []string
}

// Name returns the name of the display name of the command
func (g Staticcheck) Name() string {
	return "staticcheck"
}

// Weight returns the weight this check has in the overall average
func (g Staticcheck) Weight() float64 {
	return 0.15
}

// Percentage returns the percentage of .go files that pass
func (g Staticcheck) Percentage() (float64, []FileSummary, error) {
	return GoTool(g.Dir, g.Filenames, []string{"gometalinter", "--deadline=180s", "--disable-all", "--enable=staticcheck"})
}

// Description returns the description of Staticcheck
func (g Staticcheck) Description() string {
	return `<a href="https://staticcheck.io">Staticcheck</a> finds bugs, performance issues, and more`
}
