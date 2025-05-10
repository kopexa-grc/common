// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package nats

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kopexa-grc/common/iam/sessions"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStore_SaveLoad(t *testing.T) {
	// Start embedded NATS server
	s := startTestServer(t)

	bucket := "test_sessions_saveload_" + time.Now().Format("150405_000000")
	store, err := NewStore[string](
		WithServerURL(getTestServerURL(s)),
		WithBucketName(bucket),
		WithMaxAge(3600),
	)
	require.NoError(t, err)

	session := sessions.NewSession(store, "test")
	session.Set("key", "value")

	// Create response recorder
	w := httptest.NewRecorder()
	w.Header().Set("X-Real-IP", "127.0.0.1")
	w.Header().Set("User-Agent", "test-agent")

	// Save session
	err = store.Save(w, session)
	require.NoError(t, err)

	// Get cookie from response
	cookies := w.Result().Cookies()
	require.Len(t, cookies, 1)
	cookie := cookies[0]

	// Create new request with cookie
	r := httptest.NewRequest("GET", "/", nil)
	r.AddCookie(cookie)

	// Load session
	loaded, err := store.Load(r, "test")
	require.NoError(t, err)

	// Verify session data
	assert.Equal(t, session.ID, loaded.ID)
	assert.Equal(t, session.Name, loaded.Name)
	assert.Equal(t, "value", loaded.Get("key"))

	// Get active sessions
	activeSessions, err := store.GetActiveSessions()
	require.NoError(t, err)
	require.Len(t, activeSessions, 1)

	// Verify session metadata
	assert.Equal(t, "127.0.0.1", activeSessions[0].IP)
	assert.Equal(t, "test-agent", activeSessions[0].UserAgent)
	assert.WithinDuration(t, time.Now(), activeSessions[0].LastSeen, time.Second)
}

func TestStore_Destroy(t *testing.T) {
	// Start embedded NATS server
	s := startTestServer(t)

	bucket := "test_sessions_destroy_" + time.Now().Format("150405_000000")
	store, err := NewStore[string](
		WithServerURL(getTestServerURL(s)),
		WithBucketName(bucket),
	)
	require.NoError(t, err)

	session := sessions.NewSession(store, "test")

	// Save session
	w := httptest.NewRecorder()
	err = store.Save(w, session)
	require.NoError(t, err)

	// Get cookie
	cookies := w.Result().Cookies()
	require.Len(t, cookies, 1)
	cookie := cookies[0]

	// Create request with cookie
	r := httptest.NewRequest("GET", "/", nil)
	r.AddCookie(cookie)

	// Verify session exists
	_, err = store.Load(r, "test")
	require.NoError(t, err)

	// Destroy session
	store.Destroy(w, r, "test")

	// Verify session is gone
	_, err = store.Load(r, "test")
	assert.ErrorIs(t, err, sessions.ErrInvalidSession)
}

func TestStore_ExpiredSession(t *testing.T) {
	// Start embedded NATS server
	s := startTestServer(t)

	bucket := "test_sessions_expired_" + time.Now().Format("150405_000000")
	store, err := NewStore[string](
		WithServerURL(getTestServerURL(s)),
		WithBucketName(bucket),
	)
	require.NoError(t, err)

	session := sessions.NewSession(store, "test")
	session.ExpiresAt = time.Now().Add(-time.Hour)

	// Save expired session
	w := httptest.NewRecorder()
	err = store.Save(w, session)
	require.NoError(t, err)

	// Get cookie
	cookies := w.Result().Cookies()
	require.Len(t, cookies, 1)
	cookie := cookies[0]

	// Create request with cookie
	r := httptest.NewRequest("GET", "/", nil)
	r.AddCookie(cookie)

	// Try to load expired session
	_, err = store.Load(r, "test")
	assert.ErrorIs(t, err, sessions.ErrSessionExpired)
}

func TestStore_InvalidSession(t *testing.T) {
	// Start embedded NATS server
	s := startTestServer(t)

	bucket := "test_sessions_invalid_" + time.Now().Format("150405_000000")
	store, err := NewStore[string](
		WithServerURL(getTestServerURL(s)),
		WithBucketName(bucket),
	)
	require.NoError(t, err)

	// Test with no cookie
	r := httptest.NewRequest("GET", "/", nil)
	_, err = store.Load(r, "test")
	assert.ErrorIs(t, err, sessions.ErrInvalidSession)

	// Test with invalid cookie
	r.AddCookie(&http.Cookie{
		Name:  "test",
		Value: "invalid",
	})

	_, err = store.Load(r, "test")
	assert.ErrorIs(t, err, sessions.ErrInvalidSession)
}
