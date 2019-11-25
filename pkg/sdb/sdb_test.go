// Copyright 2019 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the MIT license.

package sdb_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"nt.web.ve/go/usersd/pkg/sdb"
)

var testDir string

func TestOpen_existing(t *testing.T) {
	dir := filepath.Join(testDir, "open-existing")
	if err := os.Mkdir(dir, 0755); err != nil {
		t.Fatal(err)
	}

	run := func() {
		db, err := sdb.Open(dir)
		if err != nil {
			t.Fatal(err)
		}

		defer db.Close()
	}

	run()
	run()
}

func TestOpen_badTemporaryFS(t *testing.T) {
	tmpDir := "TMPDIR"
	defer os.Setenv(tmpDir, os.Getenv(tmpDir))

	if err := os.Setenv(tmpDir, "/non/existent/directory"); err != nil {
		t.Fatal(err)
	}

	db, err := sdb.Open(sdb.InMemory)
	if err == nil {
		defer db.Close()
		t.Error("Database initialized on invalid temporary filesystem")
	}
}

func TestOpen_forbiddenDB(t *testing.T) {
	dir := filepath.Join(testDir, "open-forbidden-db")
	if err := os.Mkdir(dir, 0755); err != nil {
		t.Fatal(err)
	}

	err := os.Mkdir(filepath.Join(dir, sdb.DatabaseDir), 0000)
	if err != nil {
		t.Fatal(err)
	}

	db, err := sdb.Open(dir)
	if err == nil {
		defer db.Close()
		t.Error("Database initialized on forbidden database path")
	}

	if !sdb.IsBadgerError(err) {
		t.Error("The error should come from Badger", err)
	}
}

func TestOpen_forbiddenIndex(t *testing.T) {
	dir := filepath.Join(testDir, "open-forbidden-index")
	if err := os.Mkdir(dir, 0755); err != nil {
		t.Fatal(err)
	}

	err := os.Mkdir(filepath.Join(dir, sdb.SearchIndexDir), 0000)
	if err != nil {
		t.Fatal(err)
	}

	db, err := sdb.Open(dir)
	if err == nil {
		defer db.Close()
		t.Error("Database initialized on forbidden search index path")
	}

	if !sdb.IsBleveError(err) {
		t.Error("The error should come from Bleve", err)
	}
}

func TestOpenWith_fillBufferPool(t *testing.T) {
	opts, err := sdb.MemoryOptions()
	if err != nil {
		t.Fatal(err)
	}

	opts.BufferPoolSize = 5
	opts.BufferPoolFill = true

	db, err := sdb.OpenWith(opts)
	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()
}

func init() {
	var err error

	testDir, err = ioutil.TempDir("", "sdb-tests")
	if err != nil {
		panic(err)
	}
}
