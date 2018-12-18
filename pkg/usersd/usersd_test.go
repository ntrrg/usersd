// Copyright 2018 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the WTFPL.

package usersd_test

import (
	"fmt"
	"os"

	"github.com/ntrrg/usersd/pkg/usersd"
)

func ExampleNew() {
	// New database

	opts := usersd.DefaultOptions
	opts.Database = "test-db"
	defer os.RemoveAll(opts.Database)

	ud, err := usersd.New(opts)
	if err != nil {
		// Error handling
		return
	}

	// Your code here

	// Can't be deferred because the next example uses the same database
	ud.Close()
	if err := ud.Err(); err != nil {
		// Error handling
		return
	}

	// --------------------------------------------

	// Existing database

	ud2, err := usersd.New(opts)
	if err != nil {
		// Error handling
		return
	}

	defer ud2.Close()

	// Your code here

	// --------------------------------------------

	fmt.Println(ud2.IsTemp())

	// Output: false
}

func ExampleNew_temporaryStorage() {
	// Temporary storage

	opts := usersd.DefaultOptions
	ud, err := usersd.New(opts)
	if err != nil {
		// Error handling
		return
	}

	defer ud.Close()

	// Your code here

	fmt.Println(ud.IsTemp())

	// Output: true
}
