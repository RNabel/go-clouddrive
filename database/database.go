package database

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
)

const PATHS_BUCKET_NAME = "paths"

var TopLevelBuckets = [1]string{PATHS_BUCKET_NAME}

type CloudDB struct {
	db *bolt.DB
}

func (cdb *CloudDB) OpenDB() error {
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

	// Assign the db field.
	cdb.db = db

	return err
}

func (cdb *CloudDB) CloseDB() {
	defer cdb.db.Close()
}

func (cdb *CloudDB) IsInitialised() bool {
	var isInitialised = false
	cdb.db.View(func(tx *bolt.Tx) error {
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

func (cdb *CloudDB) AddElementToBucket(bucket []byte, key []byte, value []byte) error {
	return cdb.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		err := b.Put(key, value)
		return err
	})
}

func (cdb *CloudDB) GetElement(bucket string, key string) []byte {
	var output = []byte{}

	// Look up element.
	cdb.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		v := b.Get([]byte(key))
		output = v
		return nil
	})

	return output
}

// Require a set of transactions.

// TODO Add/Update a file mapping to the db
// TODO Remove a file mapping

// TODO How are changes synced with the database.
