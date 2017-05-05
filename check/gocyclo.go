package check

// GoCyclo is the check for the go cyclo command
type GoCyclo struct {
	Dir       string
	Filenames []string
}

// Name returns the name of the display name of the command
func (g GoCyclo) Name() string {
	return "gocyclo"
}

// Weight returns the weight this check has in the overall average
func (g GoCyclo) Weight() float64 {
	return .10
}

// Percentage returns the percentage of .go files that pass gofmt
func (g GoCyclo) Percentage() (float64, []FileSummary, error) {
	return GoTool(g.Dir, g.Filenames, []string{"gometalinter", "--deadline=180s", "--disable-all", "--enable=gocyclo", "--cyclo-over=15"})
}

// Description returns the description of GoCyclo
func (g GoCyclo) Description() string {
	return `<a href="https://github.com/fzipp/gocyclo">Gocyclo</a> calculates cyclomatic complexities of functions in Go source code.

The cyclomatic complexity of a function is calculated according to the following rules:

1 is the base complexity of a function
+1 for each 'if', 'for', 'case', '&&' or '||'

Go Report Card warns on functions with cyclomatic complexity > 15.`
}
