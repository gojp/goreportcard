package check

// Check describes what methods various checks (gofmt, go lint, etc.)
// should implement
type Check interface {
	Name() string
	Description() string
	Weight() float64
	// Percentage returns the passing percentage of the check,
	// as well as a map of filename to output
	Percentage() (float64, []FileSummary, error)
}
