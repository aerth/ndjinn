package database

import (
	"encoding/json"
	"log"

	"github.com/boltdb/bolt"
	_ "github.com/go-sql-driver/mysql"
)

var (
	BoltDB    *bolt.DB     // Bolt wrapper
	databases DatabaseInfo // Database info
)

type DatabaseType string

const (
	TypeBolt DatabaseType = "Bolt"
)

type DatabaseInfo struct {
	Type DatabaseType

	Bolt BoltInfo
}

// BoltInfo is the details for the database connection
type BoltInfo struct {
	Path string
}

// Connect to the database
func Connect(d DatabaseInfo) {
	var err error

	// Store the config
	databases = d

	switch d.Type {

	case TypeBolt:
		// Connect to Bolt
		if BoltDB, err = bolt.Open(d.Bolt.Path, 0600, nil); err != nil {
			log.Println("Bolt Driver Error", err)
		}

	default:
		log.Println("No registered database in config")
	}
}

// Update makes a modification to Bolt
func Update(bucket_name string, key string, dataStruct interface{}) error {
	err := BoltDB.Update(func(tx *bolt.Tx) error {
		// Create the bucket
		bucket, e := tx.CreateBucketIfNotExists([]byte(bucket_name))
		if e != nil {
			return e
		}

		// Encode the record
		encoded_record, e := json.Marshal(dataStruct)
		if e != nil {
			return e
		}

		// Store the record
		if e = bucket.Put([]byte(key), encoded_record); e != nil {
			return e
		}
		return nil
	})
	return err
}

// View retrieves a record in Bolt
func View(bucket_name string, key string, dataStruct interface{}) error {
	err := BoltDB.View(func(tx *bolt.Tx) error {
		// Get the bucket
		b := tx.Bucket([]byte(bucket_name))
		if b == nil {
			return bolt.ErrBucketNotFound
		}

		// Retrieve the record
		v := b.Get([]byte(key))
		if len(v) < 1 {
			return bolt.ErrInvalid
		}

		// Decode the record
		e := json.Unmarshal(v, &dataStruct)
		if e != nil {
			return e
		}

		return nil
	})

	return err
}

// ReadConfig returns the database information
func ReadConfig() DatabaseInfo {
	return databases
}

func Close() DatabaseInfo {
	BoltDB.Close()
	return databases
}
