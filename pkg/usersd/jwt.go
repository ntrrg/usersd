// Copyright 2018 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the WTFPL.

package usersd

import (
	"time"

	"github.com/gbrlsnchs/jwt/v2"
)

type JWTOptions struct {
	// Issuer claim for JWTs.
	Issuer string

	// Secret for signing and verifying JWTs.
	Secret string
}

// Token is a JWT token with public claims. See also:
//
// * Tx.JWT: create ready to use JWT generation.
//
// * Tx.VerifyJWT: verify and validate JWT.
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

// JWT generates a JWT for the given user. The JWT can't be used after expire,
// for no limit use 0.
func (tx *Tx) JWT(userid string, expire int64) ([]byte, error) {
	user, err := tx.GetUser(userid)
	if err != nil {
		return nil, err
	}

	jot := &Token{
		JWT: &jwt.JWT{
			Issuer:  tx.Service.opts.JWTOpts.Issuer,
			Subject: user.ID,

			ExpirationTime: expire,
			IssuedAt:       time.Now().Unix(),
		},

		User: user,
	}

	hs256 := jwt.NewHS256(tx.Service.opts.JWTOpts.Secret)
	jot.SetAlgorithm(hs256)
	payload, err := jwt.Marshal(jot)
	if err != nil {
		return nil, err
	}

	return hs256.Sign(payload)
}

// VerifyJWT returns true if the JWT was generated and signed by the service.
func (tx *Tx) VerifyJWT(token []byte) bool {
	payload, sig, err := jwt.ParseBytes(token)
	if err != nil {
		return false
	}

	hs256 := jwt.NewHS256(tx.Service.opts.JWTOpts.Secret)
	if err = hs256.Verify(payload, sig); err != nil {
		return false
	}

	jot := &Token{}
	if err = jwt.Unmarshal(payload, jot); err != nil {
		return false
	}

	validators := []jwt.ValidatorFunc{
		jwt.IssuerValidator(tx.Service.opts.JWTOpts.Issuer),
	}

	if jot.ExpirationTime != 0 {
		validators = append(validators, jwt.ExpirationTimeValidator(time.Now()))
	}

	err = jot.Validate(validators...)
	return err == nil
}
