FROM golang:1.5

RUN go get github.com/tools/godep && \
	go get golang.org/x/tools/go/vcs

COPY . $GOPATH/src/github.com/gojp/goreportcard

WORKDIR $GOPATH/src/github.com/gojp/goreportcard

EXPOSE 8080

CMD ["make", "start"]
