// Copyright 2018 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the WTFPL.

package usersd

import (
	"time"

	"github.com/gbrlsnchs/jwt/v2"
)

// Token is a JWT token with public claims. See also:
//
// * Tx.JWT: ready to use JWT generation.
// * Service.VerifyJWT: JWT verification.
type Token struct {
	*jwt.JWT
	User *User `json:"user"`
}

// UnmarshalJWT parses the given data into the JWT. Notice that UnmarshalJWT
// doesn't verify the JWT signature.
func UnmarshalJWT(data []byte) (*Token, error) {
	payload, _, err := jwt.ParseBytes(data)
	if err != nil {
		return nil, err
	}

	jot := new(Token)
	if err = jwt.Unmarshal(payload, &jot); err != nil {
		return nil, err
	}

	return jot, nil
}

// JWT generates a JWT for the given user. The JWT can't be used before
// notBefore of after expire; for no limits use 0.
func (tx *Tx) JWT(userid string, notBefore, expire int64) ([]byte, error) {
	user, err := GetUser(tx, userid)
	if err != nil {
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

		User: user,
	}

	hs256 := jwt.NewHS256(tx.Service.opts.JWTSecret)
	jot.SetAlgorithm(hs256)
	payload, err := jwt.Marshal(jot)
	if err != nil {
		return nil, err
	}

	return hs256.Sign(payload)
}

// VerifyJWT returns true if the JWT was generated and signed by the service.
func (tx *Tx) VerifyJWT(token []byte) bool {
	secret := tx.Service.opts.JWTSecret
	payload, sig, err := jwt.ParseBytes(token)
	if err != nil {
		return false
	}

	hs256 := jwt.NewHS256(secret)
	err = hs256.Verify(payload, sig)
	return err == nil
}
