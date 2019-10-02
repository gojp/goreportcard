package check

// ErrCheck is the check for the errcheck command
type ErrCheck struct {
	Dir       string
	Filenames []string
}

// Name returns the name of the display name of the command
func (c ErrCheck) Name() string {
	return "errcheck"
}

// Weight returns the weight this check has in the overall average
func (c ErrCheck) Weight() float64 {
	return .15
}

// Percentage returns the percentage of .go files that pass gofmt
func (c ErrCheck) Percentage() (float64, []FileSummary, error) {
	return GoTool(c.Dir, c.Filenames, []string{"golangci-lint", "run", "--deadline=180s", "--disable-all", "--linter='errch:errcheck {path}:PATH:LINE:MESSAGE'", "--enable=errch"})
}

// Description returns the description of gofmt
func (c ErrCheck) Description() string {
	return `<a href="https://github.com/kisielk/errcheck">errcheck</a> finds unchecked errors in go programs`
}
