// Copyright 2018 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the WTFPL.

package usersd_test

import (
	"testing"

	"github.com/ntrrg/usersd/pkg/usersd"
)

func TestNewUser(t *testing.T) {
	if err := usersd.Init(""); err != nil {
		t.Fatal(err)
	}

	defer usersd.Close()

	cases := []struct {
		in, want string
	}{
		{"Miguel", "Miguel"},
		{"Angel", "Angel"},
		{"Rivera", "Rivera"},
		{"Notararigo", "Notararigo"},
	}

	for i, c := range cases {
		user, err := usersd.NewUser(c.want)

		if err != nil {
			t.Errorf("TC#%v: %s", i, err)
		}

		if user.Name != c.want {
			t.Errorf("TC#%v: NewUser(%v) == %+v, want %v", i, c.in, user, c.want)
		}
	}
}
