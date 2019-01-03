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
			users, err := usersd.GetUsers(tx, c.q, c.sort...)
			if err != nil {
				t.Fatal(err)
			}

			if len(users) != c.want {
				t.Errorf("GetUsers(%v, %v) gets invalid data -> %v", c.q, c.sort, users)
			}
		})
	}
}

func TestUser_CheckPassword(t *testing.T) {
	ud, err := usersd.New(usersd.DefaultOptions)
	if err != nil {
		t.Fatal(err)
	}

	defer ud.Close()

	tx := ud.NewTx(false)
	defer tx.Discard()

	user := &usersd.User{
		ID:       "test",
		Email:    "test@example.com",
		Password: "1234",
	}

	if err := user.Validate(tx); err != nil {
		t.Error(err)
	}

	if !user.CheckPassword("1234") {
		t.Error("Invalid password")
	}

	user.Password = ""

	if user.CheckPassword("1234") {
		t.Error("Empty password pass the check")
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
			name: "Regular",
			user: &usersd.User{
				Password: "1234",
			},
		},

		{
			name: "ExtraData",
			user: &usersd.User{
				ID:       "test",
				Email:    "john@example.com",
				Phone:    "+12345678901",
				Password: "1234",
				Data: map[string]interface{}{
					"username": "john",
					"name":     "John Doe",
				},
			},
		},

		{
			name: "EmptyPassword",
			fail: true,
			user: &usersd.User{},
		},

		{
			name: "OAuth2",
			user: &usersd.User{
				Mode:  "oauth2",
				Email: "test@gmail.com",
			},
		},

		{
			name: "InvalidEmail",
			fail: true,
			user: &usersd.User{
				Email:    "johnexample.com",
				Password: "1234",
			},
		},

		{
			name: "ExistentEmail",
			fail: true,
			user: &usersd.User{
				Email:    "john@example.com",
				Password: "1234",
			},
		},

		{
			name: "ExistentPhone",
			fail: true,
			user: &usersd.User{
				Phone:    "+12345678901",
				Password: "1234",
			},
		},

		{
			name: "Update",
			user: &usersd.User{
				ID:       "test",
				Email:    "john@example.com",
				Phone:    "+12345678901",
				Password: "1234",
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

func TestService_GetUser(t *testing.T) {
	ud, err := usersd.New(usersd.DefaultOptions)
	if err != nil {
		t.Fatal(err)
	}

	defer ud.Close()

	tx := ud.NewTx(true)
	defer tx.Discard()

	usersFixtures(t, tx)
	tx.Commit()

	user, err := ud.GetUser("admin")
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

func TestService_GetUsers(t *testing.T) {
	ud, err := usersd.New(usersd.DefaultOptions)
	if err != nil {
		t.Fatal(err)
	}

	defer ud.Close()

	tx := ud.NewTx(true)
	defer tx.Discard()

	usersFixtures(t, tx)
	tx.Commit()

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
			users, err := ud.GetUsers(c.q, c.sort...)
			if err != nil {
				t.Fatal(err)
			}

			if len(users) != c.want {
				t.Errorf("GetUsers(%v, %v) gets invalid data -> %v", c.q, c.sort, users)
			}
		})
	}
}

func TestService_DeleteUser(t *testing.T) {
	ud, err := usersd.New(usersd.DefaultOptions)
	if err != nil {
		t.Fatal(err)
	}

	tx := ud.NewTx(true)
	defer tx.Discard()

	usersFixtures(t, tx)
	tx.Commit()

	users, err := ud.GetUsers("")
	if err != nil {
		ud.Close()
		t.Fatal(err)
	}

	if len(users) == 0 {
		ud.Close()
		t.Fatal("No users created")
	}

	for _, user := range users {
		if err = ud.DeleteUser(user.ID); err != nil {
			ud.Close()
			t.Fatal(err)
		}
	}

	users, err = ud.GetUsers("")
	if err != nil {
		ud.Close()
		t.Fatal(err)
	}

	if len(users) != 0 {
		ud.Close()
		t.Error("The database keeps users even after deleting all of them")
	}

	ud.Close()

	if err := ud.DeleteUser("invalid ID"); err == nil {
		t.Error("Deletion works without database")
	}
}

func TestService_WriteUser(t *testing.T) {
	ud, err := usersd.New(usersd.DefaultOptions)
	if err != nil {
		t.Fatal(err)
	}

	defer ud.Close()

	cases := []struct {
		name string
		fail bool
		user *usersd.User
	}{
		{
			name: "Regular",
			user: &usersd.User{
				Password: "1234",
			},
		},

		{
			name: "ExtraData",
			user: &usersd.User{
				ID:       "test",
				Email:    "john@example.com",
				Phone:    "+12345678901",
				Password: "1234",
				Data: map[string]interface{}{
					"username": "john",
					"name":     "John Doe",
				},
			},
		},

		{
			name: "EmptyPassword",
			fail: true,
			user: &usersd.User{},
		},

		{
			name: "OAuth2",
			user: &usersd.User{
				Mode:  "oauth2",
				Email: "test@gmail.com",
			},
		},

		{
			name: "InvalidEmail",
			fail: true,
			user: &usersd.User{
				Email:    "johnexample.com",
				Password: "1234",
			},
		},

		{
			name: "ExistentEmail",
			fail: true,
			user: &usersd.User{
				Email:    "john@example.com",
				Password: "1234",
			},
		},

		{
			name: "ExistentPhone",
			fail: true,
			user: &usersd.User{
				Phone:    "+12345678901",
				Password: "1234",
			},
		},

		{
			name: "Update",
			user: &usersd.User{
				ID:       "test",
				Email:    "john@example.com",
				Phone:    "+12345678901",
				Password: "1234",
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
			err := ud.WriteUser(c.user)

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

func usersFixtures(t *testing.T, tx *usersd.Tx) {
	users := []*usersd.User{
		{
			ID:       "admin",
			Email:    "admin@example.com",
			Password: "admin",
		},

		{
			Email:    "john@example.com",
			Phone:    "+12345678901",
			Password: "1234",
		},

		{
			Email:    "john2@example.com",
			Password: "1234",
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
