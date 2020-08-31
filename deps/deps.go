package deps

// import all library we wish to version check with modules but are not directly used in the project
// as this package is neither imported in the project, it does not impact the final executable file.
import (
	// Ensure go-farm versioning
	_ "github.com/dgryski/go-farm"
	// Ensure errors versioning
	_ "github.com/pkg/errors"
	// Ensure procfs versioning
	_ "github.com/prometheus/procfs"
)
