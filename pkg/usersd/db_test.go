// Copyright 2018 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the WTFPL.

package usersd_test

import (
	"fmt"

	"github.com/ntrrg/usersd/pkg/usersd"
)

func ExampleInit_newLocation() {
	if err := usersd.Init("test-db"); err != nil {
		fmt.Println(err)
	}

	defer usersd.Close()

	// Output:
}

func ExampleInit_temporaryStorage() {
	if err := usersd.Init(""); err != nil {
		fmt.Println(err)
	}

	defer usersd.Close()

	// Output:
}
