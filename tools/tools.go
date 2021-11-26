//go:build tools
// +build tools

package tools

import (
	_ "github.com/alecthomas/gocyclo"
	_ "github.com/alecthomas/gometalinter"
	_ "github.com/client9/misspell/cmd/misspell"
	_ "github.com/gordonklaus/ineffassign"
	_ "golang.org/x/lint/golint"
	_ "honnef.co/go/tools/cmd/staticcheck"
)
