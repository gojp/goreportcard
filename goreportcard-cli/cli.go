package main

import (
	"fmt"
	"github.com/gojp/goreportcard/check"
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
)

var (
	dir     = kingpin.Arg("dir", "Root directory of your Go application").Required().String()
	verbose = kingpin.Flag("verbose", "Enable verbose output").Short('v').Bool()
)

func main() {
	kingpin.Parse()
	result, err := check.CheckDir(*dir)
	if err != nil {
		log.Fatalf("Fatal error checking %s: %s", *dir, err.Error())
	}

	fmt.Printf("Grade: %s\n", result.Grade)
	fmt.Printf("Average: %f\n", result.Average)
	fmt.Printf("Files: %d\n", result.Files)
	fmt.Printf("Issues: %d\n", result.Issues)

	if *verbose {
		for _, c := range result.Checks {
			fmt.Printf("\n%s:\n", c.Name)
			for _, f := range c.FileSummaries {
				for _, e := range f.Errors {
					fmt.Printf("%s:%d\n\t%s\n", f.Filename, e.LineNumber, e.ErrorString)
				}
			}
		}
	}
}
