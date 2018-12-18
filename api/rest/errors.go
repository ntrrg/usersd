// Copyright 2018 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the WTFPL.

package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Error represents an error during the execution of any handler. It could be
// just one error or a collection of them.
type Error struct {
	Code    int     `json:"code"`
	Message string  `json:"message"`
	HTTP    int     `json:"-"`
	Errors  []Error `json:"errors,omitempty"`
}

func (e Error) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var body string

	if data, err := json.Marshal(e); err != nil {
		body = `{"code": 0, "message": "Can't parse the error"}`
	} else {
		body = string(data)
	}

	http.Error(w, body, e.HTTP)
}

func (e Error) Error() string {
	return fmt.Sprintf("(%d) %s", e.Code, e.Message)
}

// Internal errors
var (
	ErrInternal = Error{
		Code:    0,
		Message: "Internal Server Error",
		HTTP:    http.StatusInternalServerError,
	}
)

// Users errors (1XX codes).
var (
	ErrGetUsers = Error{
		Code:    100,
		Message: "Can't get the users list",
	}

	ErrCantUnmarshalUser = Error{
		Code:    101,
		Message: "Can't unmarshal the user from the request body",
		HTTP:    http.StatusBadRequest,
	}
)
