package check

// GoVet is the check for the go vet command
type GoVet struct {
	Dir       string
	Filenames []string
}

// Name returns the name of the display name of the command
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
