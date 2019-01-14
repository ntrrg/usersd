// Copyright 2018 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the WTFPL.

package usersd_test

import (
	"testing"
)

func TestTx_CheckPassword(t *testing.T) {
	ud, err := initTest("tx-check-password", true)
	if err != nil {
		t.Fatal(err)
	}

	defer ud.Close()

	tx := ud.NewTx(true)
	defer tx.Discard()

	if !tx.CheckPassword("admin", "admin") {
		t.Error("Can't verify the password")
	}

	if err := tx.SetPassword("admin", "1234"); err != nil {
		t.Errorf("Can't assign the password -> %v", err)
	}

	if tx.CheckPassword("admin", "admin") {
		t.Error("Wrong password pass")
	}
}

func TestTx_CheckPassword_emptyPassword(t *testing.T) {
	ud, err := initTest("tx-check-password-empty-password", false)
	if err != nil {
		t.Fatal(err)
	}

	defer ud.Close()

	tx := ud.NewTx(false)
	defer tx.Discard()

	if tx.CheckPassword("", "") {
		t.Error("Empty password pass")
	}
}

func TestTx_CheckPassword_nonExistentUser(t *testing.T) {
	ud, err := initTest("tx-check-password-ne-user", false)
	if err != nil {
		t.Fatal(err)
	}

	defer ud.Close()

	tx := ud.NewTx(false)
	defer tx.Discard()

	if tx.CheckPassword("", "1234") {
		t.Error("Non existent user pass")
	}
}

func TestTx_CheckPassword_userWithoutPassword(t *testing.T) {
	ud, err := initTest("tx-check-password-user-wo-password", true)
	if err != nil {
		t.Fatal(err)
	}

	defer ud.Close()

	tx := ud.NewTx(true)
	defer tx.Discard()

	users, err := tx.GetUsers(`+email:"john@example.com"`)
	if err != nil {
		t.Fatal(err)
	} else if len(users) == 0 {
		t.Fatal("Can't find the given user")
	}

	if tx.CheckPassword(users[0].ID, "1234") {
		t.Error("User without password pass")
	}
}

func TestTx_SetPassword(t *testing.T) {
	ud, err := initTest("tx-set-password", true)
	if err != nil {
		t.Fatal(err)
	}

	defer ud.Close()

	tx := ud.NewTx(true)
	defer tx.Discard()

	if err := tx.SetPassword("admin", "1234"); err != nil {
		t.Errorf("Can't assign the password -> %v", err)
	}
}

func TestService_SetPassword_emptyPassword(t *testing.T) {
	ud, err := initTest("tx-set-password-empty-password", false)
	if err != nil {
		t.Fatal(err)
	}

	defer ud.Close()

	tx := ud.NewTx(false)
	defer tx.Discard()

	if err := tx.SetPassword("", ""); err == nil {
		t.Error("Empty password assigned")
	}
}

func TestService_SetPassword_nonExistentUser(t *testing.T) {
	ud, err := initTest("tx-set-password-ne-user", false)
	if err != nil {
		t.Fatal(err)
	}

	defer ud.Close()

	tx := ud.NewTx(false)
	defer tx.Discard()

	if err := tx.SetPassword("", "1234"); err == nil {
		t.Error("Password assigned to non existent user")
	}
}
