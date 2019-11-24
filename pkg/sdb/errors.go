// Copyright 2019 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the MIT license.

package sdb

import (
	"fmt"
	"strings"

	"github.com/dgraph-io/badger/v2"
)

// Errors
var (
	ErrKeyNotFound = badgerError(badger.ErrKeyNotFound)
)

// IsBadgerError returns true if the given error is from Badger.
func IsBadgerError(err error) bool {
	return errorContains(err, "badger: ")
}

// IsBleveError returns true if the given error is from Bleve.
func IsBleveError(err error) bool {
	return errorContains(err, "bleve: ")
}

func badgerError(err error) error {
	return fmt.Errorf("badger: %w", err)
}

func bleveError(err error) error {
	return fmt.Errorf("bleve: %w", err)
}

func errorContains(err error, s string) bool {
	if err == nil {
		return false
	}

	return strings.Contains(err.Error(), s)
}
