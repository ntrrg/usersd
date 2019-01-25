// Copyright 2018 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the WTFPL.

package usersd

import (
	"path/filepath"
	"runtime"

	"github.com/blevesearch/bleve"
	"github.com/dgraph-io/badger"
)

var usersd struct {
	opts  Options
	db    *badger.DB
	index bleve.Index

	running bool
}

// DefaultOptions are commonly used options for a simple Init call.
var DefaultOptions = Options{
	PasswdOpts: PasswordOptions{
		SaltSize: 32,
		Time:     1,
		Memory:   64 * 1024,
		Threads:  byte(runtime.GOMAXPROCS(0)),
		HashSize: 32,
	},

	JWTOpts: JWTOptions{
		Issuer: "usersd",
		Secret: "secret",
	},
}

// Options are parameters for initializing a service.
type Options struct {
	// Database location.
	Database string

	// Password authentication options.
	PasswdOpts PasswordOptions

	// JWT signing and verifying options.
	JWTOpts JWTOptions
}

// Init initializes the service.
func Init(opts Options) (err error) {
	dir := opts.Database

	if usersd.db, err = openDB(filepath.Join(dir, "data")); err != nil {
		return err
	}

	if usersd.index, err = openIndex(filepath.Join(dir, "search")); err != nil {
		return err
	}

	usersd.opts = opts
	usersd.running = true
	return nil
}

// Close terminates the service.
func Close() error {
	if err := usersd.db.Close(); err != nil {
		return err
	}

	_, kvs, err := usersd.index.Advanced()
	if err != nil {
		return err
	}

	if err = kvs.Close(); err != nil {
		return err
	}

	usersd.running = false
	return nil
}
