// Copyright 2018 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the WTFPL.

package usersd_test

import (
	"testing"

	"github.com/ntrrg/usersd/pkg/usersd"
)

func TestService_CheckPassword(t *testing.T) {
	ud, err := usersd.New(usersd.DefaultOptions)
	if err != nil {
		t.Fatal(err)
	}

	defer ud.Close()

	tx := ud.NewTx(true)
	defer tx.Discard()

	user := new(usersd.User)
	if err := user.Write(tx); err != nil {
		t.Fatal(err)
	}

	if err := ud.SetPassword(tx, user.ID, "1234"); err != nil {
		t.Errorf("Can't assign the password -> %v", err)
	}

	if !ud.CheckPassword(tx, user.ID, "1234") {
		t.Error("Can't verify the password")
	}

	if ud.CheckPassword(tx, user.ID, "1235") {
		t.Error("Wrong password pass")
	}
}

func TestService_CheckPassword_emptyPassword(t *testing.T) {
	ud, err := usersd.New(usersd.DefaultOptions)
	if err != nil {
		t.Fatal(err)
	}

	defer ud.Close()

	tx := ud.NewTx(false)
	defer tx.Discard()

	if ud.CheckPassword(tx, "", "") {
		t.Error("Empty password pass")
	}
}

func TestService_CheckPassword_nonExistentUser(t *testing.T) {
	ud, err := usersd.New(usersd.DefaultOptions)
	if err != nil {
		t.Fatal(err)
	}

	defer ud.Close()

	tx := ud.NewTx(false)
	defer tx.Discard()

	if ud.CheckPassword(tx, "", "123") {
		t.Error("Non existent user pass")
	}
}

func TestService_CheckPassword_userWithoutPassword(t *testing.T) {
	ud, err := usersd.New(usersd.DefaultOptions)
	if err != nil {
		t.Fatal(err)
	}

	defer ud.Close()

	tx := ud.NewTx(true)
	defer tx.Discard()

	user := new(usersd.User)
	if err := user.Write(tx); err != nil {
		t.Fatal(err)
	}

	if ud.CheckPassword(tx, user.ID, "123") {
		t.Error("Non existent user pass")
	}
}

func TestService_SetPassword(t *testing.T) {
	ud, err := usersd.New(usersd.DefaultOptions)
	if err != nil {
		t.Fatal(err)
	}

	defer ud.Close()

	tx := ud.NewTx(true)
	defer tx.Discard()

	user := new(usersd.User)
	if err := user.Write(tx); err != nil {
		t.Fatal(err)
	}

	if err := ud.SetPassword(tx, user.ID, "1234"); err != nil {
		t.Errorf("Can't assign the password -> %v", err)
	}
}

func TestService_SetPassword_emptyPassword(t *testing.T) {
	ud, err := usersd.New(usersd.DefaultOptions)
	if err != nil {
		t.Fatal(err)
	}

	defer ud.Close()

	tx := ud.NewTx(false)
	defer tx.Discard()

	if err := ud.SetPassword(tx, "", ""); err == nil {
		t.Error("Empty password assigned")
	}
}

func TestService_SetPassword_nonExistentUser(t *testing.T) {
	ud, err := usersd.New(usersd.DefaultOptions)
	if err != nil {
		t.Fatal(err)
	}

	defer ud.Close()

	tx := ud.NewTx(false)
	defer tx.Discard()

	if err := ud.SetPassword(tx, "", "123"); err == nil {
		t.Error("Password assigned to non existent user")
	}
}
