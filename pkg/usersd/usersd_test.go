// Copyright 2018 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the WTFPL.

package usersd_test

import (
	"os"
	"path/filepath"
	"testing"

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

func TestInit_forbiddenDB(t *testing.T) {
	dir := filepath.Join(testDir, "init-forbidden-db")
	if err := os.Mkdir(dir, 0700); err != nil {
		t.Fatal(err)
	}

	if err := os.MkdirAll(filepath.Join(dir, "data"), 0000); err != nil {
		t.Fatal(err)
	}

	opts := usersd.DefaultOptions
	opts.Database = dir
	if err := usersd.Init(opts); err == nil {
		t.Error("Service initialized on forbidden data path")
	}
}

func TestInit_forbiddenIndex(t *testing.T) {
	dir := filepath.Join(testDir, "init-forbidden-index")
	if err := os.Mkdir(dir, 0700); err != nil {
		t.Fatal(err)
	}

	if err := os.MkdirAll(filepath.Join(dir, "search"), 0000); err != nil {
		t.Fatal(err)
	}

	opts := usersd.DefaultOptions
	opts.Database = dir
	if err := usersd.Init(opts); err == nil {
		t.Error("Service initialized on forbidden search index path")
	}
}
