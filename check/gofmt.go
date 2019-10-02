package check

// GoFmt is the check for the go fmt command
type GoFmt struct {
	Dir       string
	Filenames []string
}

// Name returns the name of the display name of the command
func (g GoFmt) Name() string {
	return "gofmt"
}

// Weight returns the weight this check has in the overall average
func (g GoFmt) Weight() float64 {
	return .30
}

// Percentage returns the percentage of .go files that pass gofmt
func (g GoFmt) Percentage() (float64, []FileSummary, error) {
	return GoTool(g.Dir, g.Filenames, []string{"golangci-lint", "run", "--deadline=180s", "--disable-all", "--enable=gofmt"})
	// return GoFmtNative(g.Dir, g.Filenames)
}

// Description returns the description of gofmt
func (g GoFmt) Description() string {
	return `Gofmt formats Go programs. We run <code>gofmt -s</code> on your code, where <code>-s</code> is for the <a href="https://golang.org/cmd/gofmt/#hdr-The_simplify_command">"simplify" command</a>`
}
