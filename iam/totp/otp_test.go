// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package totp

import (
	"context"
	"database/sql"
	"encoding/base32"
	"encoding/base64"
	"testing"
	"time"

	"github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockStore struct {
	values map[string][]byte
}

func newMockStore() *mockStore {
	return &mockStore{
		values: make(map[string][]byte),
	}
}

func (s *mockStore) Get(_ context.Context, key string) ([]byte, error) {
	return s.values[key], nil
}

func (s *mockStore) Set(_ context.Context, key string, value []byte) error {
	s.values[key] = value
	return nil
}

func (s *mockStore) Delete(_ context.Context, key string) error {
	delete(s.values, key)
	return nil
}

func TestOTP(t *testing.T) {
	store := newMockStore()
	otp := New(store,
		WithCodeLength(6),
		WithIssuer("test"),
		WithRecoveryCodeCount(16),
		WithRecoveryCodeLength(8),
		WithSecret(Secret{
			Version: 1,
			Key:     []byte("test-secret-key-1234567890123456"),
		}),
	)

	// Test TOTPSecret
	t.Run("TOTPSecret", func(t *testing.T) {
		user := &User{
			ID:    "test-user",
			Email: sql.NullString{String: "test@example.com", Valid: true},
		}

		secret, err := otp.TOTPSecret(user)
		require.NoError(t, err)
		assert.NotEmpty(t, secret)

		// Test decryption
		decrypted, err := otp.TOTPDecryptedSecret(secret)
		require.NoError(t, err)
		assert.NotEmpty(t, decrypted)
	})

	// Test OTPCode
	t.Run("OTPCode", func(t *testing.T) {
		address := "test@example.com"
		method := Email

		code, hash, err := otp.OTPCode(address, method)
		require.NoError(t, err)
		assert.NotEmpty(t, code)
		assert.NotEmpty(t, hash)

		// Test validation
		err = otp.ValidateOTP(code, hash)
		require.NoError(t, err)

		// Test invalid code
		err = otp.ValidateOTP("invalid", hash)
		assert.Error(t, err)
	})

	// Test ValidateTOTP
	t.Run("ValidateTOTP", func(t *testing.T) {
		user := &User{
			ID:    "test-user",
			Email: sql.NullString{String: "test@example.com", Valid: true},
		}

		secret, err := otp.TOTPSecret(user)
		require.NoError(t, err)

		user.TOTPSecret = secret

		// Entschlüssle das Secret und dekodiere es als base32
		decrypted, err := otp.TOTPDecryptedSecret(secret)
		require.NoError(t, err)
		decodedSecret, err := base32.StdEncoding.DecodeString(decrypted)
		require.NoError(t, err)
		assert.Len(t, decodedSecret, Base32SecretLength)

		// Generate a valid code
		code, err := totp.GenerateCode(decrypted, time.Now())
		require.NoError(t, err)

		// Test validation
		err = otp.ValidateTOTP(context.Background(), user, code)
		require.NoError(t, err)

		// Test invalid code
		err = otp.ValidateTOTP(context.Background(), user, "invalid")
		assert.Error(t, err)
	})

	// Test GenerateRecoveryCodes
	t.Run("GenerateRecoveryCodes", func(t *testing.T) {
		codes := otp.GenerateRecoveryCodes()
		assert.Len(t, codes, 16)

		for _, code := range codes {
			assert.Len(t, code, 8)
		}
	})

	// Test TOTPQRString
	t.Run("TOTPQRString", func(t *testing.T) {
		user := &User{
			ID:    "test-user",
			Email: sql.NullString{String: "test@example.com", Valid: true},
		}

		// Generate QR string
		qrString, err := otp.TOTPQRString(user)
		require.NoError(t, err)
		assert.Contains(t, qrString, "otpauth://totp/")
		assert.Contains(t, qrString, "test@example.com")
		assert.Contains(t, qrString, "issuer=test")
	})

	// Test DefaultOTPDelivery and DefaultName
	t.Run("UserDefaults", func(t *testing.T) {
		// Test with email
		userWithEmail := &User{
			ID:    "test-user",
			Email: sql.NullString{String: "test@example.com", Valid: true},
		}
		assert.Equal(t, Email, userWithEmail.DefaultOTPDelivery())
		assert.Equal(t, "test@example.com", userWithEmail.DefaultName())

		// Test with phone
		userWithPhone := &User{
			ID:    "test-user",
			Phone: sql.NullString{String: "+49123456789", Valid: true},
		}
		assert.Equal(t, Phone, userWithPhone.DefaultOTPDelivery())
		assert.Equal(t, "+49123456789", userWithPhone.DefaultName())

		// Test with both (email should be preferred)
		userWithBoth := &User{
			ID:    "test-user",
			Email: sql.NullString{String: "test@example.com", Valid: true},
			Phone: sql.NullString{String: "+49123456789", Valid: true},
		}
		assert.Equal(t, Email, userWithBoth.DefaultOTPDelivery())
		assert.Equal(t, "test@example.com", userWithBoth.DefaultName())
	})

	// Test WithNATS
	t.Run("WithNATS", func(t *testing.T) {
		store := newMockStore()
		otp := New(store)
		assert.NotNil(t, otp)
	})

	// Test TOTPQRString error cases
	t.Run("TOTPQRStringErrors", func(t *testing.T) {
		// Test with no secrets
		otpNoSecrets := New(store)
		user := &User{
			ID:    "test-user",
			Email: sql.NullString{String: "test@example.com", Valid: true},
		}
		_, err := otpNoSecrets.TOTPQRString(user)
		assert.Error(t, err)

		// Test with invalid secret
		otpInvalidSecret := New(store,
			WithSecret(Secret{
				Version: 1,
				Key:     []byte("invalid-key"),
			}),
		)
		secret, _ := otpInvalidSecret.TOTPSecret(user)
		user.TOTPSecret = secret
		_, err = otpInvalidSecret.TOTPQRString(user)
		assert.Error(t, err)
	})

	// Test TOTPDecryptedSecret error cases
	t.Run("TOTPDecryptedSecretErrors", func(t *testing.T) {
		// Test with no secrets
		otpNoSecrets := New(store)
		_, err := otpNoSecrets.TOTPDecryptedSecret("any-secret")
		assert.Error(t, err)

		// Test with invalid base64
		_, err = otp.TOTPDecryptedSecret("invalid-base64")
		assert.Error(t, err)

		// Test with too short ciphertext
		shortSecret := base64.StdEncoding.EncodeToString([]byte("short"))
		_, err = otp.TOTPDecryptedSecret(shortSecret)
		assert.Error(t, err)
	})
}
