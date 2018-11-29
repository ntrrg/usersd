// Copyright 2018 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the WTFPL.

package usersd

import (
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/dgraph-io/badger"
	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Errors
var (
	ErrUserIDInUse         = errors.New("The given ID is already used")
	ErrUserPasswordEmpty   = errors.New("The password is empty")
	ErrUserPasswordBadType = errors.New("Bad type for password")
	ErrUserNotFound        = errors.New("The given user doesn't exists")
)

// User is an entity that may be authenticated and authorized.
//
// There are some special keys at User.Data field:
//
// * createdAt: Creation date.
//
// * lastLogin: Last login date.
type User struct {
	ID string `json:"id"`

	mu   sync.Mutex
	Data map[string]interface{} `json:"data"`

	password string
}

// NewUser creates a user with the given arguments and populates some required
// data if missing.
func NewUser(id string, data map[string]interface{}) *User {
	if data == nil {
		data = make(map[string]interface{})
	}

	u := &User{
		ID:   id,
		Data: data,
	}

	return u
}

// NewUserJSON creates a user with the given JSON data and populates some
// required data if missing.
func NewUserJSON(data []byte) (*User, error) {
	u := NewUser("", nil)

	if err := json.Unmarshal(data, u); err != nil {
		l.Printf(lError+" Can't deserialize the user data -> %v", err)
		return nil, err
	}

	return u, nil
}

// CreateUser creates a user with the given arguments, populates some required
// data if missing and write it to the database.
func CreateUser(id string, data map[string]interface{}) (*User, error) {
	if id == "" {
		for {
			x, err := uuid.NewV4()

			if err != nil {
				return nil, err
			}

			id = x.String()

			if _, err := GetUser(id); err == badger.ErrKeyNotFound {
				break
			} else if err != nil {
				return nil, err
			}
		}
	} else {
		if _, err := GetUser(id); err == nil {
			return nil, ErrUserIDInUse
		}
	}

	u := NewUser(id, data)

	if err := u.write(); err != nil {
		return nil, err
	}

	return u, nil
}

// CreateUserJSON creates a user with the given JSON, populates some required
// data if missing and write it to the database.
func CreateUserJSON(data []byte) (*User, error) {
	u, err := NewUserJSON(data)

	if err != nil {
		return nil, err
	}

	return CreateUser(u.ID, u.Data)
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
			l.Printf(lError+" Can't fetch the user data -> %v", err)
			return nil, err
		}

		user, err := NewUserJSON(v)

		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if debug {
		l.Printf(lDebug+" %v users fetched", len(users))
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

	if err := txn.Commit(); err != nil {
		msg := lError + " Can't commit changes to the database -> %v"
		l.Printf(msg, err)
		return err
	}

	if debug {
		l.Printf(lDebug+" User (%v) data removed", u.ID)
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
// any. If the given string is empty, SetPassword reads from Data["password"],
// and if it is empty too, ErrUserPasswordEmpty is returned.
func (u *User) SetPassword(password string) error {
	if password == "" {
		value, ok := u.Data["password"]

		if !ok {
			l.Printf(lError+" The password is empty -> %v", u.ID)
			return ErrUserPasswordEmpty
		}

		password, ok = value.(string)

		if !ok {
			l.Printf(lError+" Bad type for password (%v) -> %v", u.ID, value)
			return ErrUserPasswordBadType
		}
	}

	if _, err := bcrypt.Cost([]byte(password)); err == nil {
		return nil
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)

	if err != nil {
		msg := lError + " Can't create the hash from the user (%v) password -> %v"
		l.Printf(msg, u.ID, err)
		return err
	}

	u.password = string(hash)
	u.Unset("password")
	return nil
}

// Unset removes the given key from Data field.
func (u *User) Unset(key string) {
	u.mu.Lock()
	defer u.mu.Unlock()
	delete(u.Data, key)
}

// Update updates the user data in the database.
func (u *User) Update() error {
	if _, err := GetUser(u.ID); err == badger.ErrKeyNotFound {
		l.Printf(lError+" Can't update, the user doesn't exists (%v)", u.ID)
		return ErrUserNotFound
	} else if err != nil {
		return err
	}

	return u.write()
}

// write writes the user data to the database.
func (u *User) write() error {
	if err := u.SetPassword(""); err != nil && err != ErrUserPasswordEmpty {
		return err
	}

	u.Set("id", u.ID)

	if _, ok := u.Data["createdAt"]; !ok {
		u.Set("createdAt", time.Now().Format("2006-01-02T15:04:05-0700"))
	}

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
		msg := lError + " Can't write the user (%v) to the search index -> %v"
		l.Printf(msg, u.ID, err)
		return err
	}

	if err := txn.Commit(); err != nil {
		msg := lError + " Can't commit chages to the database -> %v"
		l.Printf(msg, err)
		return err
	}

	if debug {
		l.Printf(lDebug+" User (%v) data saved -> '%v'", u.ID, string(v))
	}

	return nil
}
