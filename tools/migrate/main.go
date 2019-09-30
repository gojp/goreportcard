// migrate cache from boltdb to badger
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/boltdb/bolt"
	"github.com/dgraph-io/badger"
)

func migrate() error {
	db, err := bolt.Open("/Users/shawn/goreportcard.db", 0755, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return fmt.Errorf("could not open bolt db: %v", err)
	}
	defer db.Close()

	badgerDB, err := badger.Open(badger.DefaultOptions("/tmp/badger"))
	if err != nil {
		return err
	}

	// migrate meta info
	err = db.View(func(tx *bolt.Tx) error {
		rb := tx.Bucket([]byte("meta"))
		if rb == nil {
			return fmt.Errorf("scores bucket not found")
		}

		err = rb.ForEach(func(k, v []byte) error {
			log.Printf("Migrating %q...", string(k))
			err = badgerDB.Update(func(txn *badger.Txn) error {
				return txn.Set(k, v)
			})

			return err
		})

		if err != nil {
			return err
		}

		return nil
	})

	// migrate repo data
	err = db.View(func(tx *bolt.Tx) error {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered", r)
			}
		}()

		rb := tx.Bucket([]byte("repos"))
		if rb == nil {
			return fmt.Errorf("repos bucket not found")
		}

		err = rb.ForEach(func(k, v []byte) error {
			log.Printf("Migrating %q...", string(k))
			err = badgerDB.Update(func(txn *badger.Txn) error {
				return txn.Set([]byte("repos-"+string(k)), v)
			})

			return err
		})

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func main() {
	err := migrate()
	if err != nil {
		log.Fatal(err)
	}
}
