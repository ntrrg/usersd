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

	// Create user
	user := &usersd.User{Email: string(data)}

	w.Header().Set("Location", "/v1/users/"+user.ID)
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Printf("[ERROR][REST] Can't write the response -> %v", err)
	}
}

// ListUsers list all the users.
func ListUsers(w http.ResponseWriter, r *http.Request) {
	// Get users
	users := []*usersd.User{}

	if err := json.NewEncoder(w).Encode(users); err != nil {
		log.Printf("[ERROR][REST] Can't write the response -> %v", err)
	}
}

// GetUser gets a user.
func GetUser(w http.ResponseWriter, r *http.Request) {
	id := vestigo.Param(r, "id")

	// Get user
	user := &usersd.User{ID: id}

	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Printf("[ERROR][REST] Can't write the response -> %v", err)
	}
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

	// Update user
	user := &usersd.User{ID: id, Email: string(data)}

	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Printf("[ERROR][REST] Can't write the response -> %v", err)
	}
}

// DeleteUser gets a user.
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := vestigo.Param(r, "id")
	// Delete user
	_ = id
}
