package main

import (
	"flag"
	"fmt"
	"encoding/json"
	"log"
	"os"

	"github.com/gojp/goreportcard/check"
)

var (
	dir     = flag.String("d", ".", "Root directory of your Go application")
	verbose = flag.Bool("v", false, "Verbose output")
	th      = flag.Float64("t", 0, "Threshold of failure command")
	jsn	= flag.Bool("j", false, "JSON output. The binary will always exit with code 0")
)

func main() {
	flag.Parse()

	result, err := check.Run(*dir)
	if err != nil {
		log.Fatalf("Fatal error checking %s: %s", *dir, err.Error())
	}

	if *jsn {
		marshalledResults, _ := json.Marshal(result)
		fmt.Println(string(marshalledResults))
		os.Exit(0)
	}

	fmt.Printf("Grade: %s (%.1f%%)\n", result.Grade, result.Average*100)
	fmt.Printf("Files: %d\n", result.Files)
	fmt.Printf("Issues: %d\n", result.Issues)

	for _, c := range result.Checks {
		fmt.Printf("%s: %d%%\n", c.Name, int64(c.Percentage*100))
		if *verbose && len(c.FileSummaries) > 0 {
			for _, f := range c.FileSummaries {
				fmt.Printf("\t%s\n", f.Filename)
				for _, e := range f.Errors {
					fmt.Printf("\t\tLine %d: %s\n", e.LineNumber, e.ErrorString)
				}
			}
		}
	}

	if result.Average*100 < *th {
		os.Exit(1)
	}
}
