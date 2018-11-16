// Copyright 2018 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the WTFPL.

package usersd_test

import (
	"github.com/ntrrg/usersd/pkg/usersd"
)

var Opts = usersd.DefaultOptions

type userData struct {
	id, password string

	data map[string]interface{}
}

type userCase struct {
	in, want userData
}
