package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

func main() {
	repoSrc := filepath.Join("_repos", "src")
	files, err := ioutil.ReadDir(repoSrc)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if f.IsDir() {
			dirs, err := ioutil.ReadDir(filepath.Join(repoSrc, f.Name()))
			if err != nil {
				log.Fatal(err)
			}
			for _, d := range dirs {
				if time.Now().Sub(d.ModTime()) > 30*24*time.Hour {
					path := filepath.Join(repoSrc, f.Name(), d.Name())
					log.Printf("Deleting %s...", path)
					os.RemoveAll(path)
				}
			}
		}
	}
}
