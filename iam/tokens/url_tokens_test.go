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

func TestVerificationToken(t *testing.T) {
	t.Run("success sign/verify", func(t *testing.T) {
		vt, err := tokens.NewVerificationToken("user@example.com")
		require.NoError(t, err)
		sig, secret, err := vt.Sign()
		require.NoError(t, err)
		assert.NotEmpty(t, sig)
		assert.Len(t, secret, 128) // nonceLength + keyLength (both 64)
		err = vt.Verify(sig, secret)
		assert.NoError(t, err)
	})

	t.Run("missing email at construction", func(t *testing.T) {
		vt, err := tokens.NewVerificationToken("")
		assert.Nil(t, vt)
		assert.ErrorIs(t, err, tokens.ErrMissingEmail)
	})

	t.Run("modified email causes invalid", func(t *testing.T) {
		vt, err := tokens.NewVerificationToken("user@example.com")
		require.NoError(t, err)
		sig, secret, err := vt.Sign()
		require.NoError(t, err)

		clone := *vt
		clone.Email = "other@example.com"
		err = clone.Verify(sig, secret)
		assert.ErrorIs(t, err, tokens.ErrTokenInvalid)
	})

	t.Run("expired token", func(t *testing.T) {
		vt, err := tokens.NewVerificationToken("user@example.com")
		require.NoError(t, err)
		sig, secret, err := vt.Sign()
		require.NoError(t, err)

		expired := *vt
		expired.ExpiresAt = time.Now().Add(-time.Minute)
		err = expired.Verify(sig, secret)
		assert.ErrorIs(t, err, tokens.ErrTokenExpired)
	})
}

func TestResetToken(t *testing.T) {
	t.Run("construction requires user id", func(t *testing.T) {
		rt, err := tokens.NewResetToken("")
		assert.Nil(t, rt)
		assert.ErrorIs(t, err, tokens.ErrMissingUserID)
	})

	t.Run("sign/verify success", func(t *testing.T) {
		rt, err := tokens.NewResetToken("user-1")
		require.NoError(t, err)
		sig, secret, err := rt.Sign()
		require.NoError(t, err)
		err = rt.Verify(sig, secret)
		assert.NoError(t, err)
	})

	t.Run("tampered user id", func(t *testing.T) {
		rt, err := tokens.NewResetToken("user-1")
		require.NoError(t, err)
		sig, secret, err := rt.Sign()
		require.NoError(t, err)

		clone := *rt
		clone.UserID = "user-2"
		err = clone.Verify(sig, secret)
		assert.ErrorIs(t, err, tokens.ErrTokenInvalid)
	})

	t.Run("expired reset token", func(t *testing.T) {
		rt, err := tokens.NewResetToken("user-1")
		require.NoError(t, err)
		sig, secret, err := rt.Sign()
		require.NoError(t, err)

		expired := *rt
		expired.ExpiresAt = time.Now().Add(-time.Minute)
		err = expired.Verify(sig, secret)
		assert.ErrorIs(t, err, tokens.ErrTokenExpired)
	})

	t.Run("invalid secret length", func(t *testing.T) {
		rt, err := tokens.NewResetToken("user-1")
		require.NoError(t, err)
		sig, secret, err := rt.Sign()
		require.NoError(t, err)

		shortSecret := secret[:len(secret)-1]
		err = rt.Verify(sig, shortSecret)
		assert.ErrorIs(t, err, tokens.ErrInvalidSecret)
	})
}
