// Copyright 2018 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the WTFPL.

package usersd_test

import (
	"io/ioutil"
	"path/filepath"

	"github.com/ntrrg/usersd/pkg/usersd"
)

var testDir string

func usersFixtures(tx *usersd.Tx) error {
	entries := []struct {
		user     *usersd.User
		password string
	}{
		{
			user: &usersd.User{
				ID:    "admin",
				Email: "admin@example.com",
			},
			password: "admin",
		},

		{
			user: &usersd.User{
				Mode:  "oauth2",
				Email: "john@example.com",
				Phone: "+12345678901",
			},
		},

		{
			user: &usersd.User{
				Email: "john2@example.com",
				Data: map[string]interface{}{
					"username": "john",
					"name":     "John Doe",
				},
			},
			password: "jhon1234",
		},
	}

	for _, entry := range entries {
		if err := tx.WriteUser(entry.user); err != nil {
			return err
		}

		if entry.password != "" {
			if err := tx.SetPassword(entry.user.ID, entry.password); err != nil {
				return err
			}
		}
	}

	return nil
}

func initTest(name string, fixtures bool) error {
	opts := usersd.DefaultOptions
	opts.Database = filepath.Join(testDir, name)
	if err := usersd.Init(opts); err != nil {
		return err
	}

	if fixtures {
		tx := usersd.NewTx(true)
		defer tx.Discard()

		fns := []func(*usersd.Tx) error{
			usersFixtures,
		}

		for _, fn := range fns {
			if err := fn(tx); err != nil {
				return err
			}
		}

		if err := tx.Commit(); err != nil {
			return err
		}
	}

	return nil
}

func init() {
	var err error
	testDir, err = ioutil.TempDir("", "usersd-tests")
	if err != nil {
		panic(err)
	}
}
