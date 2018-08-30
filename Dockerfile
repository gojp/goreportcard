FROM golang:1.8-alpine


WORKDIR $GOPATH/src/github.com/gojp/goreportcard

RUN apk update && \
      apk upgrade && \
      apk add --no-cache git make && \
      go get golang.org/x/tools/go/vcs

EXPOSE 8000

CMD ["make", "start"]

COPY . .

RUN scripts/make-install.sh
