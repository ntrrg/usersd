// Copyright 2018 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the WTFPL.

package usersd

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"

	"github.com/blevesearch/bleve"
	"github.com/dgraph-io/badger"
)

var (
	db    *badger.DB
	index bleve.Index
)

// Backup writes a database backup to the given io.Writer.
func Backup(w io.Writer) error {
	_, err := db.Backup(w, 0)
	return err
}

// Close terminates with the database process.
//
// BUG(ntrrg): Since Close doesn't closes the index, testing Init with an
// existing directory is not possible.
func Close() {
	db.Close()
}

// Get is a helper for doing a simple database read.
func Get(key []byte) ([]byte, error) {
	txn := db.NewTransaction(false)
	defer txn.Discard()

	item, err := txn.Get(key)

	if err != nil {
		return nil, err
	}

	v, err := item.ValueCopy(nil)
	return v, err
}

// Init opens/creates database and indexing directories. It receives a string
// as argument, which is the path to store the data, if an empty string is
// given, a temporary storage will be used.
func Init(dir string) (err error) {
	if dir == "" {
		if dir, err = ioutil.TempDir("", "backend"); err != nil {
			return err
		}

		log.Printf("[INFO][API] Using temporary database directory at %s", dir)
	}

	indexPath := dir + "/index"
	index, err = bleve.Open(indexPath)

	if err == bleve.Error(1) { // ErrorIndexPathDoesNotExist
		mapping := bleve.NewIndexMapping()
		index, err = bleve.New(indexPath, mapping)

		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	dbPath := dir + "/database"
	opts := badger.DefaultOptions
	opts.Dir = dbPath
	opts.ValueDir = dbPath

	if db, err = badger.Open(opts); err != nil {
		return err
	}

	if err := Reset(); err != nil {
		db.Close()
		return err
	}

	return nil
}

// Reset loads default data to the database.
func Reset() error {
	empty := bytes.NewBuffer(nil)

	if err := Restore(empty); err != nil {
		return err
	}

	users := []*User{
		{
			ID:        "18dd75e9-3d4a-48e2-bafc-3c8f95a8f0d1",
			Name:      "John",
			HighScore: 322,
		},
		{
			ID:        "f9a9af78-6681-4d7d-8ae7-fc41e7a24d08",
			Name:      "Bob",
			HighScore: 21,
		},
		{
			ID:        "2d18862b-b9c3-40f5-803e-5e100a520249",
			Name:      "Alice",
			HighScore: 99332,
		},
	}

	for _, user := range users {
		if err := user.Save(); err != nil {
			log.Println(err)

			if err := Restore(empty); err != nil {
				return err
			}

			return err
		}
	}

	return nil
}

// Restore reads a database backup from the given io.Reader.
func Restore(r io.Reader) error {
	return db.Load(r)
}

// Set is a helper for doing a simple database write.
func Set(k, v []byte) error {
	return db.Update(func(txn *badger.Txn) error {
		if err := txn.Set(k, v); err != nil {
			return err
		}

		if err := txn.Commit(nil); err != nil {
			return err
		}

		return nil
	})
}
