// Copyright 2018 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the MIT license.

package middleware

import (
	"net/http"
	"strings"
)

// ReplaceURL replaces old by new from the request URL.
func ReplaceURL(old, news string) Adapter {
	return func(h http.Handler) http.Handler {
		nh := func(w http.ResponseWriter, r *http.Request) {
			r.URL.Path = strings.Replace(r.URL.Path, old, news, 1)
			h.ServeHTTP(w, r)
		}

		return http.HandlerFunc(nh)
	}
}

// StripPrefix strips the given prefix from the request URL.
func StripPrefix(prefix string) Adapter {
	return func(h http.Handler) http.Handler {
		return http.StripPrefix(prefix, h)
	}
}
