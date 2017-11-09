package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"time"
)

var real = flag.Bool("real", false, "run the deletions")

func main() {
	flag.Parse()
	files, err := ioutil.ReadDir("_repos/src")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if f.IsDir() {
			dirs, err := ioutil.ReadDir("_repos/src/" + f.Name())
			if err != nil {
				log.Fatal(err)
			}
			for _, d := range dirs {
				if time.Now().Sub(d.ModTime()) > 30*24*time.Hour {
					path := "_repos/src/" + f.Name() + "/" + d.Name()
					if *real {
						log.Printf("Deleting %s...", path)
						os.RemoveAll(path)
					} else {
						log.Printf("Would delete %s", path)
					}
				}
			}
		}
	}
}
