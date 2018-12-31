// Copyright 2018 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the WTFPL.

package usersd

import (
	"encoding/json"
	"time"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/search/query"
	"github.com/dgraph-io/badger"
	"github.com/gofrs/uuid"
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

	Roles []string `json:"roles,omitempty"`

	Data map[string]interface{} `json:"data,omitempty"`

	EmailVerified bool `json:"emailVerified"`
	PhoneVerified bool `json:"phoneVerified"`
}

// GetUser fetches a user with the given ID from the database.
func GetUser(tx *badger.Txn, id string) (*User, error) {
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
func GetUsers(tx *badger.Txn, index bleve.Index, q string, sort ...string) ([]*User, error) { // nolint: lll
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

	res, err := index.Search(req)
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
func (u *User) Delete(tx *badger.Txn, index bleve.Index) error {
	if err := tx.Delete([]byte(dbKeyPrefixUsers + u.ID)); err != nil {
		return err
	}

	return index.Delete(u.ID)
}

// Validate checks the user data and returns any errors.
func (u *User) Validate(tx *badger.Txn, index bleve.Index) error {
	old, err := GetUser(tx, u.ID)
	if err != nil && err != ErrUserIDNotFound {
		return err
	}

	errors := Errors{}

	rules := []func(tx *badger.Txn, index bleve.Index, user *User, old *User) error{ // nolint: lll
		userIDValidator,
		userModeValidator,
		userEmailValidator,
		userEmailVerifiedValidator,
		userPasswordValidator,
		userPhoneValidator,
		userPhoneVerifiedValidator,
		userCreatedAtValidator,
		userLastLoginValidator,
	}

	for _, f := range rules {
		if err := f(tx, index, u, old); err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}

// Write writes the user data to the database.
func (u *User) Write(tx *badger.Txn, index bleve.Index) error { // nolint: lll
	if err := u.Validate(tx, index); err != nil {
		return err
	}

	v, err := json.Marshal(u)
	if err != nil {
		return err
	}

	if err := tx.Set([]byte(dbKeyPrefixUsers+u.ID), v); err != nil {
		return err
	}

	return index.Index(u.ID, u)
}

func getAllUsers(tx *badger.Txn) ([]*User, error) {
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

func userIDValidator(tx *badger.Txn, index bleve.Index, user *User, old *User) error { // nolint: lll
	if user.ID == "" {
		x, err := uuid.NewV4()
		if err != nil {
			return ErrUserIDCreation.Format(err)
		}

		user.ID = x.String()
	}

	return nil
}

func userModeValidator(tx *badger.Txn, index bleve.Index, user *User, old *User) error { // nolint: lll
	if user.Mode == "" || user.EmailVerified || user.PhoneVerified {
		user.Mode = defaultUserMode
	}

	return nil
}

func userEmailValidator(tx *badger.Txn, index bleve.Index, user *User, old *User) error { // nolint: lll
	if user.Email == "" {
		return ErrUserEmailEmpty
	}

	q := `+email:"` + user.Email + `"`

	if old != nil {
		q = `-id:"` + old.ID + `" ` + q
	}

	users, err := GetUsers(tx, index, q)
	if err == nil && len(users) > 0 {
		return ErrUserEmailExists
	}

	return nil
}

func userEmailVerifiedValidator(tx *badger.Txn, index bleve.Index, user *User, old *User) error { // nolint: lll
	if old == nil || user.Email != old.Email {
		user.EmailVerified = false
	}

	return nil
}

func userPhoneValidator(tx *badger.Txn, index bleve.Index, user *User, old *User) error { // nolint: lll
	if user.Phone == "" {
		return nil
	}

	q := `+phone:"` + user.Phone + `"`

	if old != nil {
		q = `-id:"` + old.ID + `" ` + q
	}

	users, err := GetUsers(tx, index, q)
	if err == nil && len(users) > 0 {
		return ErrUserPhoneExists
	}

	return nil
}

func userPhoneVerifiedValidator(tx *badger.Txn, index bleve.Index, user *User, old *User) error { // nolint: lll
	if old == nil || user.Phone != old.Phone {
		user.PhoneVerified = false
	}

	return nil
}

func userPasswordValidator(tx *badger.Txn, index bleve.Index, user *User, old *User) error { // nolint: lll
	if user.Mode == defaultUserMode && user.Password == "" {
		return ErrUserPasswordEmpty
	} else if user.Mode != defaultUserMode && user.Password == "" {
		return nil
	}

	password := user.Password

	if _, err := bcrypt.Cost([]byte(password)); err == nil {
		return nil
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return ErrUserPasswordHash.Format(err)
	}

	user.Password = string(hash)
	return nil
}

func userCreatedAtValidator(tx *badger.Txn, index bleve.Index, user *User, old *User) error { // nolint: lll
	if old == nil {
		user.CreatedAt = time.Now().Unix()
	} else {
		user.CreatedAt = old.CreatedAt
	}

	return nil
}

func userLastLoginValidator(tx *badger.Txn, index bleve.Index, user *User, old *User) error { // nolint: lll
	if old == nil {
		user.LastLogin = 0
	}

	return nil
}
