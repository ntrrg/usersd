// Copyright 2018 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the WTFPL.

package usersd

import (
	"crypto/rand"

	"golang.org/x/crypto/argon2"
)

const passwordsDI = "passwords"

// Password errors.
var (
	ErrPasswordEmpty = Error{
		Code:    50,
		Type:    "password",
		Message: "the given password is empty",
	}
)

// PasswordOptions wraps variables that control the passwords hashing algorithm
// behavior.
type PasswordOptions struct {
	SaltSize uint32
	Time     uint32
	Memory   uint32
	Threads  byte
	HashSize uint32
}

// CheckPassword is a helper for Service.CheckPasswordTx.
// true if match.
func (s *Service) CheckPassword(userid, password string) bool {
	if password == "" {
		return false
	}

	tx := s.NewTx(false)
	defer tx.Discard()

	user, err := GetUser(tx, userid)
	if err != nil {
		return false
	}

	return s.CheckPasswordTx(tx, user, password)
}

// CheckPasswordTx compares the given password with the user password and
// returns true if match.
func (s *Service) CheckPasswordTx(tx *Tx, user *User, password string) bool {
	if password == "" {
		return false
	}

	data, err := tx.Get([]byte(passwordsDI + user.ID))
	if err != nil {
		return false
	}

	opts := s.opts.PasswdOpts
	salt, oldhash := data[:opts.SaltSize], data[opts.SaltSize:]

	hash := argon2.IDKey(
		[]byte(password), salt,
		opts.Time, opts.Memory, opts.Threads, opts.HashSize,
	)

	for i, v := range oldhash {
		if v != hash[i] {
			return false
		}
	}

	return true
}

// SetPassword is a helper for Service.SetPasswordTx.
func (s *Service) SetPassword(userid, password string) error {
	if password == "" {
		return ErrPasswordEmpty
	}

	tx := s.NewTx(true)
	defer tx.Discard()

	user, err := GetUser(tx, userid)
	if err != nil {
		return err
	}

	if err := s.SetPasswordTx(tx, user, password); err != nil {
		return err
	}

	return tx.Commit()
}

// SetPasswordTx assigns password to the given user.
func (s *Service) SetPasswordTx(tx *Tx, user *User, password string) error {
	if password == "" {
		return ErrPasswordEmpty
	}

	if _, err := GetUser(tx, user.ID); err != nil {
		return err
	}

	opts := s.opts.PasswdOpts
	salt := make([]byte, opts.SaltSize)
	rand.Read(salt)

	hash := argon2.IDKey(
		[]byte(password), salt,
		opts.Time, opts.Memory, opts.Threads, opts.HashSize,
	)

	return tx.Set([]byte(passwordsDI+user.ID), append(salt, hash...))
}
