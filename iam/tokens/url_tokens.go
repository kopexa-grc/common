// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package tokens

import (
	"time"

	"github.com/vmihailenco/msgpack/v5"
)

// NewVerificationToken creates a token struct from an email address that expires
// in 7 days
func NewVerificationToken(email string) (token *VerificationToken, err error) {
	if email == "" {
		return nil, ErrMissingEmail
	}

	token = &VerificationToken{
		Email: email,
	}

	if token.SigningInfo, err = NewSigningInfo(time.Hour * 24 * expirationDays); err != nil {
		return nil, err
	}

	return token, nil
}

// VerificationToken packages an email address with random data and an expiration
// time so that it can be serialized and hashed into a token which can be sent to users
type VerificationToken struct {
	Email string `msgpack:"email"`
	SigningInfo
}

// Sign creates a base64 encoded string from the token data so that it can be sent to
// users as part of a URL. The returned secret should be stored in the database so that
// the string can be recomputed when verifying a user provided token.
func (t *VerificationToken) Sign() (string, []byte, error) {
	data, err := msgpack.Marshal(t)
	if err != nil {
		return "", nil, err
	}

	return t.signData(data)
}

// Verify checks that a token was signed with the secret and is not expired
func (t *VerificationToken) Verify(signature string, secret []byte) (err error) {
	if t.Email == "" {
		return ErrTokenMissingEmail
	}

	if t.IsExpired() {
		return ErrTokenExpired
	}

	if len(secret) != nonceLength+keyLength {
		return ErrInvalidSecret
	}

	// Serialize the struct with the nonce from the secret
	t.Nonce = secret[0:nonceLength]

	var data []byte

	if data, err = msgpack.Marshal(t); err != nil {
		return err
	}

	return t.verifyData(data, signature, secret)
}
