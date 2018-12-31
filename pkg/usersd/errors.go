// Copyright 2018 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the WTFPL.

package usersd

import (
	"fmt"
)

// User errors.
var (
	ErrUserIDNotFound = ValidationError{
		Code:    1,
		Field:   "id",
		Message: "the given user ID doesn't exists",
	}

	ErrUserIDCreation = ValidationError{
		Code:    2,
		Field:   "id",
		Message: "can't generate the user ID -> %s",
	}

	ErrUserEmailEmpty = ValidationError{
		Code:    10,
		Field:   "email",
		Message: "the given email is empty",
	}

	ErrUserEmailInvalid = ValidationError{
		Code:    11,
		Field:   "email",
		Message: "the given email is invalid",
	}

	ErrUserEmailExists = ValidationError{
		Code:    12,
		Field:   "email",
		Message: "the given email already exists",
	}

	ErrUserPhoneEmpty = ValidationError{
		Code:    20,
		Field:   "phone",
		Message: "the given phone is empty",
	}

	ErrUserPhoneInvalid = ValidationError{
		Code:    21,
		Field:   "phone",
		Message: "the given phone is invalid",
	}

	ErrUserPhoneExists = ValidationError{
		Code:    22,
		Field:   "phone",
		Message: "the given phone already exists",
	}

	ErrUserPasswordEmpty = ValidationError{
		Code:    30,
		Field:   "password",
		Message: "the given password is empty",
	}

	ErrUserPasswordHash = ValidationError{
		Code:    31,
		Field:   "password",
		Message: "can't encrypt the password -> %s",
	}
)

// ValidationError is an error after validating an entity field.
type ValidationError struct {
	Code    int    `json:"code"`
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Error implements error.
func (e ValidationError) Error() string {
	return fmt.Sprintf("(%s) %s", e.Field, e.Message)
}

// Format returns a new error with a formated message.
func (e ValidationError) Format(args ...interface{}) error {
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
