// Copyright 2018 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the MIT license.

package middleware

import "net/http"

// AddHeader creates/appends a HTTP header before calling the http.Handler.
func AddHeader(key, value string) Adapter {
	return func(h http.Handler) http.Handler {
		nh := func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add(key, value)
			h.ServeHTTP(w, r)
		}

		return http.HandlerFunc(nh)
	}
}

// DelHeader removes a HTTP header before calling the http.Handler.
func DelHeader(key string) Adapter {
	return func(h http.Handler) http.Handler {
		nh := func(w http.ResponseWriter, r *http.Request) {
			w.Header().Del(key)
			h.ServeHTTP(w, r)
		}

		return http.HandlerFunc(nh)
	}
}

// SetHeader creates/replaces a HTTP header before calling the http.Handler.
func SetHeader(key, value string) Adapter {
	return func(h http.Handler) http.Handler {
		nh := func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set(key, value)
			h.ServeHTTP(w, r)
		}

		return http.HandlerFunc(nh)
	}
}

// Cache sets HTTP cache headers for GET requests.
//
// Directives
//
// * public/private: whether the cached response is for any or a specific user.
//
// * max-age=TIME: cache life time in seconds. The maximum value is 1 year.
//
// * s-max-age=TIME: same as max-age, but this one has effect in proxies.
//
// * must-revalidate: force expired cached response revalidation, even in
// special circumstances (like slow connections, were cached responses are used
// even after they had expired).
//
// * proxy-revalidate: same as must-revalidate, but this one has effect in
// proxies.
//
// * no-cache: disables cache.
//
// * no-store: disables cache, even in proxies.
func Cache(directives string) Adapter {
	return func(h http.Handler) http.Handler {
		nh := func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet {
				w.Header().Set("Cache-Control", directives)
			}

			h.ServeHTTP(w, r)
		}

		return http.HandlerFunc(nh)
	}
}

// JSONResponse prepares the response to be a JSON response.
func JSONResponse() Adapter {
	return SetHeader("Content-Type", "application/json; charset=utf-8")
}
