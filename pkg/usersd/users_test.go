// Copyright 2018 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the WTFPL.

package usersd_test

import (
	"testing"

	"github.com/blevesearch/bleve"
	"github.com/dgraph-io/badger"
	"github.com/gofrs/uuid"
	"github.com/ntrrg/usersd/pkg/usersd"
	"golang.org/x/crypto/bcrypt"
)

var Opts = usersd.DefaultOptions

func TestGetUser(t *testing.T) {
	ud, err := usersd.New(Opts)
	if err != nil {
		t.Fatal(err)
	}

	defer ud.Close()

	tx := ud.DB.NewTransaction(true)
	defer tx.Discard()
	index := ud.Index["users"]

	usersFixtures(t, tx, index)

	user, err := usersd.GetUser(tx, "admin")
	if err != nil {
		t.Fatal(err)
	}

	if user.Email != "admin@example.com" {
		t.Errorf("GetUser(admin).Email == %v, wants admin@example.com", user.Email)
	}

	if user.Mode != "local" {
		t.Errorf("GetUser(admin).Mode == %v, wants local", user.Mode)
	}
}

func TestGetUsers(t *testing.T) {
	ud, err := usersd.New(Opts)
	if err != nil {
		t.Fatal(err)
	}

	defer ud.Close()

	tx := ud.DB.NewTransaction(true)
	defer tx.Discard()
	index := ud.Index["users"]

	usersFixtures(t, tx, index)

	cases := []struct {
		name string
		q    string
		sort []string
		want int
	}{
		{name: "All", want: 3},

		{
			name: "AllSorted",
			want: 3,
			sort: []string{"-email"},
		},

		{
			name: "ByEmail",
			want: 1,
			q:    `email:"john@example.com"`,
		},
	}

	for _, c := range cases {
		c := c

		t.Run(c.name, func(t *testing.T) {
			users, err := usersd.GetUsers(tx, index, c.q, c.sort...)
			if err != nil {
				t.Fatal(err)
			}

			if len(users) != c.want {
				t.Errorf("GetUsers(%v, %v) gets invalid data -> %v", c.q, c.sort, users)
			}
		})
	}
}

func TestNewUser(t *testing.T) { // nolint: gocyclo
	ud, err := usersd.New(Opts)
	if err != nil {
		t.Fatal(err)
	}

	defer ud.Close()

	tx := ud.DB.NewTransaction(true)
	defer tx.Discard()
	index := ud.Index["users"]

	cases := []struct {
		name string
		fail bool

		id, email, password string

		data map[string]interface{}
	}{
		{
			name:     "Regular",
			email:    "john@example.com",
			password: "1234",
		},

		{
			name:     "ExtraData",
			id:       "test",
			email:    "john2@example.com",
			password: "1234",
			data: map[string]interface{}{
				"username": "john",
				"name":     "John Doe",
			},
		},

		{
			name:     "EmptyEmail",
			fail:     true,
			password: "1234",
		},

		{
			name:  "EmptyPassword",
			fail:  true,
			email: "john@example.com",
		},

		{
			name:     "ExistentUser",
			fail:     true,
			email:    "john@example.com",
			password: "1234",
		},
	}

	for _, c := range cases {
		c := c

		t.Run(c.name, func(t *testing.T) {
			user, err := usersd.NewUser(tx, index, c.id, c.email, c.password, c.data)

			switch {
			case err != nil && !c.fail:
				t.Fatal(err)
			case err == nil && c.fail:
				t.Fatal("User created")
			case err != nil && c.fail:
				return
			}

			if c.id == "" {
				id, err := uuid.FromString(user.ID)
				if err != nil {
					t.Fatalf("Invalid UUID (%v) -> %s", user.ID, err)
				}

				if id.Version() != 4 {
					t.Errorf("Invalid UUID version (%v)", id.Version())
				}
			} else {
				if user.ID != c.id {
					t.Errorf("User ID = %v, want %v", user.ID, c.id)
				}
			}

			if _, err := bcrypt.Cost([]byte(user.Password)); err != nil {
				t.Errorf("Invalid password hash (%v) -> %s", user.Password, err)
			}
		})
	}
}

func TestUser_Delete(t *testing.T) {
	ud, err := usersd.New(Opts)
	if err != nil {
		t.Fatal(err)
	}

	defer ud.Close()

	tx := ud.DB.NewTransaction(true)
	defer tx.Discard()
	index := ud.Index["users"]

	usersFixtures(t, tx, index)

	users, err := usersd.GetUsers(tx, index, "")
	if err != nil {
		t.Fatal(err)
	}

	if len(users) == 0 {
		t.Fatal("No users created")
	}

	for _, user := range users {
		if err = user.Delete(tx, index); err != nil {
			t.Fatal(err)
		}
	}

	users, err = usersd.GetUsers(tx, index, "")
	if err != nil {
		t.Fatal(err)
	}

	if len(users) != 0 {
		t.Error("The database keeps users even after deleting all of them")
	}
}

func TestUser_Write(t *testing.T) {
	ud, err := usersd.New(Opts)
	if err != nil {
		t.Fatal(err)
	}

	defer ud.Close()

	tx := ud.DB.NewTransaction(true)
	defer tx.Discard()
	index := ud.Index["users"]

	usersFixtures(t, tx, index)

	user, err := usersd.GetUser(tx, "admin")
	if err != nil {
		t.Fatal(err)
	}

	user.Email = "admin@test.com"

	if err = user.Write(tx, index); err != nil {
		t.Errorf("Can't write to the user -> %v", err)
	}

	newUser, err := usersd.GetUser(tx, "admin")
	if err != nil {
		t.Fatal(err)
	}

	if newUser.Email != user.Email {
		t.Errorf("Update failed, got %v,  wants %v", newUser.Email, user.Email)
	}
}

func usersFixtures(t *testing.T, tx *badger.Txn, index bleve.Index) {
	users := []struct {
		id, email, password string

		data map[string]interface{}
	}{
		{
			id:       "admin",
			email:    "admin@example.com",
			password: "admin",
		},

		{
			email:    "john@example.com",
			password: "1234",
		},

		{
			email:    "john2@example.com",
			password: "1234",
			data: map[string]interface{}{
				"username": "john",
				"name":     "John Doe",
			},
		},
	}

	for _, u := range users {
		_, err := usersd.NewUser(tx, index, u.id, u.email, u.password, u.data)

		if err != nil {
			t.Fatal(err)
		}
	}
}
