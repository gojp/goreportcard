[![Go Report Card](https://goreportcard.com/badge/gojp/goreportcard)](https://goreportcard.com/report/gojp/goreportcard) [![Build Status](https://travis-ci.org/gojp/goreportcard.svg?branch=master)](https://travis-ci.org/gojp/goreportcard) [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/gojp/goreportcard/blob/master/LICENSE)

# Go Report Card

A web application that generates a report on the quality of an open source go project. It uses several measures, including `gofmt`, `go vet`, `go lint` and `gocyclo`. To get a report on your own project, try using the hosted version of this code running at [goreportcard.com](https://goreportcard.com).

### Sponsors

Support us over on [Patreon](https://www.patreon.com/goreportcard)!

<a href="https://cooperpress.com"><img src="https://goreportcard.com/assets/cooperpress.png" width="50%" height="50%"></a>

<a href="https://www.digitalocean.com?utm_medium=opensource&utm_source=goreportcard"><img src="https://goreportcard.com/assets/digitalocean.svg" width="50%" height="50%"></a>

- [Cody Wood](https://www.linkedin.com/in/sprkyco/)

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
Running on 127.0.0.1:8000...
```

Navigate to that URL in your browser and check that you can see the front page.

### Command Line Interface

There is also a CLI available for grading applications on your local machine.

Example usage:
```
go get github.com/gojp/goreportcard/cmd/goreportcard-cli
cd $GOPATH/src/github.com/gojp/goreportcard
goreportcard-cli
```

```
Grade: A+ (99.9%)
Files: 362
Issues: 2
gofmt: 100%
go_vet: 99%
gocyclo: 99%
golint: 100%
ineffassign: 100%
license: 100%
misspell: 100%
```

Verbose output is also available:
```
goreportcard-cli -v
```

```
Grade: A+ (99.9%)
Files: 332
Issues: 2
gofmt: 100%
go_vet: 99%
go_vet  vendor/github.com/prometheus/client_golang/prometheus/desc.go:25
        error: cannot find package "github.com/prometheus/client_model/go" in any of: (vet)

gocyclo: 99%
gocyclo download/download.go:22
        warning: cyclomatic complexity 17 of function download() is high (> 15) (gocyclo)

golint: 100%
ineffassign: 100%
license: 100%
misspell: 100%
```

### Contributing

Go Report Card is an open source project run by volunteers, and contributions are welcome! Check out the [Issues](https://github.com/gojp/goreportcard/issues) page to see if your idea for a contribution has already been mentioned, and feel free to raise an issue or submit a pull request.

### Academic Citation

If you use Go Report Card for academic purposes, please use the following citation:

```
@Misc{schaaf-smith-goreportcard,
    author = {Schaaf, Herman and Smith, Shawn},
    title  = {Go Report Card: A report card for your Go application},
    year   = {2015--},
    url    = {https://www.goreportcard.com/},
    note   = {[Online; accessed <today>]}
}
```

### License

The code is licensed under the permissive Apache v2.0 licence. This means you can do what you like with the software, as long as you include the required notices. [Read this](https://tldrlegal.com/license/apache-license-2.0-(apache-2.0)) for a summary.

### Notes

We don't support Go Report Card on Windows.
