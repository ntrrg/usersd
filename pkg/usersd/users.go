// Copyright 2018 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the WTFPL.

package usersd

import (
	"bytes"
	"encoding/json"

	"github.com/dgraph-io/badger"
	"github.com/gofrs/uuid"
)

// User represents a player.
type User struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Games     uint64  `json:"gamesPlayed"`
	HighScore uint64  `json:"highscore,omitempty"`
	Friends   []*User `json:"friends,omitempty"`

	Raw []byte `json:"-"`
}

// NewUser creates a user with the given name and populates some required data.
func NewUser(name string) (*User, error) {
	user := &User{Name: name}
	id, err := uuid.NewV4()

	if err != nil {
		return nil, err
	}

	user.ID = id.String()

	if err := user.Save(); err != nil {
		return nil, err
	}

	return user, nil
}

// NewUserFromJSON creates a user from the given JSON data and populates some
// required data if missing.
func NewUserFromJSON(data []byte) (*User, error) {
	user := new(User)

	if err := user.Load(data); err != nil {
		return nil, err
	}

	if user.ID == "" {
		id, err := uuid.NewV4()

		if err != nil {
			return nil, err
		}

		user.ID = id.String()
	}

	if err := user.Save(); err != nil {
		return nil, err
	}

	return user, nil
}

// GetUser fetch a user with the given ID from the database.
func GetUser(id string) (*User, error) {
	data, err := Get([]byte(id))

	if err != nil {
		return nil, err
	}

	user := new(User)

	if err := user.Load(data); err != nil {
		return nil, err
	}

	return user, nil
}

// GetUsers fetch users that satisfies the given constraints from the database.
func GetUsers() ([]*User, error) {
	txn := db.NewTransaction(false)
	defer txn.Discard()

	var users []*User

	opts := badger.DefaultIteratorOptions
	// opts.PrefetchValues = false
	it := txn.NewIterator(opts)
	defer it.Close()

	var v []byte

	for it.Rewind(); it.Valid(); it.Next() {
		item := it.Item()

		if _, err := item.ValueCopy(v); err != nil {
			return nil, err
		}

		user := &User{ID: string(item.Key())}

		if err := user.Load(v); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

// Load fills an user with the given JSON data, but doesn't writes to the
// database.
func (u *User) Load(data []byte) error {
	if bytes.Equal(data, u.Raw) {
		return nil
	}

	if err := json.Unmarshal(data, u); err != nil {
		return err
	}

	u.Raw = data
	return nil
}

// Save writes the user data to the database.
func (u *User) Save() error {
	if err := u.Validate(); err != nil {
		return err
	}

	v, err := json.Marshal(u)

	if err != nil {
		return err
	}

	if bytes.Equal(v, u.Raw) {
		return nil
	}

	txn := db.NewTransaction(true)
	defer txn.Discard()

	if err := txn.Set([]byte(u.ID), v); err != nil {
		return err
	}

	if err := index.Index(u.ID, v); err != nil {
		return err
	}

	if err := txn.Commit(nil); err != nil {
		return err
	}

	u.Raw = v
	return nil
}

// Validate checks for invalid data from a user.
func (u *User) Validate() error {
	return nil
}
