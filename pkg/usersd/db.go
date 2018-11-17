// Copyright 2018 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the WTFPL.

package usersd

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"

	"github.com/blevesearch/bleve"
	"github.com/dgraph-io/badger"
)

// Backup writes a database backup to the given io.Writer.
func Backup(w io.Writer) error {
	if _, err := db.Backup(w, 0); err != nil {
		l.Printf(lError+" Can't create the database backup -> %v", err)
		return err
	}

	l.Print(lInfo + " Backup created")
	return nil
}

// Restore reads a database backup from the given io.Reader.
func Restore(r io.Reader) error {
	err := db.Load(r)

	if err != nil {
		l.Printf(lError+" Can't restore the database backup -> %v", err)
	} else {
		l.Print(lInfo + " Backup retored")
	}

	return err
}

// Reset loads default data to the database.
func Reset() error {
	l.Print(lInfo + " Truncating database..")
	empty := bytes.NewBuffer(nil)

	if err := Restore(empty); err != nil {
		return err
	}

	if err := admin.Save(); err != nil {
		l.Printf(lError+" Can't create the administrator user -> %v", err)
		return err
	}

	l.Print(lInfo + " Database truncated")
	return nil
}

// dbOpen opens/creates database and indexing directories. It receives a
// string as argument, which is the path to store the data, if an empty string
// is given, a temporary storage will be used.
func dbOpen(dir string) (err error) {
	if dir == "" {
		if dir, err = ioutil.TempDir("", "usersd"); err != nil {
			l.Printf(lFatal+" Can't create a temporary directory -> %v", err)
			return err
		}

		l.Print(lWarn + " Using a temporary database directory")
	}

	l.Printf(lInfo+" Database directory: %s\n", dir)

	indexPath := dir + "/search/users"
	index, err = bleve.Open(indexPath)

	if err == bleve.Error(1) { // ErrorIndexPathDoesNotExist
		mapping := bleve.NewIndexMapping()
		index, err = bleve.New(indexPath, mapping)

		if err != nil {
			l.Printf(lFatal+" Can't create the search index -> %v", err)
			return err
		}
	} else if err != nil {
		l.Printf(lFatal+" Can't open the search index -> %v", err)
		return err
	}

	dbPath := dir + "/data/users"
	opts := badger.DefaultOptions
	opts.Dir = dbPath
	opts.ValueDir = dbPath

	if db, err = badger.Open(opts); err != nil {
		if err := os.MkdirAll(dir+"/data", 0700); err != nil {
			l.Printf(lFatal+" Can't create the data folder -> %v", err)
			return err
		}

		if db, err = badger.Open(opts); err != nil {
			l.Printf(lFatal+" Can't open/create the database -> %+v", err)
			return err
		}

		Reset()
	}

	return nil
}

// dbClose closes the database and the search index.
func dbClose() error {
	db.Close()
	l.Print(lInfo + " Database closed")

	_, kvs, err := index.Advanced()

	if err != nil {
		l.Printf(lError+" Can't get the search index database -> %v", err)
		return err
	}

	if err := kvs.Close(); err != nil {
		l.Printf(lError+" Can't close the search index database -> %v", err)
		return err
	}

	l.Print(lInfo + " Search index closed")
	return nil
}
