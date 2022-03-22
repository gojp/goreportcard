all: lint build test

build:
	go build ./...

install:
	./scripts/make-install.sh

lint: fmt vet staticcheck misspell

fmt:
	./scripts/gofmt.sh

vet:
	go vet ./check ./cmd/... ./download ./handlers ./tools/...
	go vet ./main.go

staticcheck:
	@[ -x "$(shell which staticcheck)" ] || go install honnef.co/go/tools/cmd/staticcheck@master
	staticcheck ./...

test:
	 go test -cover ./...

start:
	 go run main.go

misspell:
	@[ -x "$(shell which misspell)" ] || go install ./vendor/github.com/client9/misspell/cmd/misspell
	find . -name '*.go' -not -path './vendor/*' -not -path './_repos/*' -not -path './download/test_downloads/*' -not -path './check/testdata/*' | xargs misspell -error
