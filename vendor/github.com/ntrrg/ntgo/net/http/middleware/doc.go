// Copyright 2018 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the MIT license.

/*
Package middleware provides flexibility at the HTTP request/response process.
For this purpose, the Adapter pattern is used, which consist in wrapping
handlers. An adapter may run code before and/or after the handler it is
wrapping.

	func MyAdapter(h http.Handler) http.Handler {
		nh := func(w http.ResponseWriter, r *http.Request) {
			// Code that run before
			h.ServeHTTP(w, r)
			// Code that run after
		}

		return http.HandlerFunc(nh)
	}

When a set of adapters is used, the result is a sequence of wrappers and its
execution flow depends in the order that the adapters were given.

	Adapt(h, f1, f2, f3)

1. f1 before code

2. f2 before code

3. f3 before code

4. h

5. f3 after code

6. f2 after code

7. f1 after code
*/
package middleware
