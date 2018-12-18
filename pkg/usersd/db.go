// Copyright 2018 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the WTFPL.

package usersd

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/blevesearch/bleve"
	"github.com/dgraph-io/badger"
)

// Index is a collection of Bleve search indexes.
type Index map[string]bleve.Index

// // Backup writes a database backup to the given io.Writer. Returns an error if
// // any.
// func (s *Service) Backup(w io.Writer) error {
// 	if _, err := s.DB.Backup(w, 0); err != nil {
// 		return err
// 	}
//
// 	return nil
// }
//
// // Restore reads a database backup from the given io.Reader. Returns an error
// // if any.
// func (s *Service) Restore(r io.Reader) error {
// 	if err := s.DB.Load(r); err != nil {
// 		return err
// 	}
//
// 	return nil
// }

// openDB opens/creates database and indexing directories. Return an error if
// any.
func (s *Service) openDB() (err error) {
	defer func() {
		if err != nil {
			s.closeDB()
		}
	}()

	dir := s.opts.Database
	if dir == "" {
		if dir, err = ioutil.TempDir("", "usersd"); err != nil {
			return err
		}
	}

	if s.DB, err = openDB(dir + "/data"); err != nil {
		return err
	}

	s.Index = make(Index)

	indexes := []string{
		"users",
	}

	for _, name := range indexes {
		if s.Index[name], err = openIndex(dir + "/search/" + name); err != nil {
			return err
		}
	}

	badger.SetLogger(badgerLogger)
	return nil
}

// closeDB closes the database and the search index. Returns an error if any.
func (s *Service) closeDB() error {
	if err := s.DB.Close(); err != nil {
		return err
	}

	for key, index := range s.Index {
		_, kvs, err := index.Advanced()
		if err != nil {
			return err
		}

		if err := kvs.Close(); err != nil {
			return err
		}

		delete(s.Index, key)
	}

	return nil
}

func openDB(dir string) (*badger.DB, error) {
	opts := badger.DefaultOptions
	opts.Dir = dir
	opts.ValueDir = dir

	db, err := badger.Open(opts)
	if err != nil {
		if err := os.MkdirAll(dir+"/data", 0700); err != nil {
			return nil, err
		}

		if db, err = badger.Open(opts); err != nil {
			return nil, err
		}
	}

	return db, nil
}

func openIndex(dir string) (bleve.Index, error) {
	index, err := bleve.Open(dir)
	if err == bleve.Error(1) { // ErrorIndexPathDoesNotExist
		mapping := bleve.NewIndexMapping()
		index, err = bleve.New(dir, mapping)

		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	return index, nil
}

// Badger logger

type bL struct {
	*log.Logger
}

func (l *bL) Errorf(f string, v ...interface{}) {
	l.Printf(f, v...)
}

func (l *bL) Infof(f string, v ...interface{}) {
	l.Printf(f, v...)
}

func (l *bL) Warningf(f string, v ...interface{}) {
	l.Printf(f, v...)
}

var badgerLogger = &bL{Logger: log.New(ioutil.Discard, "", log.LstdFlags)}
