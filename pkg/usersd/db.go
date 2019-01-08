// Copyright 2018 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the WTFPL.

package usersd

import (
	"io/ioutil"
	"os"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/analyzer/keyword"
	"github.com/dgraph-io/badger"
)

// Backup writes a database backup to the given io.Writer. Returns an error if
// any.
// func (s *Service) Backup(w io.Writer) error {
// 	if _, err := s.db.Backup(w, 0); err != nil {
// 		return err
// 	}
//
// 	return nil
// }

// Restore reads a database backup from the given io.Reader. Returns an error
// if any.
// func (s *Service) Restore(r io.Reader) error {
// 	if err := s.db.Load(r); err != nil {
// 		return err
// 	}
//
// 	return nil
// }

// Tx wraps a badger.Txn and a bleve.Index.
type Tx struct {
	*badger.Txn
	Index bleve.Index
}

// NewTx creates a database transaction. If writable is true, the database will
// allow modifications.
func (s *Service) NewTx(writable bool) *Tx {
	return &Tx{
		Txn:   s.db.NewTransaction(writable),
		Index: s.index,
	}
}

// closeDB closes the database and the search index. Returns an error if any.
func (s *Service) closeDB() error {
	if err := s.db.Close(); err != nil {
		return err
	}

	_, kvs, err := s.index.Advanced()
	if err != nil {
		return err
	}

	return kvs.Close()
}

// openDB opens/creates database and indexing directories. Return an error if
// any.
func (s *Service) openDB() (err error) {
	dir := s.opts.Database
	if dir == "" {
		if dir, err = ioutil.TempDir("", "usersd"); err != nil {
			return err
		}
	}

	badger.SetLogger(badgerLogger)

	if s.db, err = openDB(dir + "/data"); err != nil {
		return err
	}

	if s.index, err = openIndex(dir + "/search"); err != nil {
		s.err = s.closeDB()
		return err
	}

	return nil
}

func openDB(dir string) (*badger.DB, error) {
	opts := badger.DefaultOptions
	opts.Dir = dir
	opts.ValueDir = dir

	db, err := badger.Open(opts)
	if err != nil {
		if err = os.MkdirAll(dir+"/data", 0700); err != nil {
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
		keywordField := bleve.NewTextFieldMapping()
		keywordField.Analyzer = keyword.Name

		disabledField := bleve.NewDocumentDisabledMapping()

		users := bleve.NewDocumentMapping()
		users.AddFieldMappingsAt("documenttype", keywordField)
		users.AddFieldMappingsAt("id", keywordField)
		users.AddFieldMappingsAt("mode", keywordField)
		users.AddFieldMappingsAt("email", keywordField)
		users.AddFieldMappingsAt("phone", keywordField)
		users.AddSubDocumentMapping("password", disabledField)

		mapping := bleve.NewIndexMapping()
		mapping.AddDocumentMapping(UsersDT, users)
		mapping.TypeField = "documenttype"

		index, err = bleve.New(dir, mapping)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	return index, nil
}

type bl struct{}

func (l *bl) Errorf(f string, v ...interface{})   {}
func (l *bl) Infof(f string, v ...interface{})    {}
func (l *bl) Warningf(f string, v ...interface{}) {}

var badgerLogger = &bl{}
