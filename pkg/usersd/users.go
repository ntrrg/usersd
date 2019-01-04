// Copyright 2018 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the WTFPL.

package usersd

import (
	"encoding/json"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/search/query"
	"github.com/dgraph-io/badger"
	"golang.org/x/crypto/bcrypt"
)

const (
	dbKeyPrefixUsers = "users-"
	defaultUserMode  = "local"

	// Password hashing strength.
	bcryptCost = 10
)

// User is an entity that may be authenticated and authorized.
type User struct {
	ID        string `json:"id"`
	Mode      string `json:"mode"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Password  string `json:"password,omitempty"`
	CreatedAt int64  `json:"createdAt"`
	LastLogin int64  `json:"lastLogin"`

	Data map[string]interface{} `json:"data,omitempty"`

	EmailVerified bool `json:"emailVerified"`
	PhoneVerified bool `json:"phoneVerified"`
}

// GetUser fetches a user with the given ID from the database.
func GetUser(tx *Tx, id string) (*User, error) {
	if id == "" {
		return nil, ErrUserIDNotFound
	}

	defer func() {
		recover() // nolint: errcheck
	}()

	item, err := tx.Get([]byte(dbKeyPrefixUsers + id))
	if err == badger.ErrKeyNotFound {
		return nil, ErrUserIDNotFound
	} else if err != nil {
		return nil, err
	}

	data, err := item.ValueCopy(nil)
	if err != nil {
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
		bq = bleve.NewQueryStringQuery(q)
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

// CheckPassword compares the given password with the user password and returns
// true if match.
func (u *User) CheckPassword(password string) bool {
	if u.Password == "" {
		return false
	}

	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// Delete removes the user from the database.
func (u *User) Delete(tx *Tx) error {
	if err := tx.Delete([]byte(dbKeyPrefixUsers + u.ID)); err != nil {
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
		userPasswordValidator,
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

	v, err := json.Marshal(u)
	if err != nil {
		return err
	}

	if err := tx.Set([]byte(dbKeyPrefixUsers+u.ID), v); err != nil {
		return err
	}

	return tx.Index.Index(u.ID, u)
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

		prefix = []byte(dbKeyPrefixUsers)
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
