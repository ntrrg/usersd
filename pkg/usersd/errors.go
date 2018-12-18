// Copyright 2018 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the WTFPL.

package usersd

import (
	"fmt"
)

// User validation errors
var (
	ErrUserIDNotFound = ValidationError{
		Code:    1,
		Field:   "id",
		Message: "The given user ID doesn't exists",
	}

	ErrUserEmailEmpty = ValidationError{
		Code:    2,
		Field:   "email",
		Message: "The given email is empty",
	}

	ErrUserPasswordEmpty = ValidationError{
		Code:    3,
		Field:   "password",
		Message: "The given password is empty",
	}

	ErrUserPasswordHash = ValidationError{
		Code:    4,
		Field:   "password",
		Message: "Can't encrypt the password -> %s",
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

// ValidationErrors is a set of errors after validating an entity.
type ValidationErrors []error

func (e ValidationErrors) Error() string {
	errors := ""

	if len(e) > 0 {
		for _, err := range e {
			errors += "; " + err.Error()
		}

		return errors[2:]
	}

	return errors
}
