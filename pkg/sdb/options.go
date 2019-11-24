// Copyright 2019 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the MIT license.

package sdb

import (
	"log"
	"os"
	"path/filepath"

	"github.com/blevesearch/bleve/mapping"
	"github.com/dgraph-io/badger/v2"
)

const (
	DatabaseDir    = "database"
	SearchIndexDir = "search-index"
)

type BleveOptions struct {
	Dir          string
	DoctypeField string
	DocMappings  map[string]*mapping.DocumentMapping
}

// Options are parameters for initializing a database.
type Options struct {
	// Database location.
	Directory string

	Badger badger.Options
	Bleve  BleveOptions

	BufferPoolSize     int  // Amount of buffers.
	BufferPoolMaxBytes int  // Bytes limit per buffer.
	BufferPoolFill     bool // Fill up the pool at DB creation.

	Logger *log.Logger
}

// DefaultOptions returns commonly used options for creating a database.
func DefaultOptions(dir string) Options {
	dir = filepath.Clean(dir)

	return Options{
		Directory: dir,
		Badger:    badger.DefaultOptions(filepath.Join(dir, DatabaseDir)),

		Bleve: BleveOptions{
			Dir:          filepath.Join(dir, SearchIndexDir),
			DoctypeField: "Doctype",
			DocMappings:  make(map[string]*mapping.DocumentMapping),
		},

		BufferPoolSize:     500,
		BufferPoolMaxBytes: 5 * 1024,
		BufferPoolFill:     false,

		Logger: log.New(os.Stderr, "sdb: ", log.LstdFlags),
	}
}
