// Copyright 2018 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the WTFPL.

package rest

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/husobee/vestigo"
	"github.com/ntrrg/usersd/pkg/usersd"
)

// Healthz is the handler for the healtz endpoint.
func Healthz(w http.ResponseWriter, r *http.Request) {
}

// Reset resets the database.
func Reset(w http.ResponseWriter, r *http.Request) {
	if err := usersd.Reset(); err != nil {
		ErrInternal.ServeHTTP(w, r)
	}
}

// Users

// NewUser creates a new user.
func NewUser(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Printf("[ERROR][REST] Can't read the request body -> %v", err)
		ErrInternal.ServeHTTP(w, r)
		return
	}

	user, err := usersd.CreateUserJSON(data)

	if err != nil {
		ErrInternal.ServeHTTP(w, r)
		return
	}

	w.Header().Set("Location", "/v1/users/"+user.ID)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// ListUsers list all the users.
func ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := usersd.ListUsers()

	if err != nil {
		ErrInternal.ServeHTTP(w, r)
		return
	}

	json.NewEncoder(w).Encode(users)
}

// GetUser gets a user.
func GetUser(w http.ResponseWriter, r *http.Request) {
	id := vestigo.Param(r, "id")
	user, err := usersd.GetUser(id)

	if err != nil {
		ErrInternal.ServeHTTP(w, r)
		return
	}

	json.NewEncoder(w).Encode(user)
}

// UpdateUser creates a new user.
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := vestigo.Param(r, "id")
	data, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Printf("[ERROR][REST] Can't read the request body -> %v", err)
		ErrInternal.ServeHTTP(w, r)
		return
	}

	user, err := usersd.NewUserJSON(data)

	if err != nil {
		ErrCantUnmarshalUser.ServeHTTP(w, r)
		return
	}

	user.ID = id

	if user.Update(); err != nil {
		ErrInternal.ServeHTTP(w, r)
		return
	}

	json.NewEncoder(w).Encode(user)
}

// DeleteUser gets a user.
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := vestigo.Param(r, "id")

	if err := usersd.NewUser(id, nil).Delete(); err != nil {
		ErrInternal.ServeHTTP(w, r)
	}
}
