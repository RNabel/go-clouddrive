package database

import (
	"github.com/boltdb/bolt"
	"log"
	"fmt"
)

const PATHS_BUCKET_NAME = "paths"

var TopLevelBuckets = [1]string{PATHS_BUCKET_NAME}

func OpenDB() (*bolt.DB, error) {
	db, err := bolt.Open("files.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Create all top level buckets.
	err = db.Update(func(tx *bolt.Tx) error {
		for _, name := range TopLevelBuckets {
			_, err := tx.CreateBucketIfNotExists([]byte(name))
			if err != nil {
				return fmt.Errorf("create bucket %s", err)
			}
		}
		return nil
	})

	return db, err
}

func CloseDB(db *bolt.DB) {
	defer db.Close()
}

func IsDbInitialised(db *bolt.DB) bool {
	var isInitialised = false
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(PATHS_BUCKET_NAME))

		// Go over each key in the bucket.
		var counter = 0
		b.ForEach(func(k, v []byte) error {
			counter += 1
			return nil
		})

		isInitialised = counter != 0

		return nil
	})
	return isInitialised
}

func AddElementToBucket(db *bolt.DB, bucket []byte, key []byte, value []byte) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		err := b.Put(key, value)
		return err
	})
}

// Require a set of transactions.

// TODO Add/Update a file mapping to the db
// TODO Remove a file mapping
// TODO How are they retrieved?

// TODO How are changes synced with the database.
