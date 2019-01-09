// Copyright 2018 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the WTFPL.

package usersd

import (
	"fmt"
)

// Error is a more descriptive error.
type Error struct {
	Code    int    `json:"code"`
	Type    string `json:"field"`
	Message string `json:"message"`
}

// Error implements error.
func (e Error) Error() string {
	return fmt.Sprintf("(%d) %s: %s", e.Code, e.Type, e.Message)
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
