// Copyright 2018 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the WTFPL.

package rest

// Error represents an error during the execution of any handler. It could be
// just one error or a collection of them.
type Error struct {
	Code    int     `json:"code"`
	Message string  `json:"message"`
	Errors  []Error `json:"errors,omitempty"`
}

// ErrInternal is used for internal errors.
var ErrInternal = Error{Code: 0, Message: "Internal server error"}

// Users errors (1XX codes).
var (
	ErrGetUsers = Error{Code: 100, Message: "Can't get the users list"}

	ErrUnmarshalUser = Error{
		Code:    101,
		Message: "Can't unmarshal the user from the request body",
	}
)
