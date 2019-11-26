// Copyright 2019 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the MIT license.

package sdb

import (
	"github.com/blevesearch/bleve"
	"github.com/dgraph-io/badger/v2"
)

type DecoderFunc func(tx *Tx, key []byte) (interface{}, error)

// ReloadIndex recreates the search index, it takes a decoder function as
// argument, this is necessary since it is not possible to decode one type into
// another.
func (db *DB) ReloadIndex(f DecoderFunc) error {
	if err := db.cleanIndex(); err != nil {
		return err
	}

	tx := db.NewTx(RO)
	defer tx.Discard()

	it := tx.dbTx.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()

	for it.Rewind(); it.Valid(); it.Next() {
		key := it.Item().Key()

		val, err := f(tx, key)
		if err != nil {
			return err
		}

		if err := tx.si.Index(string(key), val); err != nil {
			return bleveError(err)
		}
	}

	return nil
}

func (db *DB) cleanIndex() error {
	bq := bleve.NewMatchAllQuery()
	req := bleve.NewSearchRequest(bq)

	res, err := db.si.Search(req)
	if err != nil {
		return bleveError(err)
	}

	for _, hit := range res.Hits {
		if err = db.si.Delete(hit.ID); err != nil {
			return bleveError(err)
		}
	}

	return nil
}
