package main

import (
	"fmt"
	"os"
)

func printUsage() {
	fmt.Println("\ngoreportcard command line tool")
	fmt.Println("\nUsage:\n\n  goreportcard [package names...]")
}

func main() {
	if len(os.Args) == 1 {
		printUsage()
	}

	names := := os.Args[1:]

}
