package main

import (
	"io/ioutil"
	"log"
	"os"
	"time"
)

func main() {
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
					log.Printf("Deleting %s...", path)
					os.RemoveAll(path)
				}
			}
		}
	}
}
