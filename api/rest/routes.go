// Copyright 2018 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the WTFPL.

package rest

import (
	"net/http"

	"github.com/ntrrg/ntgo/net/http/middleware"
)

func Mux() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/reset", Reset)

	uh := UsersMux()
	mux.Handle("/user/", middleware.Adapt(uh, middleware.StripPrefix("/user")))
	mux.Handle("/users/", middleware.Adapt(uh, middleware.StripPrefix("/users")))

	return mux
}

func UsersMux() *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/", middleware.AdaptFunc(
		Users,
		middleware.JSONResponse(),
		// middleware.Cache("max-age=3600, s-max-age=3600"),
		// middleware.Gzip(-1),
	))

	return mux
}
