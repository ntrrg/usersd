// Copyright 2018 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the WTFPL.

package usersd_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/ntrrg/usersd/pkg/usersd"
)

func TestBackup(t *testing.T) {
	if err := initTest("backup", true); err != nil {
		t.Fatal(err)
	}

	defer usersd.Close()

	if err := usersd.Backup(bytes.NewBuffer(nil)); err != nil {
		t.Error(err)
	}
}

type BadWriter struct{}

func (w BadWriter) Write(p []byte) (int, error) {
	return 0, errors.New("can't write into a bad writer")
}

func TestBackup_badWriter(t *testing.T) {
	if err := initTest("backup-bad-writer", true); err != nil {
		t.Fatal(err)
	}

	defer usersd.Close()

	if err := usersd.Backup(BadWriter{}); err == nil {
		t.Error("Backup writed into a bad writer")
	}
}

func TestRestore(t *testing.T) {
	if err := initTest("restore", true); err != nil {
		t.Fatal(err)
	}

	defer usersd.Close()

	buf := bytes.NewBuffer(nil)
	if err := usersd.Backup(buf); err != nil {
		t.Error(err)
	}

	tx := usersd.NewTx(true)

	users, err := tx.GetUsers("")
	if err != nil {
		t.Fatal(err)
	}

	n := len(users)

	if err = tx.DeleteUser(users[0].ID); err != nil {
		t.Fatal(err)
	}

	if err = tx.Commit(); err != nil {
		t.Fatal(err)
	}

	if err = usersd.Restore(buf); err != nil {
		t.Fatal(err)
	}

	tx = usersd.NewTx(false)
	defer tx.Discard()

	users, err = tx.GetUsers("")
	if err != nil {
		t.Fatal(err)
	}

	if n == len(users) {
		t.Error("Data wasn't restored")
	}
}

type BadReader struct{}

func (r BadReader) Read(p []byte) (int, error) {
	return 0, errors.New("can't read from a bad reader")
}

func TestRestore_badReader(t *testing.T) {
	if err := initTest("restore-bad-reader", false); err != nil {
		t.Fatal(err)
	}

	defer usersd.Close()

	if err := usersd.Restore(BadReader{}); err == nil {
		t.Error("Backup restored from a bad reader")
	}
}
