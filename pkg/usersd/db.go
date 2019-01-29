// Copyright 2018 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the WTFPL.

package usersd

import (
	"io"
	"os"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/analyzer/keyword"
	"github.com/dgraph-io/badger"
)

// Backup writes a database backup to the given io.Writer.
func Backup(w io.Writer) error {
	if _, err := usersd.db.Backup(w, 0); err != nil {
		return err
	}

	return nil
}

// Restore reads a database backup from the given io.Reader.
func Restore(r io.Reader) error {
	if err := usersd.db.Load(r); err != nil {
		return err
	}

	return reloadIndex()
}

// Tx wraps a complete context for doing user management operations.
type Tx struct {
	*badger.Txn
}

// NewTx creates a database transaction. If writable is true, the database will
// allow modifications.
func NewTx(writable bool) *Tx {
	return &Tx{
		Txn: usersd.db.NewTransaction(writable),
	}
}

// Get is a helper for Badger operations.
func (tx *Tx) Get(key []byte) ([]byte, error) {
	item, err := tx.Txn.Get(key)
	if err != nil {
		return nil, err
	}

	return item.ValueCopy(nil)
}

func cleanIndex() error {
	bq := bleve.NewMatchAllQuery()
	req := bleve.NewSearchRequest(bq)
	res, err := usersd.index.Search(req)
	if err != nil {
		return err
	}

	for _, hit := range res.Hits {
		if err = usersd.index.Delete(hit.ID); err != nil {
			return err
		}
	}

	return nil
}

func openDB(dir string) (*badger.DB, error) {
	opts := badger.DefaultOptions
	opts.Dir = dir
	opts.ValueDir = dir

	db, err := badger.Open(opts)
	if err != nil {
		if err = os.MkdirAll(dir, 0700); err != nil {
			return nil, err
		}

		return badger.Open(opts)
	}

	return db, nil
}

func openIndex(dir string) (bleve.Index, error) {
	index, err := bleve.Open(dir)
	if err == bleve.Error(1) { // ErrorIndexPathDoesNotExist
		keywordField := bleve.NewTextFieldMapping()
		keywordField.Analyzer = keyword.Name

		users := bleve.NewDocumentMapping()
		users.AddFieldMappingsAt("documenttype", keywordField)
		users.AddFieldMappingsAt("id", keywordField)
		users.AddFieldMappingsAt("mode", keywordField)
		users.AddFieldMappingsAt("email", keywordField)
		users.AddFieldMappingsAt("phone", keywordField)

		mapping := bleve.NewIndexMapping()
		mapping.TypeField = "documenttype"
		mapping.AddDocumentMapping(UsersDI, users)

		return bleve.New(dir, mapping)
	} else if err != nil {
		return nil, err
	}

	return index, nil
}

func reloadIndex() error {
	if err := cleanIndex(); err != nil {
		return err
	}

	tx := NewTx(true)
	defer tx.Discard()

	users, err := tx.GetUsers("")
	if err != nil {
		return err
	}

	for _, user := range users {
		if err = tx.WriteUser(user); err != nil {
			return err
		}
	}

	return nil
}

type bl struct{}

func (l *bl) Errorf(f string, v ...interface{})   {}
func (l *bl) Infof(f string, v ...interface{})    {}
func (l *bl) Warningf(f string, v ...interface{}) {}

func init() {
	var badgerLogger = &bl{}
	badger.SetLogger(badgerLogger)
}
