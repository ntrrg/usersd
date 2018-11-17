// Copyright 2018 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the WTFPL.

package usersd_test

import (
	"log"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/ntrrg/usersd/pkg/usersd"
)

type userData struct {
	id, password string

	data map[string]interface{}
}

type userCase struct {
	in, want userData
}

func TestGetUser(t *testing.T) {
	if err := usersd.Init(Opts); err != nil {
		t.Fatal(err)
	}

	defer usersd.Close()
	usersFixtures()

	cases := []struct {
		in, want string
	}{
		{"admin", "Administrator"},
	}

	for i, c := range cases {
		user, err := usersd.GetUser(c.in)

		if err != nil {
			t.Errorf("TC#%v: GetUser(%v) error -> %v", i, c.in, err)
			continue
		}

		name := user.Data["name"]

		if name != c.want {
			msg := "TC#%v: GetUser(%v).Data[name] == %v, wants %v"
			t.Errorf(msg, i, c.in, name, c.want)
		}
	}
}

func TestListUsers(t *testing.T) {
	if err := usersd.Init(Opts); err != nil {
		t.Fatal(err)
	}

	defer usersd.Close()
	usersFixtures()

	users, err := usersd.ListUsers()

	if err != nil {
		t.Fatal(err)
	}

	if len(users) < 1 {
		t.Error("ListUsers() doesn't fetch any data.")
	}
}

func TestUser_Save(t *testing.T) {
	if err := usersd.Init(Opts); err != nil {
		t.Fatal(err)
	}

	defer usersd.Close()

	cases := []userCase{
		{
			userData{id: "admin"},
			userData{id: "admin"},
		},

		{
			userData{
				data: map[string]interface{}{
					"id": "ntrrg",
				},
			},
			userData{
				data: map[string]interface{}{
					"id": "ntrrg",
				},
			},
		},
	}

	for i, c := range cases {
		user, err := usersd.NewUser(c.in.id, c.in.password, c.in.data)

		if err != nil {
			t.Errorf("TC#%v: %s", i, err)
		}

		if err := user.Save(); err != nil {
			t.Errorf("TC#%v: NewUser(%+v).Save() error -> %v", i, c.in, err)
		}

		if c.want.id == "" {
			id, err := uuid.FromString(user.ID)

			if err != nil {
				t.Errorf(
					"TC#%v: NewUser(%+v).ID invalid UUID (%v) -> %s",
					i, c.in, user.ID, err,
				)
			}

			if id.Version() != 4 {
				t.Errorf(
					"TC#%v: NewUser(%+v).ID invalid UUID version (%v)",
					i, c.in, id.Version(),
				)
			}
		} else {
			if user.ID != c.want.id {
				t.Errorf(
					"TC#%v: NewUser(%v).ID == %+v, want %v",
					i, c.in, user.ID, c.want.id,
				)
			}
		}
	}
}

func TestUser_Delete(t *testing.T) {
	if err := usersd.Init(Opts); err != nil {
		t.Fatal(err)
	}

	defer usersd.Close()
	usersFixtures()

	users, err := usersd.ListUsers()

	if err != nil {
		t.Fatal(err)
	}

	x := len(users)

	for _, user := range users {
		if err := user.Delete(); err != nil {
			t.Errorf("User(%+v).Delete() error -> %v", user, err)
		}
	}

	users, err = usersd.ListUsers()

	if err != nil {
		t.Fatal(err)
	}

	if len(users) >= x {
		msg := "The users list keeps data even after calling User.Delete()"
		t.Error(msg)
	}
}

func usersFixtures() {
	users := []userData{
		{password: "1234"},

		{
			data: map[string]interface{}{
				"username": "ntrrg",
				"name":     "Miguel Angel Rivera Notararigo",
			},
		},
	}

	for _, u := range users {
		user, err := usersd.NewUser(u.id, u.password, u.data)

		if err != nil {
			log.Fatal(err)
		}

		if err := user.Save(); err != nil {
			log.Fatal(err)
		}
	}
}
