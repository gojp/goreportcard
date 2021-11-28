package main

import (
	"flag"
	"log"
	"strings"

	"github.com/dgraph-io/badger/v2"
)

var (
	deleteRepoName       = flag.String("deleterepo", "", "repo to delete from badger cache")
	removeDuplicatesFlag = flag.Bool("removeduplicates", false, "remove non-lowercase duplicates from badger cache")
	dryRun               = flag.Bool("dryrun", false, "dry run mode")
)

func main() {
	flag.Parse()

	db, err := badger.Open(badger.DefaultOptions("/usr/local/badger").WithTruncate(true))
	if err != nil {
		log.Fatal("ERROR: could not open badger db: ", err)
	}

	defer db.Close()

	if *deleteRepoName != "" {
		err := deleteRepo(db, *deleteRepoName, *dryRun)
		if err != nil {
			log.Fatal(err)
		}
	}

	if *removeDuplicatesFlag {
		err := removeDuplicates(db, *dryRun)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// delete a repo from badger cache
func deleteRepo(db *badger.DB, repo string, dryRun bool) error {
	if !dryRun {
		return db.Update(func(txn *badger.Txn) error {
			return txn.Delete([]byte("repos-" + repo))
		})
	}

	log.Println("would delete", "repos-"+repo)

	return nil
}

// remove duplicates from badger cache
func removeDuplicates(db *badger.DB, dryRun bool) error {
	var toRemove []string

	err := db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := string(item.Key())

			if strings.HasPrefix(k, "repos-") {
				if strings.ToLower(k) != k {
					toRemove = append(toRemove, k)
				}
			}
		}

		return nil
	})

	for _, k := range toRemove {
		err := db.View(func(txn *badger.Txn) error {
			_, err := txn.Get([]byte(strings.ToLower(k)))

			return err
		})

		if err == badger.ErrKeyNotFound {
			continue
		}

		if err != nil {
			return err
		}

		if !dryRun {
			err = db.Update(func(txn *badger.Txn) error {
				log.Println("Deleting", k)
				return txn.Delete([]byte(k))
			})

			if err != nil {
				return err
			}
		} else {
			log.Println("would delete", k)
		}
	}

	return err
}
