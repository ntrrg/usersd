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

// NewUser creates a new user.
func NewUser(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrInternal)
		return
	}

	user, err := usersd.NewUserJSON(data)

	if err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrUnmarshalUser)
		return
	}

	if user.Save(); err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrInternal)
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
		log.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrGetUsers)
		return
	}

	json.NewEncoder(w).Encode(users)
}

// GetUser gets a user.
func GetUser(w http.ResponseWriter, r *http.Request) {
	id := vestigo.Param(r, "id")
	user, err := usersd.GetUser(id)

	if err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrGetUsers)
		return
	}

	json.NewEncoder(w).Encode(user)
}

// UpdateUser creates a new user.
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := vestigo.Param(r, "id")
	data, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrInternal)
		return
	}

	user, err := usersd.NewUserJSON(data)

	if err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrUnmarshalUser)
		return
	}

	user.ID = id
	user.Set("id", id)

	if user.Save(); err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrInternal)
		return
	}

	json.NewEncoder(w).Encode(user)
}

// DeleteUser gets a user.
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := vestigo.Param(r, "id")
	user := new(usersd.User)
	user.ID = id

	if err := user.Delete(); err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrGetUsers)
		return
	}
}

// Reset resets the database.
func Reset(w http.ResponseWriter, r *http.Request) {
	if err := usersd.Reset(); err != nil {
		log.Println(err)
		http.Error(w, "Can't reset the database", http.StatusInternalServerError)
		return
	}

	if _, err := w.Write([]byte("Done.")); err != nil {
		log.Println(err)
	}
}
