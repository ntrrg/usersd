// Copyright 2018 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the WTFPL.

package usersd_test

import (
	"fmt"
	"testing"

	"github.com/ntrrg/usersd/pkg/usersd"
)

func TestError_Error(t *testing.T) {
	err := usersd.ErrUserPasswordEmpty
	want := fmt.Sprintf("(%d) %s: %s", err.Code, err.Field, err.Message)
	got := err.Error()

	if got != want {
		t.Errorf("Bad error. Want: %v; got: %v", want, got)
	}
}

func TestError_Format(t *testing.T) {
	err := usersd.ErrUserPasswordHash
	extra := "Invalid input"

	want := fmt.Sprintf(
		"(%d) %s: %s",
		err.Code, err.Field, fmt.Sprintf(err.Message, extra),
	)

	got := err.Format(extra).Error()

	if got != want {
		t.Errorf("Bad error formatting. Want: %v; got: %v", want, got)
	}
}

func TestErrors_Error(t *testing.T) {
	err := usersd.Errors{
		usersd.ErrUserIDNotFound,
		usersd.ErrUserEmailEmpty,
		usersd.ErrUserPasswordEmpty,
	}

	want := "(1) id: the given user ID doesn't exists; "
	want += "(10) email: the given email is empty; "
	want += "(30) password: the given password is empty"

	got := err.Error()

	if got != want {
		t.Errorf("Bad errors formating. Want: %v; got: %v", want, got)
	}
}

func TestErrors_Error_empty(t *testing.T) {
	err := usersd.Errors{}
	got := err.Error()

	if got != "" {
		t.Errorf("Bad errors formating. Want an empty string; got: %v", got)
	}
}
