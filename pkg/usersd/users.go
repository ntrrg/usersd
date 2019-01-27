// Copyright 2018 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the WTFPL.

package usersd

import (
	"encoding/json"

	"github.com/blevesearch/bleve"
	"github.com/dgraph-io/badger"
)

// Users documents identifier.
const UsersDI = "users"

const defaultUserMode = "local"

// User is an entity that may be authenticated and authorized. See also:
//
// * Tx.WriteUser: write the user.
//
// * Tx.GetUser: get user by ID.
//
// * Tx.GetUsers: get users that satisfies some constraints.
//
// * Tx.DeleteUser: delete the given user.
type User struct {
	ID        string `json:"id"`
	Mode      string `json:"mode"`
	CreatedAt int64  `json:"createdAt"`
	LastLogin int64  `json:"lastLogin"`

	Data map[string]interface{} `json:"data,omitempty"`

	// Verification methods.
	Email         string `json:"email"`
	Phone         string `json:"phone"`
	EmailVerified bool   `json:"emailVerified"`
	PhoneVerified bool   `json:"phoneVerified"`

	Roles []string `json:"roles"`
}

// GetUser fetches a user with the given ID from the DB.
func (tx *Tx) GetUser(id string) (*User, error) {
	if id == "" {
		return nil, ErrUserNotFound
	}

	data, err := tx.Get([]byte(UsersDI + id))
	if err == badger.ErrKeyNotFound {
		return nil, ErrUserNotFound
	} else if err != nil {
		return nil, err
	}

	user := new(User)
	if err := json.Unmarshal(data, user); err != nil {
		return nil, err
	}

	return user, nil
}

// GetUsers fetches users from the DB that satisfies the given constraints.
func (tx *Tx) GetUsers(q string, sort ...string) ([]*User, error) {
	if q == "" && len(sort) == 0 {
		return getAllUsers(tx)
	}

	users := []*User{}
	bq := bleve.NewQueryStringQuery(q + " +documenttype:" + UsersDI)
	req := bleve.NewSearchRequest(bq)
	req.SortBy(sort)

	res, err := usersd.index.Search(req)
	if err != nil {
		return nil, err
	}

	for _, hit := range res.Hits {
		user, err := tx.GetUser(hit.ID)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

// DeleteUser removes the user from the DB.
func (tx *Tx) DeleteUser(id string) error {
	if err := tx.Delete([]byte(UsersDI + id)); err != nil {
		return err
	}

	if err := tx.Delete([]byte(PasswordsDI + id)); err != nil {
		return err
	}

	return usersd.index.Delete(id)
}

// ValidateUser checks the user data and returns any errors.
func (tx *Tx) ValidateUser(user *User) error {
	old, err := tx.GetUser(user.ID)
	if err != nil && err != ErrUserNotFound {
		return err
	}

	var errors Errors

	rules := []func(*Tx, *User, *User) error{
		userIDValidator,
		userEmailValidator,
		userEmailVerifiedValidator,
		userPhoneValidator,
		userPhoneVerifiedValidator,
		userModeValidator,
		userCreatedAtValidator,
		userLastLoginValidator,
	}

	for _, f := range rules {
		if err := f(tx, user, old); err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) != 0 {
		return errors
	}

	return nil
}

// WriteUser writes the user data to the DB.
func (tx *Tx) WriteUser(user *User) error {
	if err := tx.ValidateUser(user); err != nil {
		return err
	}

	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	if err := tx.Set([]byte(UsersDI+user.ID), data); err != nil {
		return err
	}

	v := make(map[string]interface{})
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	v["documenttype"] = UsersDI
	return usersd.index.Index(user.ID, v)
}

func getAllUsers(tx *Tx) ([]*User, error) {
	var users []*User

	it := tx.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()

	var (
		v []byte

		prefix = []byte(UsersDI)
	)

	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		var err error
		item := it.Item()
		v, err = item.ValueCopy(v)
		if err != nil {
			return nil, err
		}

		user := new(User)
		if err := json.Unmarshal(v, user); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}
