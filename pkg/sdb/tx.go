// Copyright 2019 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the MIT license.

package sdb

import (
	"encoding/json"
	"errors"

	"github.com/blevesearch/bleve"
	"github.com/dgraph-io/badger/v2"
)

const (
	RW = true
	RO = false
)

// Tx is a transaction object which provides data management methods. The
// search index doesn't support transactions yet, so indexing operations just
// take effect after committing the transaction.
type Tx struct {
	db   *DB
	dbTx *badger.Txn
	si   bleve.Index
	rw   bool

	// Search index operations to be done when the transaction is committed.
	operations map[string]interface{}
}

// NewTx creates a database transaction. If rw is false, the new transaction
// will be read-only.
func (db *DB) NewTx(rw bool) *Tx {
	return &Tx{
		db:   db,
		dbTx: db.db.NewTransaction(rw),
		si:   db.si,
		rw:   rw,
	}
}

// Commit writes the transaction operations to the database. If a Bleve error
// is returned, the search index should be reloaded (see DB.ReloadIndex), keep
// the amount of operations per transaction low to avoid this.
func (tx *Tx) Commit() error {
	if err := tx.dbTx.Commit(); err != nil {
		return badgerError(err)
	}

	for id, data := range tx.operations {
		if data != nil {
			if err := tx.si.Index(id, data); err != nil {
				return bleveError(err)
			}
		} else {
			if err := tx.si.Delete(id); err != nil {
				return bleveError(err)
			}
		}
	}

	return nil
}

// Delete deletes the given key. This operation happens in memory, it will be
// written to the database once Commit is called.
func (tx *Tx) Delete(key []byte) error {
	if err := tx.dbTx.Delete(key); err != nil {
		return badgerError(err)
	}

	if tx.operations == nil {
		tx.operations = make(map[string]interface{})
	}

	tx.operations[string(key)] = nil

	return nil
}

// Discard drops all the pending modifications and set the transactions as
// discarded.
func (tx *Tx) Discard() {
	if tx.rw {
		tx.operations = nil
	}

	tx.dbTx.Discard()
}

// Find fetches the keys from the values that satisfies the given constraints.
// See http://blevesearch.com/docs/Query-String-Query/ for more info about the
// the query language syntax. sort is a list of field names used for sorting,
// any field prefixed by a hyphen (-) will user reverse order.
func (tx *Tx) Find(q string, sort ...string) ([][]byte, error) {
	if q == "" && len(sort) == 0 {
		return nil, nil
	}

	result := [][]byte{}
	bq := bleve.NewQueryStringQuery(q)
	req := bleve.NewSearchRequest(bq)
	req.SortBy(sort)

	res, err := tx.si.Search(req)
	if err != nil {
		return nil, bleveError(err)
	}

	for _, hit := range res.Hits {
		result = append(result, []byte(hit.ID))
	}

	return result, nil
}

// Get reads the value from the given key and decodes it into v (must be a
// pointer).
func (tx *Tx) Get(key []byte, v interface{}) error {
	item, err := tx.dbTx.Get(key)
	if errors.Is(err, badger.ErrKeyNotFound) {
		return ErrKeyNotFound
	} else if err != nil {
		return badgerError(err)
	}

	buf := tx.db.buffers.Get()
	defer tx.db.buffers.Add(buf)

	data, err := item.ValueCopy(buf.Bytes())
	if err != nil {
		return err
	}

	return json.Unmarshal(data, v)
}

// Set set val as value of the given key. This operation happens in memory, it
// will be written to the database once Commit is called. Since values are
// encoded with JSON, struct tags will define what fields are stored.
func (tx *Tx) Set(key []byte, val interface{}) error {
	buf := tx.db.buffers.Get()
	defer tx.db.buffers.Add(buf)

	e := json.NewEncoder(buf)
	if err := e.Encode(val); err != nil {
		return err
	}

	if err := tx.dbTx.Set(key, buf.Bytes()); err != nil {
		return badgerError(err)
	}

	if tx.operations == nil {
		tx.operations = make(map[string]interface{})
	}

	tx.operations[string(key)] = val

	return nil
}
