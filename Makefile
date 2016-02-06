all: lint test

install: 
	./scripts/make-install.sh

lint:
	golint ./...
	go vet ./...
	find . -name '*.go' | xargs gofmt -w -s

test: 
	godep go test -cover ./...

start:
	godep go run main.go
