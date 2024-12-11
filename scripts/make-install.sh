#!/bin/sh

go install ./vendor/github.com/alecthomas/gometalinter

go install ./vendor/github.com/fzipp/gocyclo/cmd/gocyclo
go install ./vendor/github.com/gordonklaus/ineffassign
go install ./vendor/github.com/client9/misspell/cmd/misspell
