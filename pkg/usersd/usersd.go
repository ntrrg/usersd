// Copyright 2018 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the WTFPL.

package usersd

import (
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/blevesearch/bleve"
	"github.com/dgraph-io/badger"
)

// Logger prefixes.
const (
	lDebug = "[DEBUG][USERSD]"
	lInfo  = "[INFO][USERSD]"
	lWarn  = "[WARN][USERSD]"
	lError = "[ERROR][USERSD]"
	lFatal = "[FATAL][USERSD]"
)

var (
	admin      *User
	db         *badger.DB
	index      bleve.Index
	bcryptCost int
	l          *log.Logger
	debug      bool
)

// Options are parameters for initializing the API.
type Options struct {
	// Administrator user. Format: "ID[:PASSWORD]".
	Admin string

	// Database location.
	Database string

	// Password hashing strength. See https://godoc.org/golang.org/x/crypto/bcrypt#GenerateFromPassword.
	HashingStrength int

	// Debugging mode.
	Verbose bool
	Logger  *log.Logger
	Debug   bool
}

// DefaultOptions are the commonly used options for a simple Init call.
var DefaultOptions = Options{
	Admin:           "admin:admin",
	HashingStrength: 10,
}

// Init sets the API up. It receives an Options instance as argument and
// returns an error if any.
func Init(opts Options) (err error) {
	if opts.Debug {
		debug = true
		opts.Verbose = true
	}

	if opts.Verbose && opts.Logger != nil {
		l = opts.Logger
	} else if opts.Verbose && opts.Logger == nil {
		l = log.New(os.Stderr, "", log.LstdFlags)
	} else {
		l = log.New(ioutil.Discard, "", log.LstdFlags)
	}

	bcryptCost = opts.HashingStrength

	var username, password string
	creds := strings.SplitN(opts.Admin, ":", 2)

	if creds[0] == "" {
		username, password = "admin", "admin"
	} else if len(creds) == 1 {
		username, password = creds[0], creds[0]
	} else {
		username, password = creds[0], creds[1]
	}

	admin = NewUser("admin", map[string]interface{}{
		"username": username,
		"password": password,
		"name":     "Administrator",
	})

	if err := dbOpen(opts.Database); err != nil {
		return err
	}

	badger.SetLogger(badgerLogger)
	return nil
}

// Close terminates with the API processes (database, search index, etc...).
func Close() error {
	if err := dbClose(); err != nil {
		return err
	}

	l.Print(lInfo + " API closed")
	return nil
}

// Badger logger

type bL struct {
	*log.Logger
}

func (l *bL) Errorf(f string, v ...interface{}) {
	l.Printf("[ERROR][BADGER] "+f, v...)
}

func (l *bL) Infof(f string, v ...interface{}) {
	l.Printf("[INFO][BADGER] "+f, v...)
}

func (l *bL) Warningf(f string, v ...interface{}) {
	l.Printf("[WARN][BADGER] "+f, v...)
}

var badgerLogger = &bL{Logger: log.New(ioutil.Discard, "", log.LstdFlags)}
