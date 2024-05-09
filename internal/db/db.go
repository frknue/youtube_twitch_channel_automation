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
func CheckClipID(clipID string) bool {
	// Open the Badger database located in the /tmp/badger directory.
	db, err := badger.Open(badger.DefaultOptions("/tmp/badger"))
	if err != nil {
		log.Printf("Error opening database: %v", err)
		return false // return false to indicate the failure of the database operation
	}
	defer db.Close()

	// Use View function for read-only operation
	found := false
	err = db.View(func(txn *badger.Txn) error {
		_, err := txn.Get([]byte(clipID))
		if err == badger.ErrKeyNotFound {
			log.Println("Clip ID not found")
			return nil // return nil to continue normally
		}
		if err != nil {
			return err // return an actual error that might have occurred
		}
		log.Println("Clip ID found")
		found = true // set found to true if no errors occurred and clip ID is found
		return nil
	})

	if err != nil {
		log.Printf("Error checking clip ID: %v", err)
		return false
	}

	return found
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

func CleanUpDB() error {
	// Open the Badger database located in the specified directory
	db, err := badger.Open(badger.DefaultOptions("/tmp/badger"))
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Use Update function for transactional operation
	err = db.Update(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			key := item.Key()
			err := item.Value(func(val []byte) error {
				if string(val) == "processed" {
					err := txn.Delete(key)
					if err != nil {
						return err
					}
					log.Printf("Deleted key: %s\n", key)
				}
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Failed to clean up database: %v", err)
	}

	return nil
}

// Create a lock to prevent multiple instances of the application from running concurrently
func CreateLock() error {
	// Open the Badger database located in the specified directory
	db, err := badger.Open(badger.DefaultOptions("/tmp/badger"))
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Use Update function for transactional operation
	err = db.Update(func(txn *badger.Txn) error {
		_, err := txn.Get([]byte("lock"))
		if err == badger.ErrKeyNotFound {
			// Key not found, create the lock
			e := badger.NewEntry([]byte("lock"), []byte("locked"))
			err = txn.SetEntry(e)
			if err != nil {
				return err
			}
			log.Println("Lock created")
		} else if err != nil {
			return err
		} else {
			log.Println("Lock already exists")
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Failed to create lock: %v", err)
	}

	return nil
}

// Remove the lock to allow other instances of the application to run
func RemoveLock() error {
	// Open the Badger database located in the specified directory
	db, err := badger.Open(badger.DefaultOptions("/tmp/badger"))
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Use Update function for transactional operation
	err = db.Update(func(txn *badger.Txn) error {
		err := txn.Delete([]byte("lock"))
		if err != nil {
			return err
		}
		log.Println("Lock removed")
		return nil
	})

	if err != nil {
		log.Fatalf("Failed to remove lock: %v", err)
	}

	return nil
}
