// Copyright 2018 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the WTFPL.

package rest

type Error struct {
	Code    int     `json:"code"`
	Message string  `json:"message"`
	Errors  []Error `json:"errors,omitempty"`
}

func NewError(code int, msg string) Error {
	return Error{
		Code:    code,
		Message: msg,
	}
}

var (
	// ErrGetUsers is used as body response when is not possible to get the users
	// list.
	ErrGetUsers = NewError(100, "Can't get the users")
)
