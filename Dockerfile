# syntax=docker/dockerfile:1
FROM golang:1.16.3-alpine

RUN apk add --update --no-cache ca-certificates git make \
    && go get golang.org/x/tools/go/vcs

RUN apk add build-base
RUN mkdir /app && mkdir /app/data/ && touch /app/data/.git-credentials

COPY . $GOPATH/src/github.com/gojp/goreportcard

WORKDIR $GOPATH/src/github.com/gojp/goreportcard

RUN ./scripts/make-install.sh

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/goreportcard .
RUN rm -rf $GOPATH/src/github.com/gojp/goreportcard

WORKDIR /app

EXPOSE 8000

ENV GIT_TERMINAL_PROMPT 0

CMD ["./goreportcard", "-b", "./data/badger/"]
