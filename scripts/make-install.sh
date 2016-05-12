#!/bin/bash

go get github.com/tools/godep

go get github.com/alecthomas/gometalinter
gometalinter --install --update
go get github.com/client9/misspell/cmd/misspell
