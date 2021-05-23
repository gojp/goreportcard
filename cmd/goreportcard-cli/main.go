package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gojp/goreportcard/check"
)

var (
	dir     = flag.String("d", ".", "Root directory of your Go application")
	verbose = flag.Bool("v", false, "Verbose output")
	th      = flag.Float64("t", 0, "Threshold of failure command")
	jsn     = flag.Bool("j", false, "JSON output. The binary will always exit with code 0")
	post    = flag.String("p", "", "Post local generate checks to cache")
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

	if *post != "" {
		respBytes, err := json.Marshal(result)
		if err != nil {
			log.Fatalf("Fatal could not marshal json: %v", err)
		}

		log.Println("Total grade: ", result.Average)

		request, err := http.NewRequest("POST", *post, bytes.NewBuffer(respBytes))
		if err != nil {
			log.Fatalf("Fatal error create request: %s", err.Error())
		}

		request.Header.Set("Content-Type", "application/json; charset=UTF-8")

		client := &http.Client{
			Timeout:   30 * time.Second,
			Transport: http.DefaultTransport,
		}
		response, err := client.Do(request)
		if err != nil {
			log.Fatalf("Fatal error do request: %s", err.Error())
		}
		defer response.Body.Close()

		log.Println("Send status", response.Status)

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatalf("Fatal error read request: %s", err.Error())
		}
		log.Println("response Body:", string(body))

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
