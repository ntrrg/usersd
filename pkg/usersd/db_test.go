// Copyright 2018 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the WTFPL.

package usersd_test

import (
	"testing"

	"github.com/ntrrg/usersd/pkg/usersd"
)

// TODO: Test backup creation.
func TestBackup(t *testing.T) {
}

// TODO: Test backup loading and search index updating.
func TestRestore(t *testing.T) {
}

func TestBL(t *testing.T) {
	bl := usersd.BL{}
	bl.Errorf("")
	bl.Infof("")
	bl.Warningf("")
}
