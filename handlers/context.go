package handlers

// Context defines the interface certain requests expect as an additional
// argument.
type Context interface {
	Suggest(string) ([]string, error)
}
