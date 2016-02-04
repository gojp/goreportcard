[![Go Report Card](http://goreportcard.com/badge/gojp/goreportcard)](http://goreportcard.com/report/gojp/goreportcard) [![Build Status](https://travis-ci.org/gojp/goreportcard.svg?branch=master)](https://travis-ci.org/gojp/goreportcard)

# Go Report Card

A web application that generates a report on the quality of an open source go project. It uses several measures, including `gofmt`, `go vet`, `go lint` and `gocyclo`. To get a report on your own project, try using the hosted version of this code running at [goreportcard.com](http://goreportcard.com). Currently this is limited to projects hosted on Github, but there is an [open issue](https://github.com/gojp/goreportcard/issues/51) to support all publicly hosted repos.

### Installation

Assuming you already have a recent version of Go installed, pull down the code with `go get`:

```
go get github.com/gojp/goreportcard
```

Go into the source directory and pull down the project dependencies:

```
cd $GOPATH/src/github.com/gojp/goreportcard
make install
```

Now run

```
make start
```

and you should see

```
Running on 127.0.0.1:8080...
```

Navigate to that URL in your browser and check that you can see the front page.


### Contributing

Go Report Card is an open source project run by volunteers, and contributions are welcome! Check out the [Issues](https://github.com/gojp/goreportcard/issues) page to see if your idea for a contribution has already been mentioned, and feel free to raise an issue or submit a pull request.

### License

The code is licensed under the permissive Apache v2.0 licence. This means you can do what you like with the software, as long as you include the required notices. [Read this](https://tldrlegal.com/license/apache-license-2.0-(apache-2.0)) for a summary.
