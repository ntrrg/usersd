// Copyright 2018 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the WTFPL.

package usersd_test

import (
	"testing"

	"github.com/ntrrg/usersd/pkg/usersd"
)

func TestTx_GetUser(t *testing.T) {
	if err := initTest("tx-get-user", true); err != nil {
		t.Fatal(err)
	}

	defer usersd.Close()

	tx := usersd.NewTx(true)
	defer tx.Discard()

	if _, err := tx.GetUser("admin"); err != nil {
		t.Errorf("Can't fetch the user data -> %v", err)
	}
}

func TestTx_GetUser_discartedTx(t *testing.T) {
	if err := initTest("tx-get-user-discarted", false); err != nil {
		t.Fatal(err)
	}

	defer usersd.Close()

	tx := usersd.NewTx(false)
	tx.Discard()

	if _, err := tx.GetUser("admin"); err == nil {
		t.Error("Getting user with discarted transaction")
	}
}

func TestTx_GetUser_malformedData(t *testing.T) {
	if err := initTest("tx-get-user-malformed", false); err != nil {
		t.Fatal(err)
	}

	defer usersd.Close()

	tx := usersd.NewTx(true)
	defer tx.Discard()

	if err := tx.Set([]byte(usersd.UsersDI+"admin"), []byte{1, 2}); err != nil {
		t.Fatal(err)
	}

	if _, err := tx.GetUser("admin"); err == nil {
		t.Error("Getting user with malformed data")
	}
}

func TestTx_GetUsers(t *testing.T) {
	if err := initTest("tx-get-users", true); err != nil {
		t.Fatal(err)
	}

	defer usersd.Close()

	tx := usersd.NewTx(true)
	defer tx.Discard()

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
			users, err := tx.GetUsers(c.q)
			if err != nil {
				t.Fatal(err)
			}

			if len(users) != c.want {
				t.Errorf("tx.GetUsers(%v) gets invalid data -> %v", c.q, users)
			}
		})
	}
}

func TestTx_GetUsers_sorted(t *testing.T) {
	if err := initTest("tx-get-users-sorted", true); err != nil {
		t.Fatal(err)
	}

	defer usersd.Close()

	tx := usersd.NewTx(true)
	defer tx.Discard()

	users, err := tx.GetUsers("", "-email")
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

func TestTx_GetUsers_outdatedIndex(t *testing.T) {
	if err := initTest("tx-get-users-outdated-index", true); err != nil {
		t.Fatal(err)
	}

	defer usersd.Close()

	tx := usersd.NewTx(true)
	defer tx.Discard()

	if err := tx.Delete([]byte(usersd.UsersDI + "admin")); err != nil {
		t.Fatal(err)
	}

	if _, err := tx.GetUsers("+id:admin"); err == nil {
		t.Error("Getting users with an outdated search index")
	}
}

func TestTx_GetUsers_malformedData(t *testing.T) {
	if err := initTest("tx-get-users-malformed", true); err != nil {
		t.Fatal(err)
	}

	defer usersd.Close()

	tx := usersd.NewTx(true)
	defer tx.Discard()

	if err := tx.Set([]byte(usersd.UsersDI+"admin"), []byte{1, 2}); err != nil {
		t.Fatal(err)
	}

	if _, err := tx.GetUsers(""); err == nil {
		t.Error("Getting users with malformed data")
	}
}

func TestTx_DeleteUser(t *testing.T) {
	if err := initTest("tx-delete-user", true); err != nil {
		t.Fatal(err)
	}

	defer usersd.Close()

	tx := usersd.NewTx(true)
	defer tx.Discard()

	users, err := tx.GetUsers("")
	if err != nil {
		t.Fatal(err)
	}

	if len(users) == 0 {
		t.Fatal("No users created")
	}

	for _, user := range users {
		if err = tx.DeleteUser(user.ID); err != nil {
			t.Fatal(err)
		}
	}

	users, err = tx.GetUsers("")
	if err != nil {
		t.Fatal(err)
	}

	if len(users) != 0 {
		t.Error("The database keeps users even after deleting all of them")
	}
}

func TestTx_DeleteUser_discartedTx(t *testing.T) {
	if err := initTest("tx-delete-user-discarted", false); err != nil {
		t.Fatal(err)
	}

	defer usersd.Close()

	tx := usersd.NewTx(false)
	tx.Discard()

	if err := tx.DeleteUser(""); err == nil {
		t.Error("Removing user with discarted transaction")
	}
}

func TestTx_DeleteUser_roTx(t *testing.T) {
	if err := initTest("tx-delete-user-ro", false); err != nil {
		t.Fatal(err)
	}

	defer usersd.Close()

	tx := usersd.NewTx(false)
	defer tx.Discard()

	if err := tx.DeleteUser(""); err == nil {
		t.Error("Removing user with read-only transaction")
	}
}

func TestTx_WriteUser(t *testing.T) {
	if err := initTest("tx-write-user", false); err != nil {
		t.Fatal(err)
	}

	defer usersd.Close()

	tx := usersd.NewTx(true)
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
			err := tx.WriteUser(c.user)

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

func TestTx_WriteUser_discartedTx(t *testing.T) {
	if err := initTest("tx-write-user-discarted", false); err != nil {
		t.Fatal(err)
	}

	defer usersd.Close()

	tx := usersd.NewTx(true)
	tx.Discard()

	user := new(usersd.User)
	if err := tx.WriteUser(user); err == nil {
		t.Error("Writing user with discarted transaction")
	}
}

func TestTx_WriteUser_roTx(t *testing.T) {
	if err := initTest("tx-write-user-ro", false); err != nil {
		t.Fatal(err)
	}

	defer usersd.Close()

	tx := usersd.NewTx(false)
	defer tx.Discard()

	user := new(usersd.User)
	if err := tx.WriteUser(user); err == nil {
		t.Error("Writing user with read-only transaction")
	}
}
