// Copyright 2018 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the WTFPL.

package usersd

import (
	"regexp"
	"time"

	"github.com/blevesearch/bleve"
	"github.com/dgraph-io/badger"
	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Regular expressions
var (
	emailRegexp = regexp.MustCompile(`^[a-zA-Z0-9.!#$%&'*+\/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`) // nolint: lll
)

func userCreatedAtValidator(tx *badger.Txn, index bleve.Index, user *User, old *User) error { // nolint: lll
	if old == nil {
		user.CreatedAt = time.Now().Unix()
	} else {
		user.CreatedAt = old.CreatedAt
	}

	return nil
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

func userEmailValidator(tx *badger.Txn, index bleve.Index, user *User, old *User) error { // nolint: lll
	if user.Email == "" {
		return nil
	}

	if !emailRegexp.MatchString(user.Email) {
		return ErrUserEmailInvalid
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

func userLastLoginValidator(tx *badger.Txn, index bleve.Index, user *User, old *User) error { // nolint: lll
	if old == nil {
		user.LastLogin = 0
	}

	return nil
}

func userModeValidator(tx *badger.Txn, index bleve.Index, user *User, old *User) error { // nolint: lll
	if user.Mode == "" || user.EmailVerified || user.PhoneVerified {
		user.Mode = defaultUserMode
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
