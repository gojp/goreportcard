#!/bin/sh

go get github.com/golangci/golangci-lint/cmd/golangci-lint
gometalinter --install --update
