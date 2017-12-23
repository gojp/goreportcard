package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
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
				path := "_repos/src/" + f.Name() + "/" + d.Name()
				if time.Now().Sub(d.ModTime()) > 30*24*time.Hour {
					if *real {
						log.Printf("Deleting %s (repo is old)...", path)
						os.RemoveAll(path)
						continue
					} else {
						log.Printf("Would delete %s (repo is old)", path)
					}
				}

				size, err := DirSize(path)
				if err != nil {
					log.Fatal(err)
				}
				if size < 15*1000*1000 {
					if *real {
						log.Printf("Deleting %s (repo size < 15M)...", path)
						os.RemoveAll(path)
					} else {
						log.Printf("Would delete %s (repo size < 15M)", path)
					}
				}
			}
		}
	}
}

func DirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}
