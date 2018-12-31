// Copyright 2018 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the WTFPL.

package usersd

import (
	"time"

	"github.com/gbrlsnchs/jwt/v2"
)

// Token is a JWT token with public claims.
type Token struct {
	*jwt.JWT
	User User `json:"user"`
}

// UnmarshalJWT parses the JWT and returns a Token if possible. Notice that
// UnmarshalJWT doesn't verify the JWT signature.
func UnmarshalJWT(token []byte) (*Token, error) {
	payload, _, err := jwt.ParseBytes(token)
	if err != nil {
		return nil, err
	}

	jot := new(Token)
	if err = jwt.Unmarshal(payload, jot); err != nil {
		return nil, err
	}

	return jot, nil
}

// JWT generates a JWT for the given user. The JWT can't be used before
// notBefore of after expire, for no limits use 0.
func (s *Service) JWT(user *User, notBefore, expire int64) ([]byte, error) {
	tx := s.DB.NewTransaction(false)
	defer tx.Discard()
	index := s.Index["users"]

	if err := user.Validate(tx, index); err != nil {
		return nil, err
	}

	jot := &Token{
		JWT: &jwt.JWT{
			Issuer:  "usersd",
			Subject: user.ID,

			ExpirationTime: expire,
			NotBefore:      notBefore,
			IssuedAt:       time.Now().Unix(),
		},

		User: *user,
	}

	hs256 := jwt.NewHS256(s.opts.JWTSecret)
	jot.SetAlgorithm(hs256)
	jot.User.Password = ""
	payload, err := jwt.Marshal(jot)
	if err != nil {
		return nil, err
	}

	return hs256.Sign(payload)
}

// VerifyJWT returns true if the JWT was signed by the service.
func (s *Service) VerifyJWT(token []byte) bool {
	return VerifyJWT(s.opts.JWTSecret, token)
}

// VerifyJWT returns true if the JWT was signed with the given secret.
func VerifyJWT(secret string, token []byte) bool {
	payload, sig, err := jwt.ParseBytes(token)
	if err != nil {
		return false
	}

	hs256 := jwt.NewHS256(secret)
	err = hs256.Verify(payload, sig)
	return err == nil
}
