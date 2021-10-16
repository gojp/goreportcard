package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/scanner"
	"go/token"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/gordonklaus/ineffassign/pkg/ineffassign"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/packages"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix(ineffassign.Analyzer.Name + ": ")

	flag.Usage = func() {
		paras := strings.Split(ineffassign.Analyzer.Doc, "\n\n")
		fmt.Fprintf(os.Stderr, "%s: %s\n\n", ineffassign.Analyzer.Name, paras[0])
		fmt.Fprintf(os.Stderr, "Usage: %s [-flag] [package]\n\n", ineffassign.Analyzer.Name)
		if len(paras) > 1 {
			fmt.Fprintln(os.Stderr, strings.Join(paras[1:], "\n\n"))
		}
		fmt.Fprintln(os.Stderr, "\nFlags:")
		flag.PrintDefaults()
	}
	flag.Parse()

	patterns := flag.Args()
	if len(patterns) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	pkgs, err := load(patterns...)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	diagnostics := make([][]analysis.Diagnostic, len(pkgs))
	var wg sync.WaitGroup
	for i, pkg := range pkgs {
		i, pkg := i, pkg
		wg.Add(1)
		go func() {
			ineffassign.Analyzer.Run(&analysis.Pass{
				Files: pkg.Syntax,
				Report: func(d analysis.Diagnostic) {
					diagnostics[i] = append(diagnostics[i], d)
				},
			})
			wg.Done()
		}()
	}
	wg.Wait()

	exitcode := 0
	for i, pkg := range pkgs {
		for _, d := range diagnostics[i] {
			pos := pkg.Fset.Position(d.Pos)
			fmt.Fprintf(os.Stderr, "%s: %s\n", pos, d.Message)
			exitcode = 3
		}
	}
	os.Exit(exitcode)
}

func load(patterns ...string) ([]*packages.Package, error) {
	conf := &packages.Config{
		// We would use NeedSyntax, but then Package.Fset is nil (https://github.com/golang/go/issues/48226).
		Mode:  packages.NeedFiles,
		Tests: true,
	}
	pkgs, err := packages.Load(conf, patterns...)
	if err != nil {
		return nil, err
	}

	for _, pkg := range pkgs {
		pkg.Fset = token.NewFileSet()
		pkg.Syntax = make([]*ast.File, len(pkg.GoFiles))
		for i, f := range pkg.GoFiles {
			var err error
			pkg.Syntax[i], err = parser.ParseFile(pkg.Fset, f, nil, parser.AllErrors|parser.ParseComments)
			if errlist, ok := err.(scanner.ErrorList); ok {
				for _, err := range errlist {
					pkg.Errors = append(pkg.Errors, packages.Error{
						Pos:  err.Pos.String(),
						Msg:  err.Msg,
						Kind: packages.ParseError,
					})
				}
			} else if err != nil {
				return nil, err
			}
		}
	}

	if n := packages.PrintErrors(pkgs); n > 1 {
		return nil, fmt.Errorf("%d errors during loading", n)
	} else if n == 1 {
		return nil, fmt.Errorf("error during loading")
	} else if len(pkgs) == 0 {
		return nil, fmt.Errorf("%s matched no packages", strings.Join(patterns, " "))
	}

	return pkgs, nil
}

func init() {
	flag.Bool("n", false, "no effect (deprecated)")
}
