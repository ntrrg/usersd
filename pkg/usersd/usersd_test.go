// Copyright 2018 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the WTFPL.

package usersd_test

import (
	"path/filepath"

	"github.com/ntrrg/usersd/pkg/usersd"
)

func ExampleInit() {
	// New database

	opts := usersd.DefaultOptions
	opts.Database = filepath.Join(testDir, "example-new")
	if err := usersd.Init(opts); err != nil {
		// Error handling
		return
	}

	// Can't be deferred because the code below uses the same database.
	// defer usersd.Close()

	// Your code here

	if err := usersd.Close(); err != nil {
		return
	}

	// --------------------------------------------

	// Existing database

	if err := usersd.Init(opts); err != nil {
		// Error handling
		return
	}

	defer usersd.Close()

	// Your code here

	// Output:
}
