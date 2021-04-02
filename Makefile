all: lint build test

build:
	go build ./...

install: 
	./scripts/make-install.sh

lint:
	golangci-lint run --skip-dirs=repos --disable-all --enable=golint --enable=vet --enable=gofmt ./...
	find . -name '*.go' | xargs gofmt -w -s

test: 
	 go test -cover ./check ./handlers

start:
	 go run main.go

misspell:
	find . -name '*.go' -not -path './vendor/*' -not -path './_repos/*' | xargs misspell -error
