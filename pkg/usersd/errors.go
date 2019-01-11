// Copyright 2018 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the WTFPL.

package usersd

import (
	"fmt"
)

// User errors.
var (
	ErrUserNotFound = Error{
		Code:    1,
		Message: "the given user doesn't exists",
	}

	ErrUserIDCreation = Error{
		Code:    2,
		Message: "can't generate a unique user ID -> %s",
	}

	ErrUserEmailEmpty = Error{
		Code:    10,
		Field:   "email",
		Message: "the given email is empty",
	}

	ErrUserEmailInvalid = Error{
		Code:    11,
		Field:   "email",
		Message: "the given email is invalid",
	}

	ErrUserEmailExists = Error{
		Code:    12,
		Field:   "email",
		Message: "the given email already exists",
	}

	ErrUserPhoneEmpty = Error{
		Code:    20,
		Field:   "phone",
		Message: "the given phone is empty",
	}

	ErrUserPhoneInvalid = Error{
		Code:    21,
		Field:   "phone",
		Message: "the given phone is invalid",
	}

	ErrUserPhoneExists = Error{
		Code:    22,
		Field:   "phone",
		Message: "the given phone already exists",
	}
)

// Password errors.
var (
	ErrPasswordEmpty = Error{
		Code:    50,
		Field:   "password",
		Message: "the given password is empty",
	}
)

// Error is a more descriptive error.
type Error struct {
	Code    int    `json:"code"`
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Error implements error.
func (e Error) Error() string {
	if e.Field == "" {
		return fmt.Sprintf("(%d) %s", e.Code, e.Message)
	}

	return fmt.Sprintf("(%d) %s: %s", e.Code, e.Field, e.Message)
}

// Format returns a new error with a formated message.
func (e Error) Format(args ...interface{}) error {
	e.Message = fmt.Sprintf(e.Message, args...)
	return e
}

// Errors is a set of errors wrapped into a single error.
type Errors []error

// Error implements error.
func (e Errors) Error() string {
	errors := ""

	if len(e) > 0 {
		for _, err := range e {
			errors += "; " + err.Error()
		}

		return errors[2:]
	}

	return errors
}
