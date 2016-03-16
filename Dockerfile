FROM golang:1.6

RUN go get golang.org/x/tools/go/vcs

COPY . $GOPATH/src/github.com/gojp/goreportcard

WORKDIR $GOPATH/src/github.com/gojp/goreportcard

EXPOSE 8080

CMD ["make", "start"]
