// Copyright 2018 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the WTFPL.

package usersd

import (
	"encoding/json"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/search/query"
	"github.com/dgraph-io/badger"
)

const (
	usersDI         = "users"
	defaultUserMode = "local"
)

// User errors.
var (
	ErrUserIDNotFound = Error{
		Code:    1,
		Type:    "id",
		Message: "the given user ID doesn't exists",
	}

	ErrUserIDCreation = Error{
		Code:    2,
		Type:    "id",
		Message: "can't generate the user ID -> %s",
	}

	ErrUserEmailEmpty = Error{
		Code:    10,
		Type:    "email",
		Message: "the given email is empty",
	}

	ErrUserEmailInvalid = Error{
		Code:    11,
		Type:    "email",
		Message: "the given email is invalid",
	}

	ErrUserEmailExists = Error{
		Code:    12,
		Type:    "email",
		Message: "the given email already exists",
	}

	ErrUserPhoneEmpty = Error{
		Code:    20,
		Type:    "phone",
		Message: "the given phone is empty",
	}

	ErrUserPhoneInvalid = Error{
		Code:    21,
		Type:    "phone",
		Message: "the given phone is invalid",
	}

	ErrUserPhoneExists = Error{
		Code:    22,
		Type:    "phone",
		Message: "the given phone already exists",
	}
)

// User is an entity that may be authenticated and authorized.
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
}

// GetUser fetches a user with the given ID from the database.
func GetUser(tx *Tx, id string) (*User, error) {
	if id == "" {
		return nil, ErrUserIDNotFound
	}

	data, err := tx.Get([]byte(usersDI + id))
	if err == badger.ErrKeyNotFound {
		return nil, ErrUserIDNotFound
	} else if err != nil {
		return nil, err
	}

	user := new(User)
	if err := json.Unmarshal(data, user); err != nil {
		return nil, err
	}

	return user, nil
}

// GetUsers fetches users that satisfies the given constraints.
func GetUsers(tx *Tx, q string, sort ...string) ([]*User, error) {
	if q == "" && len(sort) == 0 {
		return getAllUsers(tx)
	}

	var (
		users = []*User{}

		bq query.Query
	)

	if q != "" {
		bq = bleve.NewQueryStringQuery(q + " +documenttype:" + usersDI)
	} else {
		bq = bleve.NewMatchAllQuery()
	}

	req := bleve.NewSearchRequest(bq)
	req.SortBy(sort)

	res, err := tx.Index.Search(req)
	if err != nil {
		return nil, err
	}

	for _, hit := range res.Hits {
		user, err := GetUser(tx, hit.ID)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

// Delete removes the user from the database.
func (u *User) Delete(tx *Tx) error {
	if err := tx.Delete([]byte(usersDI + u.ID)); err != nil {
		return err
	}

	return tx.Index.Delete(u.ID)
}

// Validate checks the user data and returns any errors.
func (u *User) Validate(tx *Tx) error {
	old, err := GetUser(tx, u.ID)
	if err != nil && err != ErrUserIDNotFound {
		return err
	}

	errors := Errors{}

	rules := []func(tx *Tx, user *User, old *User) error{
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
		if err := f(tx, u, old); err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}

// Write writes the user data to the database.
func (u *User) Write(tx *Tx) error {
	if err := u.Validate(tx); err != nil {
		return err
	}

	data, err := json.Marshal(u)
	if err != nil {
		return err
	}

	if err := tx.Set([]byte(usersDI+u.ID), data); err != nil {
		return err
	}

	v := make(map[string]interface{})
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	v["documenttype"] = usersDI
	return tx.Index.Index(u.ID, v)
}

// GetUser is a helper method for GetUser.
func (s *Service) GetUser(id string) (*User, error) {
	tx := s.NewTx(false)
	defer tx.Discard()
	return GetUser(tx, id)
}

// GetUsers is a helper method for GetUsers.
func (s *Service) GetUsers(q string, sort ...string) ([]*User, error) {
	tx := s.NewTx(false)
	defer tx.Discard()
	return GetUsers(tx, q, sort...)
}

// DeleteUser is a helper method for User.Delete.
func (s *Service) DeleteUser(id string) error {
	user := &User{ID: id}
	tx := s.NewTx(true)
	defer tx.Discard()

	if err := user.Delete(tx); err != nil {
		return err
	}

	return tx.Commit()
}

// WriteUser is a helper method for User.Write.
func (s *Service) WriteUser(user *User) error {
	tx := s.NewTx(true)
	defer tx.Discard()

	if err := user.Write(tx); err != nil {
		return err
	}

	return tx.Commit()
}

func getAllUsers(tx *Tx) ([]*User, error) {
	var users []*User

	it := tx.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()

	var (
		v []byte

		prefix = []byte(usersDI)
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
