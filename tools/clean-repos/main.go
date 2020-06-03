package main

import (
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/yeqown/log"
)

var real = flag.Bool("real", false, "run the deletions")

func main() {
	flag.Parse()
	files, err := ioutil.ReadDir("_repos/src")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if f.Name() != "github.com" {
			continue
		}
		if f.IsDir() {
			dirs, err := ioutil.ReadDir("_repos/src/" + f.Name())
			if err != nil {
				log.Fatal(err)
			}
			for _, d := range dirs {
				repos, err := ioutil.ReadDir("_repos/src/" + f.Name() + "/" + d.Name())
				if err != nil {
					log.Fatal(err)
				}
				for _, repo := range repos {
					path := "_repos/src/" + f.Name() + "/" + d.Name() + "/" + repo.Name()
					if time.Since(d.ModTime()) > 30*24*time.Hour {
						if *real {
							log.Infof("Deleting %s (repo is old)...", path)
							os.RemoveAll(path)
							continue
						} else {
							log.Infof("Would delete %s (repo is old)", path)
						}
					}

					size, err := DirSize(path)
					if err != nil {
						log.Fatal(err)
					}
					if size < 50*1000*1000 {
						if *real {
							log.Infof("Deleting %s (dir size < 50M)...", path)
							os.RemoveAll(path)
						} else {
							log.Infof("Would delete %s (dir size < 50M)", path)
						}
					}
				}
			}
		}
	}
}

// DirSize returns the size of a directory
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
