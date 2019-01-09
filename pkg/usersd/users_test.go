// Copyright 2018 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the WTFPL.

package usersd_test

import (
	"testing"

	"github.com/ntrrg/usersd/pkg/usersd"
)

func TestGetUser(t *testing.T) {
	ud, err := usersd.New(usersd.DefaultOptions)
	if err != nil {
		t.Fatal(err)
	}

	defer ud.Close()

	tx := ud.NewTx(true)
	defer tx.Discard()

	usersFixtures(t, tx)

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

func TestGetUser_discartedTx(t *testing.T) {
	ud, err := usersd.New(usersd.DefaultOptions)
	if err != nil {
		t.Fatal(err)
	}

	defer ud.Close()

	tx := ud.NewTx(false)
	tx.Discard()

	if _, err = usersd.GetUser(tx, "admin"); err == nil {
		t.Fatal("Getting user with discarted transaction")
	}
}

func TestGetUser_malformedData(t *testing.T) {
	ud, err := usersd.New(usersd.DefaultOptions)
	if err != nil {
		t.Fatal(err)
	}

	defer ud.Close()

	tx := ud.NewTx(true)
	defer tx.Discard()

	if err = tx.Set([]byte("usersadmin"), []byte{1, 2}); err != nil {
		t.Fatal(err)
	}

	if _, err = usersd.GetUser(tx, "admin"); err == nil {
		t.Fatal("Getting user with malformed data")
	}
}

func TestGetUsers(t *testing.T) {
	ud, err := usersd.New(usersd.DefaultOptions)
	if err != nil {
		t.Fatal(err)
	}

	defer ud.Close()

	tx := ud.NewTx(true)
	defer tx.Discard()

	usersFixtures(t, tx)

	cases := []struct {
		name string
		q    string
		want int
	}{
		{name: "All", want: 3},

		{
			name: "ByEmail",
			want: 1,
			q:    `+email:john@example.com`,
		},

		{
			name: "ByExtraData",
			want: 1,
			q:    `+data.name:"John Doe"`,
		},
	}

	for _, c := range cases {
		c := c

		t.Run(c.name, func(t *testing.T) {
			users, err := usersd.GetUsers(tx, c.q)
			if err != nil {
				t.Fatal(err)
			}

			if len(users) != c.want {
				t.Errorf("GetUsers(%v) gets invalid data -> %v", c.q, users)
			}
		})
	}
}

func TestGetUsers_sorted(t *testing.T) {
	ud, err := usersd.New(usersd.DefaultOptions)
	if err != nil {
		t.Fatal(err)
	}

	defer ud.Close()

	tx := ud.NewTx(true)
	defer tx.Discard()

	usersFixtures(t, tx)

	users, err := usersd.GetUsers(tx, "", "-email")
	if err != nil {
		t.Fatal("Can't fetch the users")
	}

	emails := []string{
		"john@example.com",
		"john2@example.com",
		"admin@example.com",
	}

	for i, email := range emails {
		if email != users[i].Email {
			t.Errorf("Bad order: %q - %q", email, users[i].Email)
		}
	}
}

func TestGetUsers_outdatedIndex(t *testing.T) {
	ud, err := usersd.New(usersd.DefaultOptions)
	if err != nil {
		t.Fatal(err)
	}

	defer ud.Close()

	tx := ud.NewTx(true)
	defer tx.Discard()

	usersFixtures(t, tx)

	if err = tx.Delete([]byte("usersadmin")); err != nil {
		t.Fatal(err)
	}

	if _, err = usersd.GetUsers(tx, "+id:admin"); err == nil {
		t.Error("Getting users with an outdated search index")
	}
}

func TestGetUsers_malformedData(t *testing.T) {
	ud, err := usersd.New(usersd.DefaultOptions)
	if err != nil {
		t.Fatal(err)
	}

	defer ud.Close()

	tx := ud.NewTx(true)
	defer tx.Discard()

	if err = tx.Set([]byte("usersadmin"), []byte{1, 2}); err != nil {
		t.Fatal(err)
	}

	if _, err = usersd.GetUsers(tx, ""); err == nil {
		t.Error("Getting users with malformed data")
	}
}

func TestUser_Delete(t *testing.T) {
	ud, err := usersd.New(usersd.DefaultOptions)
	if err != nil {
		t.Fatal(err)
	}

	defer ud.Close()

	tx := ud.NewTx(true)
	defer tx.Discard()

	usersFixtures(t, tx)

	users, err := usersd.GetUsers(tx, "")
	if err != nil {
		t.Fatal(err)
	}

	if len(users) == 0 {
		t.Fatal("No users created")
	}

	for _, user := range users {
		if err = user.Delete(tx); err != nil {
			t.Fatal(err)
		}
	}

	users, err = usersd.GetUsers(tx, "")
	if err != nil {
		t.Fatal(err)
	}

	if len(users) != 0 {
		t.Error("The database keeps users even after deleting all of them")
	}
}

func TestUser_Delete_discartedTx(t *testing.T) {
	ud, err := usersd.New(usersd.DefaultOptions)
	if err != nil {
		t.Fatal(err)
	}

	defer ud.Close()

	tx := ud.NewTx(true)
	tx.Discard()

	user := &usersd.User{ID: "test"}

	if err = user.Delete(tx); err == nil {
		t.Error("Removing user with discarted transaction")
	}
}

func TestUser_Delete_roTx(t *testing.T) {
	ud, err := usersd.New(usersd.DefaultOptions)
	if err != nil {
		t.Fatal(err)
	}

	defer ud.Close()

	tx := ud.NewTx(false)
	defer tx.Discard()

	user := &usersd.User{ID: "test"}

	if err = user.Delete(tx); err == nil {
		t.Error("Removing user with read-only transaction")
	}
}

func TestUser_Write(t *testing.T) {
	ud, err := usersd.New(usersd.DefaultOptions)
	if err != nil {
		t.Fatal(err)
	}

	defer ud.Close()

	tx := ud.NewTx(true)
	defer tx.Discard()

	cases := []struct {
		name string
		fail bool
		user *usersd.User
	}{
		{
			name: "Simple",
			user: new(usersd.User),
		},

		{
			name: "ExtraData",
			user: &usersd.User{
				ID:    "test",
				Email: "john@example.com",
				Phone: "+12345678901",
				Data: map[string]interface{}{
					"username": "john",
					"name":     "John Doe",
				},
			},
		},

		{
			name: "InvalidEmail",
			fail: true,
			user: &usersd.User{Email: "johnexample.com"},
		},

		{
			name: "ExistentEmail",
			fail: true,
			user: &usersd.User{Email: "john@example.com"},
		},

		{
			name: "ExistentPhone",
			fail: true,
			user: &usersd.User{Phone: "+12345678901"},
		},

		{
			name: "Update",
			user: &usersd.User{
				ID:    "test",
				Email: "john@example.com",
				Phone: "+12345678901",
				Data: map[string]interface{}{
					"username": "john",
					"name":     "John Doe",
					"age":      26,
				},
			},
		},
	}

	for _, c := range cases {
		c := c

		t.Run(c.name, func(t *testing.T) {
			err := c.user.Write(tx)

			switch {
			case err != nil && !c.fail:
				t.Fatal(err)
			case err == nil && c.fail:
				t.Fatal("User created")
			case err != nil && c.fail:
				return
			}
		})
	}
}

func TestUser_Write_discartedTx(t *testing.T) {
	ud, err := usersd.New(usersd.DefaultOptions)
	if err != nil {
		t.Fatal(err)
	}

	defer ud.Close()

	tx := ud.NewTx(true)
	tx.Discard()

	user := &usersd.User{}

	if err = user.Write(tx); err == nil {
		t.Error("Writing user with discarted transaction")
	}
}

func TestUser_Write_roTx(t *testing.T) {
	ud, err := usersd.New(usersd.DefaultOptions)
	if err != nil {
		t.Fatal(err)
	}

	defer ud.Close()

	tx := ud.NewTx(false)
	defer tx.Discard()

	user := new(usersd.User)
	if err = user.Write(tx); err == nil {
		t.Error("Writing user with read-only transaction")
	}
}

func TestService_GetUser(t *testing.T) {
	ud, err := usersd.New(usersd.DefaultOptions)
	if err != nil {
		t.Fatal(err)
	}

	defer ud.Close()

	tx := ud.NewTx(true)
	defer tx.Discard()

	usersFixtures(t, tx)

	if err = tx.Commit(); err != nil {
		t.Fatal(err)
	}

	if _, err := ud.GetUser("admin"); err != nil {
		t.Fatal(err)
	}
}

func TestService_GetUsers(t *testing.T) {
	ud, err := usersd.New(usersd.DefaultOptions)
	if err != nil {
		t.Fatal(err)
	}

	defer ud.Close()

	tx := ud.NewTx(true)
	defer tx.Discard()

	usersFixtures(t, tx)

	if err = tx.Commit(); err != nil {
		t.Fatal(err)
	}

	users, err := ud.GetUsers("")
	if err != nil {
		t.Fatal(err)
	}

	if len(users) == 0 {
		t.Error("0 usersd fetched")
	}
}

func TestService_DeleteUser(t *testing.T) {
	ud, err := usersd.New(usersd.DefaultOptions)
	if err != nil {
		t.Fatal(err)
	}

	defer ud.Close()

	tx := ud.NewTx(true)
	defer tx.Discard()

	usersFixtures(t, tx)

	if err = tx.Commit(); err != nil {
		t.Fatal(err)
	}

	if err = ud.DeleteUser("admin"); err != nil {
		t.Fatal(err)
	}
}

func TestService_WriteUser(t *testing.T) {
	ud, err := usersd.New(usersd.DefaultOptions)
	if err != nil {
		t.Fatal(err)
	}

	defer ud.Close()

	user := new(usersd.User)
	if err = ud.WriteUser(user); err != nil {
		t.Fatal(err)
	}
}

func usersFixtures(t *testing.T, tx *usersd.Tx) {
	users := []*usersd.User{
		{
			ID:    "admin",
			Email: "admin@example.com",
		},

		{
			Email: "john@example.com",
			Phone: "+12345678901",
		},

		{
			Email: "john2@example.com",
			Data: map[string]interface{}{
				"username": "john",
				"name":     "John Doe",
			},
		},
	}

	for _, user := range users {
		if err := user.Write(tx); err != nil {
			t.Fatal(err)
		}
	}
}
