// Copyright 2019 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the MIT license.

package sdb

import (
	"io/ioutil"
	"log"

	"github.com/blevesearch/bleve"
	"github.com/dgraph-io/badger/v2"
	"github.com/dgraph-io/badger/v2/options"
	ntbytes "nt.web.ve/go/ntgo/bytes"
)

// DB is a database object which provides database management methods, for data
// management see Tx.
type DB struct {
	opts Options

	db *badger.DB
	si bleve.Index

	buffers *ntbytes.BufferPool
	logger  *log.Logger
}

// Open initializes a database in the given directory.
func Open(dir string) (*DB, error) {
	var err error

	if dir == "" {
		dir, err = ioutil.TempDir("", "sdb")
		if err != nil {
			return nil, err
		}
	}

	opts := DefaultOptions(dir)

	return OpenWith(opts)
}

// OpenWith initializes a database with the given options.
func OpenWith(opts Options) (*DB, error) {
	var err error

	db := new(DB)

	db.buffers = ntbytes.NewBufferPool(
		opts.BufferPoolSize,
		opts.BufferPoolMaxBytes,
	)

	if opts.BufferPoolFill {
		db.buffers.Fill()
	}

	db.logger = opts.Logger

	if db.db, err = openDB(opts.Badger); err != nil {
		return nil, badgerError(err)
	}

	if db.si, err = openSearchIndex(opts.Bleve); err != nil {
		return nil, bleveError(err)
	}

	db.opts = opts

	return db, nil
}

// Close terminates the database.
func (db *DB) Close() error {
	if err := db.db.Close(); err != nil {
		return badgerError(err)
	}

	_, kvs, err := db.si.Advanced()
	if err != nil {
		return bleveError(err)
	}

	if err = kvs.Close(); err != nil {
		return bleveError(err)
	}

	return nil
}

func openDB(opts badger.Options) (*badger.DB, error) {
	opts.Logger = &bl{}
	opts.Compression = options.Snappy

	return badger.Open(opts)
}

func openSearchIndex(opts BleveOptions) (bleve.Index, error) {
	index, err := bleve.Open(opts.Dir)
	if err == bleve.Error(1) { // ErrorIndexPathDoesNotExist
		mapping := bleve.NewIndexMapping()
		mapping.TypeField = opts.DoctypeField

		for t, m := range opts.DocMappings {
			mapping.AddDocumentMapping(t, m)
		}

		return bleve.New(opts.Dir, mapping)
	}

	return index, err
}

type bl struct{}

func (l *bl) Errorf(f string, v ...interface{})   {}
func (l *bl) Warningf(f string, v ...interface{}) {}
func (l *bl) Infof(f string, v ...interface{})    {}
func (l *bl) Debugf(f string, v ...interface{})   {}