// Copyright 2018 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the WTFPL.

package rest

import (
	"net/http"

	"github.com/husobee/vestigo"
	"github.com/ntrrg/ntgo/net/http/middleware"
)

// Mux is the main mux.
func Mux() http.Handler {
	mux := vestigo.NewRouter()

	mux.HandleFunc("/healthz", Healthz)
	mux.HandleFunc("/reset", Reset)

	// Users

	mux.Post("/v1/users", middleware.AdaptFunc(
		NewUser,
		middleware.JSONRequest(""),
	).ServeHTTP)

	mux.Get("/v1/users", middleware.AdaptFunc(
		ListUsers,
		middleware.Cache("max-age=3600, s-max-age=3600"),
		middleware.Gzip(-1),
	).ServeHTTP)

	mux.Get("/v1/users/:id", middleware.AdaptFunc(
		GetUser,
		middleware.Cache("max-age=3600, s-max-age=3600"),
		middleware.Gzip(-1),
	).ServeHTTP)

	mux.Put("/v1/users/:id", middleware.AdaptFunc(
		UpdateUser,
		middleware.JSONRequest(""),
	).ServeHTTP)

	mux.Delete("/v1/users/:id", DeleteUser)

	return middleware.Adapt(
		mux,
		middleware.JSONResponse(),
	)
}
