// Copyright 2018 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the WTFPL.

package usersd

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/dgraph-io/badger"
	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

// User is an entity that may be authenticated and authorized.
//
// There are some special keys at User.Data field:
//
// * createdAt: Creation date.
//
// * lastLogin: Last login date.
type User struct {
	mu       sync.Mutex
	ID       string                 `json:"id"`
	Password string                 `json:"password"`
	Data     map[string]interface{} `json:"data"`
}

// NewUser creates a user with the given arguments and populates some required
// data if missing.
func NewUser(id, password string, data map[string]interface{}) (*User, error) {
	if data == nil {
		data = make(map[string]interface{})
	}

	u := &User{
		ID:   id,
		Data: data,
	}

	if err := u.SetPassword(password); err != nil {
		l.Printf(lError+" Can't create the password hash -> %v", err)
		return nil, err
	}

	if debug {
		l.Printf(lDebug+" Temporary user created (%v) -> %+v", u.ID, u)
	}

	return u, nil
}

// NewUserJSON creates a user with the given JSON data and populates some
// required data if missing.
func NewUserJSON(data []byte) (*User, error) {
	u := new(User)

	if err := json.Unmarshal(data, u); err != nil {
		return nil, err
	}

	u.SetPassword(u.Password)
	return u, nil
}

// ListUsers fetches users that satisfies the given constraints.
func ListUsers() ([]*User, error) {
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

		v, err := item.ValueCopy(v)

		if err != nil {
			return nil, err
		}

		user, err := NewUserJSON(v)

		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

// GetUser fetches a user with the given ID.
func GetUser(id string) (*User, error) {
	txn := db.NewTransaction(false)
	defer txn.Discard()

	item, err := txn.Get([]byte(id))

	if err != nil {
		l.Printf(lError+" Can't find the given user (%s) -> %v", id, err)
		return nil, err
	}

	data, err := item.ValueCopy(nil)

	if err != nil {
		l.Printf(lError+" Can't fetch the user data -> %v", err)
		return nil, err
	}

	user, err := NewUserJSON(data)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// Delete removes the user from the database.
func (u *User) Delete() error {
	txn := db.NewTransaction(true)
	defer txn.Discard()

	if err := txn.Delete([]byte(u.ID)); err != nil {
		msg := lError + " Can't remove the user (%v) from the database -> %v"
		l.Printf(msg, u.ID, err)
		return err
	}

	if err := index.Delete(u.ID); err != nil {
		msg := lError + " Can't remove the user (%v) from the search index -> %v"
		l.Printf(msg, u.ID, err)
		return err
	}

	if err := txn.Commit(nil); err != nil {
		msg := lError + " Can't commit changes to the database -> %v"
		l.Printf(msg, err)
		return err
	}

	if debug {
		l.Printf(lDebug+" User (%v) data removed", u.ID)
	}

	return nil
}

// Save writes the user data to the database.
func (u *User) Save() error {
	if u.ID == "" {
		x, err := uuid.NewV4()

		if err != nil {
			return err
		}

		u.ID = x.String()
	}

	u.Set("createdAt", time.Now().Format("2006-01-02T15:04:05-0700"))
	v, err := json.Marshal(u)

	if err != nil {
		l.Printf(lError+" Can't serialize the user (%v) -> %v", u.ID, err)
		return err
	}

	txn := db.NewTransaction(true)
	defer txn.Discard()

	if err := txn.Set([]byte(u.ID), v); err != nil {
		msg := lError + " Can't write the user (%v) to the database -> %v"
		l.Printf(msg, u.ID, err)
		return err
	}

	if err := index.Index(u.ID, u); err != nil {
		msg := lError + " Can't add the user (%v) to the search index -> %v"
		l.Printf(msg, u.ID, err)
		return err
	}

	if err := txn.Commit(nil); err != nil {
		msg := lError + " Can't commit chages to the database -> %v"
		l.Printf(msg, err)
		return err
	}

	if debug {
		l.Printf(lDebug+" User (%v) data saved -> '%v'", u.ID, string(v))
	}

	return nil
}

// Set sets the given value at the given key at Data field.
func (u *User) Set(key string, value interface{}) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.Data[key] = value
}

// SetPassword sets the user password from a string and returns an error if
// any.
func (u *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)

	if err != nil {
		return err
	}

	u.Password = string(hash)
	return nil
}
