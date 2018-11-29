// Copyright 2018 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the MIT license.

package middleware

import (
	"net/http"
	"net/url"
	"strings"
)

// Replace replaces old by new from the request URL.
func Replace(old, New string) Adapter {
	return func(h http.Handler) http.Handler {
		nh := func(w http.ResponseWriter, r *http.Request) {
			r2 := new(http.Request)
			*r2 = *r
			r2.URL = new(url.URL)
			*r2.URL = *r.URL
			r2.URL.Path = strings.Replace(r.URL.Path, old, New, 1)
			h.ServeHTTP(w, r2)
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
