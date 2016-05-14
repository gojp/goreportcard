all: lint test

install: 
	./scripts/make-install.sh

lint:
	gometalinter --exclude=vendor --exclude=repos --disable-all --enable=golint --enable=vet ./...
	find . -name '*.go' | xargs gofmt -w -s

test: 
	 go test -cover ./...

start:
	 go run main.go
