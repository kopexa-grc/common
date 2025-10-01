// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package tokens

import (
	"time"

	"github.com/vmihailenco/msgpack/v5"
)

// NewVerificationToken creates a token struct from an email address that expires
// in expirationDays (default: 7) days.
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
// for email ownership verification.
type VerificationToken struct {
	Email string `msgpack:"email"`
	SigningInfo
}

// Sign creates a base64 URL encoded signature for the token's msgpack representation.
// The returned secret MUST be stored securely; without it the signature cannot be
// recomputed for verification. The secret concatenates nonce||key.
func (t *VerificationToken) Sign() (string, []byte, error) {
	data, err := msgpack.Marshal(t)
	if err != nil {
		return "", nil, err
	}

	return t.signData(data)
}

// Verify checks that a token was signed with the secret, required fields are present,
// and it has not expired.
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

// ResetToken packages a user ID with random data and an expiration time so that it can
// be serialized and hashed into a token which can be sent to users for password resets.
type ResetToken struct {
	UserID string `msgpack:"user_id"`
	SigningInfo
}

// NewResetToken creates a token struct from a user ID that expires in resetTokenExpirationMinutes.
func NewResetToken(id string) (token *ResetToken, err error) {
	if id == "" {
		return nil, ErrMissingUserID
	}

	token = &ResetToken{
		UserID: id,
	}

	if token.SigningInfo, err = NewSigningInfo(time.Minute * resetTokenExpirationMinutes); err != nil {
		return nil, err
	}

	return token, nil
}

// Sign creates a base64 URL encoded signature for the reset token. See VerificationToken.Sign.
func (t *ResetToken) Sign() (string, []byte, error) {
	return t.SignToken(t)
}

// Validate checks that the token has required fields (UserID must be non-empty).
func (t *ResetToken) Validate() error {
	if t.UserID == "" {
		return ErrTokenMissingUserID
	}

	return nil
}

// SetNonce sets the nonce for verification (implements URLToken contract).
func (t *ResetToken) SetNonce(nonce []byte) {
	t.Nonce = nonce
}

// Verify performs full validation (required fields, expiration, signature) for a ResetToken.
func (t *ResetToken) Verify(signature string, secret []byte) error {
	if err := t.Validate(); err != nil {
		return err
	}

	return t.VerifyToken(t, signature, secret)
}
