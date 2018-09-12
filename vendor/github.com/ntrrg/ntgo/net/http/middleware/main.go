// Copyright 2018 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the MIT license.

package middleware

import "net/http"

// Adapter is a wrapper for http.Handler. Takes a http.Handler as argument and
// creates a new one that may run code before and/or after calling the given.
type Adapter func(http.Handler) http.Handler

// Adapt wraps a http.Handler into a list of Adapters. Takes a http.Handler
// and a list of Adapters, which will be wrapped from right to left (and ran
// left to right).
func Adapt(h http.Handler, a ...Adapter) http.Handler {
	for i := len(a) - 1; i >= 0; i-- {
		h = a[i](h)
	}

	return h
}

// AdaptFunc works as Adapt but for http.HandlerFunc.
func AdaptFunc(h http.HandlerFunc, a ...Adapter) http.Handler {
	return Adapt(h, a...)
}
