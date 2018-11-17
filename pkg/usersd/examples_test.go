// Copyright 2018 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the WTFPL.

package usersd_test

import (
	"fmt"
	"log"
	"os"

	"github.com/ntrrg/usersd/pkg/usersd"
)

func ExampleInit() {
	// Temporary storage

	opts := usersd.DefaultOptions

	if err := usersd.Init(opts); err != nil {
		fmt.Println(err)
	}

	// Your code here

	if err := usersd.Close(); err != nil {
		fmt.Println(err)
	}

	// --------------------------------------------

	// New database

	opts.Database = "test-db"
	defer os.RemoveAll(opts.Database)

	if err := usersd.Init(opts); err != nil {
		fmt.Println(err)
	}

	// Your code here

	if err := usersd.Close(); err != nil {
		fmt.Println(err)
	}

	// --------------------------------------------

	// Existing database

	if err := usersd.Init(opts); err != nil {
		fmt.Println(err)
	}

	// Your code here

	if err := usersd.Close(); err != nil {
		fmt.Println(err)
	}

	// --------------------------------------------

	// Output:
}

func ExampleInit_verbose() {
	opts := usersd.DefaultOptions
	opts.Verbose = true
	opts.Logger = log.New(os.Stdout, "", 0)
	opts.Database = "test-db"
	defer os.RemoveAll(opts.Database)

	if err := usersd.Init(opts); err != nil {
		fmt.Println(err)
	}

	// Your code here

	if err := usersd.Close(); err != nil {
		fmt.Println(err)
	}

	// Output:
	// [INFO][USERSD] Database directory: test-db
	// [INFO][USERSD] Truncating database..
	// [INFO][USERSD] Backup retored
	// [INFO][USERSD] Database truncated
	// [INFO][USERSD] Database closed
	// [INFO][USERSD] Search index closed
	// [INFO][USERSD] API closed
}

func ExampleInit_debug() {
	opts := usersd.DefaultOptions
	opts.Admin = ""
	opts.Debug = true

	if err := usersd.Init(opts); err != nil {
		fmt.Println(err)
	}

	// Your code here

	if err := usersd.Close(); err != nil {
		fmt.Println(err)
	}

	// Output:
}
