// Copyright 2019 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the MIT license.

package sdb_test

import (
	"testing"

	"nt.web.ve/go/usersd/pkg/sdb"
)

func TestIsBadgerError_nilError(t *testing.T) {
	sdb.IsBadgerError(nil)
}

func TestIsBleveError_nilError(t *testing.T) {
	sdb.IsBleveError(nil)
}
