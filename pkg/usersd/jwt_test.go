// Copyright 2018 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the WTFPL.

package usersd_test

import (
	"testing"

	"github.com/ntrrg/usersd/pkg/usersd"
)

func TestUnmarshalJWT(t *testing.T) {
	token := []byte("invalid jwt format")

	if _, err := usersd.UnmarshalJWT(token); err == nil {
		t.Error("Invalid JWT parsed")
	}

	token = []byte("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJ1c2Vyc2QiLCJzdWIiOiJ0ZXN0IiwiaWF0IjoxNTQ2MjI1MTk0LCJ1c2VyIjp7ImlkIjoxMjM0LCJlbWFpbCI6InRlc3RAZXhhbXBsZS5jb20iLCJtb2RlIjoibG9jYWwiLCJ2ZXJpZmllZCI6ZmFsc2UsImNyZWF0ZWRBdCI6MTU0NjIyNTE5NCwibGFzdExvZ2luIjowfX0.OQrbnjdYk9glBP9i5OWhAdReOh_8i8zd5JJtcnOrfL0") // nolint: lll

	if _, err := usersd.UnmarshalJWT(token); err == nil {
		t.Error("Invalid JWT unmarshaled")
	}

	token = []byte("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJ1c2Vyc2QiLCJzdWIiOiJ0ZXN0IiwiaWF0IjoxNTQ2MjI1MTk0LCJ1c2VyIjp7ImlkIjoidGVzdCIsImVtYWlsIjoidGVzdEBleGFtcGxlLmNvbSIsIm1vZGUiOiJsb2NhbCIsInZlcmlmaWVkIjpmYWxzZSwiY3JlYXRlZEF0IjoxNTQ2MjI1MTk0LCJsYXN0TG9naW4iOjB9fQ.CE6an7tDnzsEsq2aexjln5uUuG5Rtju6ObDqgbTLDro") // nolint: lll

	if _, err := usersd.UnmarshalJWT(token); err != nil {
		t.Errorf("Can't unmarshal valid JWT -> %v", err)
	}
}

func TestTx_JWT(t *testing.T) {
	if err := initTest("tx-jwt", true); err != nil {
		t.Fatal(err)
	}

	defer usersd.Close()

	tx := usersd.NewTx(true)
	defer tx.Discard()

	if _, err := tx.JWT("non-existent-user", 0); err == nil {
		t.Error("JWT generated for non existent user")
	}

	if _, err := tx.JWT("admin", 0); err != nil {
		t.Errorf("Can't generate the JWT -> %v", err)
	}
}

func TestTx_VerifyJWT(t *testing.T) {
	if err := initTest("tx-verify-jwt", true); err != nil {
		t.Fatal(err)
	}

	defer usersd.Close()

	tx := usersd.NewTx(false)
	defer tx.Discard()

	token := []byte("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJ1c2Vyc2QiLCJzdWIiOiJ0ZXN0IiwiaWF0IjoxNTQ2MjI1MTk0LCJ1c2VyIjp7ImlkIjoidGVzdCIsImVtYWlsIjoidGVzdEBleGFtcGxlLmNvbSIsIm1vZGUiOiJsb2NhbCIsInZlcmlmaWVkIjpmYWxzZSwiY3JlYXRlZEF0IjoxNTQ2MjI1MTk0LCJsYXN0TG9naW4iOjB9fQ.CE6an7tDnzsEsq2aexjln5uUuG5Rtju6ObDqgbTLDro") // nolint: lll

	if !tx.VerifyJWT(token) {
		t.Error("Can't verify valid JWT")
	}
}
