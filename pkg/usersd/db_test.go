// Copyright 2018 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the WTFPL.

package usersd_test

import (
	"bytes"
	"testing"

	"github.com/ntrrg/usersd/pkg/usersd"
)

func TestBackup(t *testing.T) {
	if err := usersd.Init(Opts); err != nil {
		t.Fatal(err)
	}

	defer usersd.Close()
	backup := bytes.NewBuffer(nil)

	if err := usersd.Backup(backup); err != nil {
		t.Error(err)
	}
}
