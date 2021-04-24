#!/bin/bash

need_gofmt=$(gofmt -s -l `find . -name '*.go' | grep -v vendor | grep -v _repos`)

if [[ -n ${need_gofmt} ]]; then
    echo "These files fail gofmt -s:"
    echo "${need_gofmt}"
    exit 1
fi


