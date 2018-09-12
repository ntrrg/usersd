// Copyright 2018 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the WTFPL.

package rest

import (
	nthttp "github.com/ntrrg/ntgo/net/http"
)

var Server *nthttp.Server

func init() {
	Server = nthttp.NewServer(nthttp.Config{
		Handler: Mux(),
	})
}
