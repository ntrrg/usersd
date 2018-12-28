// Copyright 2018 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the WTFPL.

package usersd_test

import (
	"fmt"
	"testing"

	"github.com/ntrrg/usersd/pkg/usersd"
)

func TestValidationError_Error(t *testing.T) {
	err := usersd.ErrUserPasswordEmpty
	want := fmt.Sprintf("(%s) %s", err.Field, err.Message)
	got := err.Error()

	if got != want {
		t.Errorf("Bad error. Want: %v; got: %v", want, got)
	}
}

func TestValidationError_Format(t *testing.T) {
	err := usersd.ErrUserPasswordHash
	extra := "Invalid input"
	want := fmt.Sprintf("(%s) %s", err.Field, fmt.Sprintf(err.Message, extra))
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

	want := "The given user ID doesn't exists; "
	want += "(email) The given email is empty; "
	want += "(password) The given password is empty"

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
