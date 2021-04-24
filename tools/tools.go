// +build tools

package tools

import (
	_ "github.com/client9/misspell/cmd/misspell"
	_ "golang.org/x/lint/golint"
	_ "honnef.co/go/tools/cmd/staticcheck"
)
