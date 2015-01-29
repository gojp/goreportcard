package check

type GoVet struct {
	Dir       string
	Filenames []string
}

func (g GoVet) Name() string {
	return "go_vet"
}

// Percentage returns the percentage of .go files that pass go vet
func (g GoVet) Percentage() (float64, []FileSummary, error) {
	return GoTool(g.Dir, g.Filenames, []string{"go", "tool", "vet"})
}

// Description returns the description of go lint
func (g GoVet) Description() string {
	return `<code>go vet</code> examines Go source code and reports suspicious constructs, such as Printf calls whose arguments do not align with the format string.`
}
