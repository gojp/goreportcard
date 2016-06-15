package check

// errcheck is the check for the errcheck command
type Errcheck struct {
	Dir       string
	Filenames []string
}

// Name returns the name of the display name of the command
func (e Errcheck) Name() string {
	return "errcheck"
}

// Weight returns the weight this check has in the overall average
func (e Errcheck) Weight() float64 {
	return .10
}

// Percentage returns the percentage of .go files that pass errcheck
func (e Errcheck) Percentage() (float64, []FileSummary, error) {
	return GoTool(e.Dir, e.Filenames, []string{"gometalinter", "--deadline=180s", "--disable-all", "--enable=errcheck", "--min-confidence=0.85"})
}

// Description returns the description of errcheck
func (e Errcheck) Description() string {
	return `errcheck checks that you checked errors.`
}
