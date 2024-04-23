[![Go Report Card](https://goreportcard.com/badge/gojp/goreportcard)](https://goreportcard.com/report/gojp/goreportcard) [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/gojp/goreportcard/blob/master/LICENSE)

# Go Report Card

A web application that generates a report on the quality of an open source Go project. It uses several measures, including `gofmt`, `go vet`, `go lint` and `gocyclo`. To get a report on your own project, try [goreportcard.com](https://goreportcard.com).

### Sponsors

Support us over on [Patreon](https://www.patreon.com/goreportcard)!

<a href="https://www.bairesdev.com/sponsoring-open-source-projects/"><img src="https://goreportcard.com/assets/bairesdev.png" width="50%" height="50%"></a>

<a href="https://www.digitalocean.com?utm_medium=opensource&utm_source=goreportcard"><img src="https://goreportcard.com/assets/digitalocean.svg" width="50%" height="50%"></a>

- [Cody Wood](https://www.linkedin.com/in/sprkyco/)
- Pascal Wenger
- Jonas Kwiedor
- [PhotoPrism](https://photoprism.app)
- Kia Farhang
- [Patrick DeVivo](https://twitter.com/patrickdevivo) ([MergeStat](https://github.com/mergestat/mergestat))
- [Alexis Geoffrey](https://github.com/alexisgeoffrey)

### Installation

```
git clone https://github.com/gojp/goreportcard.git
cd goreportcard
make install
```

Now run:

```
GRC_DATABASE_PATH=./db make start
```

and you should see

```
Running on :8000...
```

Navigate to `localhost:8000` and you should see the Go Report Card front page.

### Command Line Interface

There is also a CLI available for grading applications on your local machine.

Example usage:
```
git clone https://github.com/gojp/goreportcard.git
cd goreportcard
make install
go install ./cmd/goreportcard-cli
goreportcard-cli
```

```
Grade ........... A+ 98.1%
Files ................ 962
Issues ............... 158
gofmt ............... 100%
go_vet .............. 100%
ineffassign ......... 100%
gocyclo .............. 83%
license ............. 100%
misspell ............ 100%
```

Verbose output:

```
goreportcard-cli -v
```

```
Grade ........... A+ 98.1%
Files ................ 962
Issues ............... 158
gofmt ............... 100%
go_vet .............. 100%
ineffassign ......... 100%
gocyclo .............. 83%
        vendor/honnef.co/go/tools/staticcheck/structtag.go
                Line 8: warning: cyclomatic complexity 19 of function parseStructTag() is high (> 15) (gocyclo)

license ............. 100%
misspell ............ 100%
```

### Contributing

Go Report Card is an open source project run by volunteers, and contributions are welcome! Check out the [Issues](https://github.com/gojp/goreportcard/issues) page to see if your idea has already been mentioned. Feel free to raise an issue or submit a pull request.

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

The code is licensed under the permissive Apache v2.0 license. [Read this](https://tldrlegal.com/license/apache-license-2.0-(apache-2.0)) for a summary.
