// Copyright 2019 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the MIT license.

package sdb_test

import (
	"testing"

	"nt.web.ve/go/usersd/pkg/sdb"
)

func TestDB_ReloadIndex(t *testing.T) {
	db, err := initTest("reload-search-index")
	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()

	tx := db.NewTx(sdb.RW)
	defer tx.Discard()

	key := []byte("test")
	data := struct{ Data string }{Data: "lorem ipsum"}

	if errSet := tx.Set(key, data); err != nil {
		t.Fatal(errSet)
	}

	if errCommit := tx.Commit(); err != nil {
		t.Fatal(errCommit)
	}

	if errReload := db.ReloadIndex(); err != nil {
		t.Fatal(errReload)
	}

	tx2 := db.NewTx(sdb.RO)
	defer tx2.Discard()

	keys, errFind := tx2.Find("Data:lorem")
	if err != nil {
		t.Fatal(errFind)
	}

	if len(keys) != 1 {
		t.Error("Search index not reloaded. Want: 1 result, Got: ", len(keys))
	}
}
