// Copyright 2018 Miguel Angel Rivera Notararigo. All rights reserved.
// This source code was released under the WTFPL.

package usersd

import (
	"crypto/rand"
	"encoding/base64"
	"net/smtp"
)

// Email verification code document identifier.
const EmailVerificationCodeDI = "email-verification-codes"

// GetEmailVerificationCode generates a special code for email verification.
func (tx *Tx) GetEmailVerificationCode(userid string) error {
	user, err := tx.GetUser(userid)
	if err != nil {
		return err
	}

	if user.Email == "" {
		return ErrUserEmailEmpty
	}

	hash := make([]byte, 9)
	if _, err = rand.Read(hash); err != nil {
		return err
	}

	code := base64.URLEncoding.EncodeToString(hash)
	ttl := usersd.opts.EmailVerificationCodesTTL
	err = tx.SetWithTTL([]byte(EmailVerificationCodeDI+user.ID), code, ttl)
	if err != nil {
		return err
	}

	auth := smtp.PlainAuth("", opts.Email, opts.Password, opts.Host)
	msg := []byte(
		"To: " + email + "\r\n" +
			"Subject: MI_RI password recovery\r\n" +
			"\r\n" +
			"Go to " + link + ", you will be asked for access with your new " +
			"random password (" + password + ")\r\n",
	)

	err = smtp.SendMail(
		SMTPHost+":"+SMTPPort,
		auth,
		SMTPEmail,
		[]string{email},
		msg,
	)
}

// VerifyEmail.
func (tx *Tx) VerifyEmail(user *User) ([]byte, error) {
	if err := tx.ValidateUser(user); err != nil {
		return nil, err
	}

	email := user.Email
	if email == "" {
		return nil, ErrUserEmailEmpty
	}

	hash := make([]byte, opts.SaltSize)
	if _, err := rand.Read(hash); err != nil {
		return err
	}

	tx, err := newTxHTTP(w, r)
	if err != nil {
		return
	}

	defer txRollback(tx)

	var id string
	q := "SELECT id FROM users WHERE email=?"
	err = txQueryRow(tx, q, email).Scan(&id)
	if err != nil {
		http.Error(w, "", http.StatusNotFound)
		return
	}

	b := make([]byte, 33)
	rand.Read(b)
	password := base64.URLEncoding.EncodeToString(b)
	key, err := NewRecoverLinkKey(email + ":" + password)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	redirect := r.FormValue("redirect")
	link := fmt.Sprintf(
		"%v%v/recover/?key=%v&redirect=%v",
		BaseURL, v[1:], key, redirect,
	)

	lDebug("[DEBUG][SMTP] Sending recovery link to %v -> %v", email, link)
	auth := smtp.PlainAuth("", SMTPEmail, SMTPPassword, SMTPHost)
	msg := []byte(
		"To: " + email + "\r\n" +
			"Subject: MI_RI password recovery\r\n" +
			"\r\n" +
			"Go to " + link + ", you will be asked for access with your new " +
			"random password (" + password + ")\r\n",
	)

	err = smtp.SendMail(
		SMTPHost+":"+SMTPPort,
		auth,
		SMTPEmail,
		[]string{email},
		msg,
	)

	if err != nil {
		log.Printf("[ERROR][SMTP] Can't send the email -> %v", err)
		http.Error(w, "", http.StatusInternalServerError)
	}
}

