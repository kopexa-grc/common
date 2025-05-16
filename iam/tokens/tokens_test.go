// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package tokens_test

import (
	"testing"
	"time"

	"github.com/kopexa-grc/common/iam/tokens"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSigningInfo(t *testing.T) {
	t.Run("NewSigningInfo", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			expires := time.Hour
			info, err := tokens.NewSigningInfo(expires)
			require.NoError(t, err)
			assert.NotNil(t, info.Nonce)
			assert.Equal(t, nonceLength, len(info.Nonce))
			assert.True(t, info.ExpiresAt.After(time.Now()))
			assert.True(t, info.ExpiresAt.Before(time.Now().Add(expires+time.Minute)))
		})

		t.Run("zero expiration", func(t *testing.T) {
			_, err := tokens.NewSigningInfo(0)
			assert.ErrorIs(t, err, tokens.ErrExpirationIsRequired)
		})

		t.Run("negative expiration", func(t *testing.T) {
			info, err := tokens.NewSigningInfo(-time.Hour)
			require.NoError(t, err)
			assert.True(t, info.IsExpired())
		})
	})

	t.Run("IsExpired", func(t *testing.T) {
		t.Run("not expired", func(t *testing.T) {
			info, err := tokens.NewSigningInfo(time.Hour)
			require.NoError(t, err)
			assert.False(t, info.IsExpired())
		})

		t.Run("expired", func(t *testing.T) {
			info := tokens.SigningInfo{
				ExpiresAt: time.Now().Add(-time.Hour),
			}
			assert.True(t, info.IsExpired())
		})

		t.Run("expired at now", func(t *testing.T) {
			info := tokens.SigningInfo{
				ExpiresAt: time.Now(),
			}
			assert.True(t, info.IsExpired())
		})
	})
}

func TestOrganizationInviteToken(t *testing.T) {
	t.Run("NewOrganizationInviteToken", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			email := "test@example.com"
			orgID := "org123"
			token, err := tokens.NewOrganizationInviteToken(email, orgID)
			require.NoError(t, err)
			assert.Equal(t, email, token.Email)
			assert.Equal(t, orgID, token.OrganizationID)
			assert.NotNil(t, token.Nonce)
			assert.Equal(t, nonceLength, len(token.Nonce))
			assert.True(t, token.ExpiresAt.After(time.Now()))
		})

		t.Run("empty email", func(t *testing.T) {
			_, err := tokens.NewOrganizationInviteToken("", "org123")
			assert.ErrorIs(t, err, tokens.ErrInviteTokenMissingEmail)
		})

		t.Run("empty organization ID", func(t *testing.T) {
			token, err := tokens.NewOrganizationInviteToken("test@example.com", "")
			require.NoError(t, err)
			assert.Empty(t, token.OrganizationID)
		})
	})

	t.Run("Sign and Verify", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			token, err := tokens.NewOrganizationInviteToken("test@example.com", "org123")
			require.NoError(t, err)

			signature, secret, err := token.Sign()
			require.NoError(t, err)
			assert.NotEmpty(t, signature)
			assert.Equal(t, nonceLength+keyLength, len(secret))

			err = token.Verify(signature, secret)
			assert.NoError(t, err)
		})

		t.Run("expired token", func(t *testing.T) {
			token, err := tokens.NewOrganizationInviteToken("test@example.com", "org123")
			require.NoError(t, err)

			signature, secret, err := token.Sign()
			require.NoError(t, err)

			expiredToken := *token
			expiredToken.ExpiresAt = time.Now().Add(-time.Hour)
			err = expiredToken.Verify(signature, secret)
			assert.ErrorIs(t, err, tokens.ErrTokenExpired)
		})

		t.Run("invalid secret length", func(t *testing.T) {
			token, err := tokens.NewOrganizationInviteToken("test@example.com", "org123")
			require.NoError(t, err)

			signature, _, err := token.Sign()
			require.NoError(t, err)

			invalidSecret := make([]byte, nonceLength+keyLength-1)
			err = token.Verify(signature, invalidSecret)
			assert.ErrorIs(t, err, tokens.ErrInvalidSecret)
		})

		t.Run("empty email", func(t *testing.T) {
			token, err := tokens.NewOrganizationInviteToken("test@example.com", "org123")
			require.NoError(t, err)

			signature, secret, err := token.Sign()
			require.NoError(t, err)

			emptyEmailToken := *token
			emptyEmailToken.Email = ""
			err = emptyEmailToken.Verify(signature, secret)
			assert.ErrorIs(t, err, tokens.ErrInviteTokenMissingEmail)
		})

		t.Run("modified token data", func(t *testing.T) {
			token, err := tokens.NewOrganizationInviteToken("test@example.com", "org123")
			require.NoError(t, err)

			signature, secret, err := token.Sign()
			require.NoError(t, err)

			modifiedToken := *token
			modifiedToken.Email = "modified@example.com"
			err = modifiedToken.Verify(signature, secret)
			assert.ErrorIs(t, err, tokens.ErrTokenInvalid)
		})

		t.Run("invalid signature", func(t *testing.T) {
			token, err := tokens.NewOrganizationInviteToken("test@example.com", "org123")
			require.NoError(t, err)

			_, secret, err := token.Sign()
			require.NoError(t, err)

			err = token.Verify("invalid-signature", secret)
			assert.Error(t, err)
		})

		t.Run("modified organization ID", func(t *testing.T) {
			token, err := tokens.NewOrganizationInviteToken("test@example.com", "org123")
			require.NoError(t, err)

			signature, secret, err := token.Sign()
			require.NoError(t, err)

			modifiedToken := *token
			modifiedToken.OrganizationID = "org456"
			err = modifiedToken.Verify(signature, secret)
			assert.ErrorIs(t, err, tokens.ErrTokenInvalid)
		})
	})
}

const (
	nonceLength = 64
	keyLength   = 64
)
