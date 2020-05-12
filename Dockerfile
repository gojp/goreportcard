FROM golang:1.12.5-alpine

RUN apk add --update --no-cache git make \
    && go get golang.org/x/tools/go/vcs

COPY . $GOPATH/src/github.com/gojp/goreportcard

WORKDIR $GOPATH/src/github.com/gojp/goreportcard

RUN ./scripts/make-install.sh

RUN go build \
    -o $GOPATH/bin/goreportcard-cli \
    $GOPATH/src/github.com/gojp/goreportcard/cmd/goreportcard-cli

EXPOSE 8000

CMD ["make", "start"]
