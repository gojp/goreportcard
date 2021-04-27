FROM golang:1.16.3-alpine

RUN apk add --update --no-cache git make \
    && go get golang.org/x/tools/go/vcs

RUN apk add build-base

COPY . $GOPATH/src/github.com/gojp/goreportcard

WORKDIR $GOPATH/src/github.com/gojp/goreportcard

RUN ./scripts/make-install.sh

EXPOSE 8000

CMD ["make", "start"]
