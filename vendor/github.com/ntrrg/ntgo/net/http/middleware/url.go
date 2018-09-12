// Copyright 2018 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the MIT license.

package middleware

import "net/http"

// StripPrefix strips the given prefix from the request URL.
func StripPrefix(prefix string) Adapter {
	return func(h http.Handler) http.Handler {
		return http.StripPrefix(prefix, h)
	}
}
