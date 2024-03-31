package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gojp/goreportcard/check"
)

var (
	dir     = flag.String("d", ".", "Root directory of your Go application")
	verbose = flag.Bool("v", false, "Verbose output")
	th      = flag.Float64("t", 0, "Threshold of failure command")
	jsn     = flag.Bool("j", false, "JSON output. The binary will always exit with code 0")
)

// dotPrintf fills in the blank space between two strings with dots. The total
// length displayed is indicated by fullLen. The left string is specified by
// lfStr. The right string is formatted. At least two dots are shown, even if
// this means the total length exceeds fullLen.
func dotPrintf(fullLen int, lfStr, rtFmtStr string, args ...interface{}) {
	rtStr := fmt.Sprintf(rtFmtStr, args...)
	dotLen := fullLen - len(lfStr) - len(rtStr)
	if dotLen < 2 {
		dotLen = 2
	}
	fmt.Printf("%s %s %s\n", lfStr, strings.Repeat(".", dotLen), rtStr)
}

func main() {
	flag.Parse()

	result, err := check.Run(*dir, true)
	if err != nil {
		log.Fatalf("Fatal error checking %s: %s", *dir, err.Error())
	}

	if *jsn {
		marshalledResults, _ := json.Marshal(result)
		fmt.Println(string(marshalledResults))
		os.Exit(0)
	}

	dotPrintf(24, "Grade", "%s %.1f%%", result.Grade, result.Average*100)
	dotPrintf(24, "Files", "%d", result.Files)
	dotPrintf(24, "Issues", "%d", result.Issues)

	for _, c := range result.Checks {
		dotPrintf(24, c.Name, "%d%%", int64(c.Percentage*100))
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
