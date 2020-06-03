all: lint build test

build:
	go build ./cmd/goreportcard-cli

install: 
	./scripts/make-install.sh

lint:
	golangci-lint run --skip-dirs=repos --disable-all \
		--enable=golint --enable=vet --enable=gofmt --enable=misspell ./...

test: 
	go test -cover ./internal

start:
	go run ./cmd/goreportcard-cli/ start-web

