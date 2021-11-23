package main

import (
	"log"

	"github.com/dgraph-io/badger/v2"
)

func main() {
	// delete a repo from badger cache
	db, err := badger.Open(badger.DefaultOptions("/usr/local/badger").WithTruncate(true))
	if err != nil {
		log.Fatal("ERROR: could not open badger db: ", err)
	}

	defer db.Close()

	repo := "[repo]"

	err = db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte("repos-" + repo))
	})

	if err != nil {
		log.Fatal("Badger writing error:", err)
	}
}
