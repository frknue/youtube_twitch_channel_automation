package db

import (
	"fmt"
	"log"

	badger "github.com/dgraph-io/badger/v4"
)

func SaveClipID(clipID string) error {
	// Open the Badger database located in the /tmp/badger directory.
	db, err := badger.Open(badger.DefaultOptions("/tmp/badger"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Use Update function for transactional operation
	err = db.Update(func(txn *badger.Txn) error {
		_, err := txn.Get([]byte(clipID))
		if err == badger.ErrKeyNotFound {
			// Key not found, process the clip
			e := badger.NewEntry([]byte(clipID), []byte("processed"))
			err = txn.SetEntry(e)
			if err != nil {
				return err
			}
			log.Println("Clip ID processed and stored")
		} else if err != nil {
			return err
		} else {
			log.Println("Clip ID already processed")
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	return nil
}

// Check if the clip ID is already processed
func CheckClipID(clipID string) (bool, error) {
	// Open the Badger database located in the /tmp/badger directory.
	db, err := badger.Open(badger.DefaultOptions("/tmp/badger"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Use View function for read-only operation
	err = db.View(func(txn *badger.Txn) error {
		_, err := txn.Get([]byte(clipID))
		if err == badger.ErrKeyNotFound {
			log.Println("Clip ID not found")
			return err
		}
		log.Println("Clip ID found")
		return nil
	})

	if err != nil {
		log.Fatal(err)
		return false, err
	}

	return true, nil
}

func PrintClipIDs() error {
	// Open the Badger database located in the specified directory
	db, err := badger.Open(badger.DefaultOptions("/tmp/badger"))
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Iterate over all key-value pairs in the database
	err = db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			key := item.Key()
			err := item.Value(func(val []byte) error {
				fmt.Printf("Key: %s, Value: %s\n", key, val)
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Failed to iterate over database: %v", err)
	}

	return nil
}
