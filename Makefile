install: 
	./scripts/make-install.sh

lint:
	golint ./...
	go vet ./...
	find . -name '*.go' | xargs gofmt -w -s

start:
	godep go run main.go
