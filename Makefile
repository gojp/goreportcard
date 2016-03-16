all: lint test

install: 
	./scripts/make-install.sh

lint:
	gometalinter --debug --exclude=vendor --disable-all --enable=lint,vet ./...
	find . -name '*.go' | xargs gofmt -w -s

test: 
	 go test -cover ./...

start:
	 go run main.go
