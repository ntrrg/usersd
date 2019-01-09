// Copyright 2018 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the WTFPL.

package usersd

import (
	"runtime"

	"github.com/blevesearch/bleve"
	"github.com/dgraph-io/badger"
)

// DefaultOptions are the commonly used options for a simple Init call.
var DefaultOptions = Options{
	JWTSecret: "secret",
	PasswdOpts: PasswordOptions{
		SaltSize: 32,
		Time:     1,
		Memory:   64 * 1024,
		Threads:  byte(runtime.GOMAXPROCS(0)),
		HashSize: 32,
	},
}

// Options are parameters for initializing a service.
type Options struct {
	// Database location, if an empty string is given, a temporary storage will
	// be used.
	Database string

	// Password hashing options.
	PasswdOpts PasswordOptions

	// Secret for signing and verifying JWT.
	JWTSecret string
}

// Service is an authentication and authorization service.
type Service struct {
	opts  Options
	err   error
	db    *badger.DB
	index bleve.Index

	closed bool
}

// New creates and starts a service. Receives an Options instance as argument
// and returns a Service instance and an error if any.
func New(opts Options) (*Service, error) {
	s := new(Service)
	s.opts = opts

	if err := s.Start(); err != nil {
		return nil, err
	}

	return s, nil
}

// Close terminates the service (databases, search indexes, etc...). Any error
// closing the service will be stored at Service.err and will be accessible
// from Service.Err().
func (s *Service) Close() {
	if s.closed {
		return
	}

	s.err = s.closeDB()
	s.closed = true
}

// Err checks if any error occurred during some processes (closing, etc...).
func (s *Service) Err() error {
	return s.err
}

// IsTemp returns true if the service persistent storage is temporary.
func (s *Service) IsTemp() bool {
	return s.opts.Database == ""
}

// Start initialize the service (databases, search indexes, etc...). Returns an
// error if any.
func (s *Service) Start() error {
	return s.openDB()
}
