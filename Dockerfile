FROM golang:1.10-alpine

COPY . $GOPATH/src/github.com/gojp/goreportcard

WORKDIR $GOPATH/src/github.com/gojp/goreportcard

RUN apk update && apk upgrade && apk add --no-cache git make \
        && go get golang.org/x/tools/go/vcs \
        && ./scripts/make-install.sh

EXPOSE 8000

CMD ["make", "start"]
