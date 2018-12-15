package main

import (
	"fmt"
	"github.com/gojp/goreportcard/check"
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
)

var (
	dir     = kingpin.Arg("dir", "Root directory of your Go application").Required().String()
	verbose = kingpin.Flag("verbose", "Verbose output").Short('v').Bool()
)

func main() {
	kingpin.Parse()
	result, err := check.CheckDir(*dir)
	if err != nil {
		log.Fatalf("Fatal error checking %s: %s", *dir, err.Error())
	}

	fmt.Printf("Grade: %s (%.1f%%)\n", result.Grade, result.Average*100)
	fmt.Printf("Files: %d\n", result.Files)
	fmt.Printf("Issues: %d\n", result.Issues)

	for _, c := range result.Checks {
		fmt.Printf("%s: %d%%\n", c.Name, int64(c.Percentage*100))
		if *verbose && len(c.FileSummaries) > 0 {
			for _, f := range c.FileSummaries {
				for _, e := range f.Errors {
					fmt.Printf("%s\t%s:%d\n\t%s\n", c.Name, f.Filename, e.LineNumber, e.ErrorString)
				}
			}
			fmt.Println()
		}
	}
}
