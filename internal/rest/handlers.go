// Copyright 2018 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the WTFPL.

package rest

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ntrrg/usersd/pkg/usersd"
)

func Users(w http.ResponseWriter, r *http.Request) {
	users, err := usersd.GetUsers()

	if err != nil {
		log.Println(err)
		http.Error(w, "", 500)
		json.NewEncoder(w).Encode(ErrGetUsers)
		return
	}

	result := map[string]interface{}{"users": users}
	json.NewEncoder(w).Encode(result)
}

func Reset(w http.ResponseWriter, r *http.Request) {
	if err := usersd.Reset(); err != nil {
		log.Println(err)
		http.Error(w, "Can't reset the database", 500)
		return
	}

	if _, err := w.Write([]byte("Done.")); err != nil {
		log.Println(err)
	}
}
