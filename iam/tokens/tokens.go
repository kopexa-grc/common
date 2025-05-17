// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package tokens

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"time"

	"github.com/vmihailenco/msgpack/v5"
)

// SigningInfo contains the cryptographic information needed to sign and verify tokens.
// It includes an expiration time and a nonce for additional security.
type SigningInfo struct {
	// ExpiresAt is the UTC timestamp when the token expires.
	ExpiresAt time.Time `msgpack:"expires_at"`
	// Nonce is a random value used to prevent token reuse.
	Nonce []byte `msgpack:"nonce"`
}

// NewSigningInfo creates a new SigningInfo instance with the specified expiration duration.
// It generates a random nonce and sets the expiration time.
//
// Parameters:
//   - expires: The duration until the token expires. Must be greater than 0.
//
// Returns:
//   - SigningInfo: The created signing info
//   - error: If the expiration is 0 or nonce generation fails
func NewSigningInfo(expires time.Duration) (SigningInfo, error) {
	if expires == 0 {
		return SigningInfo{}, ErrExpirationIsRequired
	}

	info := SigningInfo{
		ExpiresAt: time.Now().UTC().Add(expires).Truncate(time.Microsecond),
		Nonce:     make([]byte, nonceLength),
	}

	if _, err := rand.Read(info.Nonce); err != nil {
		return info, ErrFailedSigning.With(err)
	}

	return info, nil
}

// IsExpired checks if the token has passed its expiration time.
//
// Returns:
//   - bool: true if the token is expired, false otherwise
func (d SigningInfo) IsExpired() bool {
	return d.ExpiresAt.Before(time.Now())
}

// signData signs the provided data using HMAC-SHA256 and returns the signature and secret.
//
// Parameters:
//   - data: The data to sign
//
// Returns:
//   - string: The base64-encoded signature
//   - []byte: The secret containing the nonce and key
//   - error: If signing fails
func (d SigningInfo) signData(data []byte) (string, []byte, error) {
	key := make([]byte, keyLength)
	if _, err := rand.Read(key); err != nil {
		return "", nil, ErrFailedSigning.With(err)
	}

	mac := hmac.New(sha256.New, key)
	if _, err := mac.Write(data); err != nil {
		return "", nil, ErrFailedSigning.With(err)
	}

	secret := make([]byte, nonceLength+keyLength)
	copy(secret[:nonceLength], d.Nonce)
	copy(secret[nonceLength:], key)

	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil)), secret, nil
}

// verifyData verifies the signature of the provided data using the secret.
//
// Parameters:
//   - data: The data to verify
//   - signature: The base64-encoded signature to verify against
//   - secret: The secret containing the nonce and key
//
// Returns:
//   - error: If verification fails
func (d SigningInfo) verifyData(data []byte, signature string, secret []byte) error {
	var err error

	mac := hmac.New(sha256.New, secret[nonceLength:])
	if _, err = mac.Write(data); err != nil {
		return err
	}

	var token []byte

	if token, err = base64.RawURLEncoding.DecodeString(signature); err != nil {
		return err
	}

	if !hmac.Equal(mac.Sum(nil), token) {
		return ErrTokenInvalid
	}

	return nil
}

// OrganizationInviteToken represents a token used for inviting users to an organization.
// It contains the email address of the invitee and the organization ID.
type OrganizationInviteToken struct {
	// Email is the email address of the user being invited.
	Email string `msgpack:"email"`
	// OrganizationID is the ID of the organization the user is being invited to.
	OrganizationID string `msgpack:"organization_id"`
	// SigningInfo contains the cryptographic information for the token.
	SigningInfo
}

// NewOrganizationInviteToken creates a new organization invite token.
//
// Parameters:
//   - email: The email address of the user being invited
//   - organizationID: The ID of the organization
//
// Returns:
//   - *OrganizationInviteToken: The created token
//   - error: If token creation fails
func NewOrganizationInviteToken(email string, organizationID string) (*OrganizationInviteToken, error) {
	var err error

	if email == "" {
		return nil, ErrInviteTokenMissingEmail
	}

	token := &OrganizationInviteToken{
		Email:          email,
		OrganizationID: organizationID,
	}

	if token.SigningInfo, err = NewSigningInfo(time.Hour * 24 * inviteExpirationDays); err != nil {
		return nil, err
	}

	return token, nil
}

// Sign signs the token and returns its signature and secret.
//
// Returns:
//   - string: The base64-encoded signature
//   - []byte: The secret containing the nonce and key
//   - error: If signing fails
func (t *OrganizationInviteToken) Sign() (string, []byte, error) {
	data, err := msgpack.Marshal(t)
	if err != nil {
		return "", nil, err
	}

	return t.signData(data)
}

// Verify verifies the token's signature and checks its validity.
//
// Parameters:
//   - signature: The base64-encoded signature to verify
//   - secret: The secret containing the nonce and key
//
// Returns:
//   - error: If verification fails
func (t *OrganizationInviteToken) Verify(signature string, secret []byte) error {
	if t.Email == "" {
		return ErrInviteTokenMissingEmail
	}

	if t.IsExpired() {
		return ErrTokenExpired
	}

	if len(secret) != nonceLength+keyLength {
		return ErrInvalidSecret
	}

	t.Nonce = secret[0:nonceLength]

	var data []byte

	var err error

	if data, err = msgpack.Marshal(t); err != nil {
		return err
	}

	return t.verifyData(data, signature, secret)
}
