#!/bin/bash

go get github.com/tools/godep
go get github.com/boltdb/bolt

go get github.com/alecthomas/gometalinter
gometalinter --install --update
