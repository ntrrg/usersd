// Copyright 2018 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the WTFPL.

package usersd

import (
	"runtime"

	"github.com/blevesearch/bleve"
	"github.com/dgraph-io/badger"
)

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

// Service is an authentication and authorization service.
type Service struct {
	opts  Options
	db    *badger.DB
	index bleve.Index

	running bool
}

// Create creates a service, but doesn't initialize it.
func Create(opts Options) *Service {
	return &Service{
		opts: opts,
	}
}

// New creates and initialize a service.
func New(opts Options) (*Service, error) {
	s := Create(opts)

	if err := s.Start(); err != nil {
		return nil, err
	}

	return s, nil
}

// Close terminates the service (databases, search indexes, etc...).
func (s *Service) Close() error {
	if !s.running {
		return nil
	}

	s.running = false
	return s.closeDB()
}

// Start initialize the service (databases, search indexes, etc...).
func (s *Service) Start() error {
	s.running = true
	return s.openDB()
}
