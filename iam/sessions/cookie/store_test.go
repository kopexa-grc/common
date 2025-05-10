// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package cookie

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
	store, err := NewStore[string](
		WithSigningKey("12345678901234567890123456789012"),
		WithEncryptionKey("12345678901234567890123456789012"),
		WithMaxAge(3600),
		WithSecure(true),
		WithHTTPOnly(true),
		WithSameSite(sessions.CookieSameSiteLax),
	)
	require.NoError(t, err)

	session := sessions.NewSession(store, "test")
	session.Set("key", "value")

	// Create response recorder
	w := httptest.NewRecorder()

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
}

func TestStore_Destroy(t *testing.T) {
	store, err := NewStore[string](
		WithSigningKey("12345678901234567890123456789012"),
		WithEncryptionKey("12345678901234567890123456789012"),
	)
	require.NoError(t, err)

	session := sessions.NewSession(store, "test")

	// Save session
	w := httptest.NewRecorder()
	err = store.Save(w, session)
	require.NoError(t, err)

	// Destroy session
	store.Destroy(w, nil, "test")

	// Verify at least one cookie is deleted (MaxAge == -1)
	cookies := w.Result().Cookies()
	found := false
	for _, c := range cookies {
		if c.MaxAge == -1 {
			found = true
			break
		}
	}
	assert.True(t, found, "no deleted cookie (MaxAge == -1) found")
}

