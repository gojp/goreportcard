#!/bin/sh

go install ./vendor/github.com/alecthomas/gometalinter

go install ./vendor/golang.org/x/lint/golint
go install ./vendor/github.com/fzipp/gocyclo/cmd/gocyclo
go install ./vendor/github.com/gordonklaus/ineffassign
go install ./vendor/github.com/client9/misspell
