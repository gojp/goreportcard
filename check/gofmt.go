package check

type GoFmt struct {
	Dir       string
	Filenames []string
}

func (g GoFmt) Name() string {
	return "gofmt"
}

// Percentage returns the percentage of .go files that pass gofmt
func (g GoFmt) Percentage() (float64, []FileSummary, error) {
	return GoTool(g.Dir, g.Filenames, []string{"gofmt", "-s", "-l"})
}

// Description returns the description of gofmt
func (g GoFmt) Description() string {
	return `Gofmt formats Go programs. We run <code>gofmt -s</code> on your code, where <code>-s</code> is for the <a href="https://golang.org/cmd/gofmt/#hdr-The_simplify_command">"simplify" command</a>`
}