func TestStore_InvalidSession(t *testing.T) {
	store, err := NewStore[string](
		WithSigningKey("12345678901234567890123456789012"),
		WithEncryptionKey("12345678901234567890123456789012"),
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
	assert.Error(t, err)
}

func TestStore_ExpiredSession(t *testing.T) {
	store, err := NewStore[string](
		WithSigningKey("12345678901234567890123456789012"),
		WithEncryptionKey("12345678901234567890123456789012"),
	)
	require.NoError(t, err)

	session := sessions.NewSession(store, "test")
	session.ExpiresAt = time.Now().Add(-time.Hour)

	// Save expired session
	w := httptest.NewRecorder()
	err = store.Save(w, session)
	require.NoError(t, err)

	// Try to load expired session
	r := httptest.NewRequest("GET", "/", nil)
	r.AddCookie(w.Result().Cookies()[0])
	_, err = store.Load(r, "test")
	assert.ErrorIs(t, err, sessions.ErrSessionExpired)
}

// Security validation tests
func TestStore_InvalidConfig(t *testing.T) {
	_, err := NewStore[string](
		WithSigningKey("short"),
		WithEncryptionKey("short"),
	)
	assert.Error(t, err)

	_, err = NewStore[string](
		WithSigningKey("12345678901234567890123456789012"),
		WithEncryptionKey("12345678901234567890123456789012"),
		WithMaxAge(0),
	)
	assert.Error(t, err)

	_, err = NewStore[string](
		WithSigningKey("12345678901234567890123456789012"),
		WithEncryptionKey("12345678901234567890123456789012"),
		WithSameSite("invalid"),
	)
	assert.Error(t, err)

	_, err = NewStore[string](
		WithSigningKey("12345678901234567890123456789012"),
		WithEncryptionKey("12345678901234567890123456789012"),
		WithSecure(false),
	)
	assert.Error(t, err)

	_, err = NewStore[string](
		WithSigningKey("12345678901234567890123456789012"),
		WithEncryptionKey("12345678901234567890123456789012"),
		WithHTTPOnly(false),
	)
	assert.Error(t, err)
}

func TestStore_SubdomainSupport(t *testing.T) {
	// Test mit Domain ohne Punkt
	store, err := NewStore[string](
		WithSigningKey("12345678901234567890123456789012"),
		WithEncryptionKey("12345678901234567890123456789012"),
		WithDomain("kopexa.com"),
		WithMaxAge(3600),
		WithSecure(true),
		WithHTTPOnly(true),
		WithSameSite(sessions.CookieSameSiteLax),
	)
	require.NoError(t, err)

	session := sessions.NewSession(store, "test")
	session.Set("key", "value")

	// Test Save
	w := httptest.NewRecorder()
	err = store.Save(w, session)
	require.NoError(t, err)

	// Überprüfe Cookie-Einstellungen
	cookies := w.Result().Cookies()
	require.Len(t, cookies, 1)
	cookie := cookies[0]
	assert.Equal(t, "kopexa.com", cookie.Domain, "Cookie domain should be normalized (dot is not present in Set-Cookie)")
	assert.True(t, cookie.Secure, "Cookie should be secure")
	assert.True(t, cookie.HttpOnly, "Cookie should be httpOnly")

	// Test Load von verschiedenen Subdomains
	domains := []string{
		"console.kopexa.com",
		"auth.kopexa.com",
		"api.kopexa.com",
	}

	for _, domain := range domains {
		r := httptest.NewRequest("GET", "https://"+domain+"/", nil)
		r.AddCookie(cookie)

		loaded, err := store.Load(r, "test")
		require.NoError(t, err, "Should load session from %s", domain)
		assert.Equal(t, "value", loaded.Get("key"), "Session value should be preserved")
	}

	// Test mit Domain mit Punkt
	store2, err := NewStore[string](
		WithSigningKey("12345678901234567890123456789012"),
		WithEncryptionKey("12345678901234567890123456789012"),
		WithDomain(".kopexa.com"),
		WithMaxAge(3600),
		WithSecure(true),
		WithHTTPOnly(true),
		WithSameSite(sessions.CookieSameSiteLax),
	)
	require.NoError(t, err)

	session2 := sessions.NewSession(store2, "test2")
	session2.Set("key", "value2")

	w2 := httptest.NewRecorder()
	err = store2.Save(w2, session2)
	require.NoError(t, err)

	cookies2 := w2.Result().Cookies()
	require.Len(t, cookies2, 1)
	cookie2 := cookies2[0]
	assert.Equal(t, "kopexa.com", cookie2.Domain, "Cookie domain should be preserved (dot is not present in Set-Cookie)")
}

func TestStore_InvalidDomain(t *testing.T) {
	// Test mit ungültiger Domain (leerer String ist erlaubt)
	_, err := NewStore[string](
		WithSigningKey("12345678901234567890123456789012"),
		WithEncryptionKey("12345678901234567890123456789012"),
		WithDomain(""),
	)
	assert.NoError(t, err)
}

func TestStore_DevMode(t *testing.T) {
	// Test mit Entwicklungsmodus
	store, err := NewStore[string](
		WithSigningKey("12345678901234567890123456789012"),
		WithEncryptionKey("12345678901234567890123456789012"),
		WithDevMode(true),
		WithSecure(false),       // Im Dev-Modus erlaubt
		WithHTTPOnly(false),     // Im Dev-Modus erlaubt
		WithDomain("localhost"), // Im Dev-Modus erlaubt
	)
	require.NoError(t, err)

	session := sessions.NewSession(store, "test")
	session.Set("key", "value")

	// Test Save
	w := httptest.NewRecorder()
	err = store.Save(w, session)
	require.NoError(t, err)

	// Überprüfe Cookie-Einstellungen
	cookies := w.Result().Cookies()
	require.Len(t, cookies, 1)
	cookie := cookies[0]
	assert.Equal(t, "localhost", cookie.Domain)
	assert.False(t, cookie.Secure)
	assert.False(t, cookie.HttpOnly)
}

func TestStore_DevModeValidation(t *testing.T) {
	// Test: Dev-Modus erlaubt unsichere Konfiguration
	_, err := NewStore[string](
		WithSigningKey("12345678901234567890123456789012"),
		WithEncryptionKey("12345678901234567890123456789012"),
		WithDevMode(true),
		WithSecure(false),
		WithHTTPOnly(false),
		WithDomain("localhost"),
	)
	require.NoError(t, err, "Dev-Modus sollte unsichere Konfiguration erlauben")

	// Test: Ohne Dev-Modus wird unsichere Konfiguration abgelehnt
	_, err = NewStore[string](
		WithSigningKey("12345678901234567890123456789012"),
		WithEncryptionKey("12345678901234567890123456789012"),
		WithSecure(false),
		WithHTTPOnly(false),
		WithDomain("localhost"),
	)
	assert.Error(t, err, "Produktionsmodus sollte unsichere Konfiguration ablehnen")
}
