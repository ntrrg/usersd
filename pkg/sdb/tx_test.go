// Copyright 2019 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the MIT license.

package sdb_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/dgraph-io/badger/v2"
	"nt.web.ve/go/usersd/pkg/sdb"
)

func TestTx_bigTransaction(t *testing.T) {
	db, err := sdb.Open(sdb.InMemory)
	if err != nil {
		t.Fatal(err)
	}

	tx := db.NewTx(sdb.RW)
	defer tx.Discard()

	data := "lorem ipsum"

	for i := 0; ; i++ {
		if err = tx.Set([]byte(fmt.Sprintf("test-%v", i)), &data); err != nil {
			break
		}
	}

	err = tx.Set([]byte("test"), &data)
	if !errors.Is(err, badger.ErrTxnTooBig) {
		t.Error("Big transaction allowed")
	}
}

func TestTx_Set_nonPointer(t *testing.T) {
	db, err := sdb.Open(sdb.InMemory)
	if err != nil {
		t.Fatal(err)
	}

	tx := db.NewTx(sdb.RW)
	defer tx.Discard()

	err = tx.Set([]byte("test"), "non pointer")
	if err != sdb.ErrValMustBePointer {
		t.Error("can't encode data if its is not a pointer")
	}
}

func BenchmarkTx_Set(b *testing.B) {
	db, err := sdb.Open(sdb.InMemory)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	for i := 0; i <= b.N; i++ {
		tx := db.NewTx(sdb.RW)

		if err := tx.Set([]byte("test"), &p); err != nil {
			b.Fatal(err)
		}

		tx.Discard()
	}
}

func BenchmarkTx_Get(b *testing.B) {
	db, err := sdb.Open(sdb.InMemory)
	if err != nil {
		b.Fatal(err)
	}

	tx := db.NewTx(sdb.RW)

	if err := tx.Set([]byte("test"), &p); err != nil {
		b.Fatal(err)
	}

	if err := tx.Commit(); err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	for i := 0; i <= b.N; i++ {
		var p2 Person

		tx := db.NewTx(sdb.RO)

		if err := tx.Get([]byte("test"), &p2); err != nil {
			b.Fatal(err)
		}

		tx.Discard()
	}
}
