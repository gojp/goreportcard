all: lint test

install: 
	./scripts/make-install.sh

lint:
	golint ./...
	go vet ./...
	find . -name '*.go' | xargs gofmt -w -s

test: 
	 go test -cover ./...

start:
	 go run main.go
