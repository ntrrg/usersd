// Copyright 2018 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the WTFPL.

package usersd_test

import (
	"testing"

	"github.com/ntrrg/usersd/pkg/usersd"
)

func TestService_Password(t *testing.T) {
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

	if err := ud.SetPasswordTx(tx, user, "1234"); err != nil {
		t.Errorf("Can't assign the password -> %v", err)
	}

	if !ud.CheckPasswordTx(tx, user, "1234") {
		t.Errorf("Can't verify the password")
	}
}
