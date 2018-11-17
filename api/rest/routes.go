// Copyright 2018 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the WTFPL.

package rest

import (
	"net/http"

	"github.com/ntrrg/ntgo/net/http/middleware"
)

func Mux() *http.ServeMux {
	v := "/v1"
	mux := http.NewServeMux()

	mux.HandleFunc("/reset", Reset)

	mux.Handle(
		v+"/users/",
		middleware.Adapt(UsersMux(), middleware.StripPrefix(v+"/users")),
	)

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
