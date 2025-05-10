// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package sessions

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// mockStore is a mock implementation of the Store interface
type mockStore[T any] struct {
	sessions  map[string]*Session[T]
	saveErr   error
	loadErr   error
	destroyed bool
	session   *Session[T]
}

func newMockStore[T any]() *mockStore[T] {
	return &mockStore[T]{
		sessions: make(map[string]*Session[T]),
	}
}

func (s *mockStore[T]) Save(_ http.ResponseWriter, session *Session[T]) error {
	if s.saveErr != nil {
		return s.saveErr
	}

	s.sessions[session.Name] = session

	return nil
}

func (s *mockStore[T]) Load(_ *http.Request, _ string) (*Session[T], error) {
	if s.loadErr != nil {
		return nil, s.loadErr
	}

	return s.session, nil
}

func (s *mockStore[T]) Destroy(_ http.ResponseWriter, _ *http.Request, _ string) {
	s.destroyed = true
}

func TestNewSession(t *testing.T) {
	store := newMockStore[string]()
	session := NewSession(store, "test")

	assert.NotEmpty(t, session.ID)
	assert.Equal(t, "test", session.Name)
	assert.NotNil(t, session.Values)
	assert.False(t, session.CreatedAt.IsZero())
	assert.False(t, session.ExpiresAt.IsZero())
}

func TestSession_SetGet(t *testing.T) {
	store := newMockStore[string]()
	session := NewSession(store, "test")

	// Test Set and Get
	session.Set("key", "value")
	assert.Equal(t, "value", session.Get("key"))

	// Test GetOk
	value, ok := session.GetOk("key")
	assert.True(t, ok)
	assert.Equal(t, "value", value)

	// Test non-existent key
	value, ok = session.GetOk("nonexistent")
	assert.False(t, ok)
	assert.Empty(t, value)
}

func TestSession_DeleteClear(t *testing.T) {
	store := newMockStore[string]()
	session := NewSession(store, "test")

	// Set some values
	session.Set("key1", "value1")
	session.Set("key2", "value2")

	// Test Delete
	session.Delete("key1")
	_, ok := session.GetOk("key1")
	assert.False(t, ok)
	value, ok := session.GetOk("key2")
	assert.True(t, ok)
	assert.Equal(t, "value2", value)

	// Test Clear
	session.Clear()
	_, ok = session.GetOk("key2")
	assert.False(t, ok)
}

func TestSession_IsExpired(t *testing.T) {
	store := newMockStore[string]()
	session := NewSession(store, "test")

	// Session should not be expired
	assert.False(t, session.IsExpired())

	// Set expiration to past
	session.ExpiresAt = time.Now().Add(-time.Hour)
	assert.True(t, session.IsExpired())
}

func TestGenerateSessionID(t *testing.T) {
	id1 := GenerateSessionID()
	id2 := GenerateSessionID()

	assert.NotEmpty(t, id1)
	assert.NotEmpty(t, id2)
	assert.NotEqual(t, id1, id2)
}

func TestValidateKeyLength(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		wantErr bool
	}{
		{
			name:    "valid 16 bytes",
			key:     "1234567890123456",
			wantErr: false,
		},
		{
			name:    "valid 32 bytes",
			key:     "12345678901234567890123456789012",
			wantErr: false,
		},
		{
			name:    "valid 64 bytes",
			key:     "1234567890123456789012345678901234567890123456789012345678901234",
			wantErr: false,
		},
		{
			name:    "invalid length",
			key:     "12345",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateKeyLength(tt.key)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSession_Rotate(t *testing.T) {
	store := newMockStore[string]()
	session := NewSession(store, "test")
	oldID := session.ID
	oldCreatedAt := session.CreatedAt

	// Wait a bit to ensure timestamps are different
	time.Sleep(time.Millisecond)

	// Rotate session
	session.Rotate()

	// Verify new ID and timestamp
	assert.NotEqual(t, oldID, session.ID, "session ID should change")
	assert.NotEqual(t, oldCreatedAt, session.CreatedAt, "created timestamp should update")
	assert.True(t, session.CreatedAt.After(oldCreatedAt), "new timestamp should be after old timestamp")
}

func TestSession_Expiration(t *testing.T) {
	store := newMockStore[string]()
	session := NewSession(store, "test")

	// Test standard expiration
	assert.False(t, session.IsExpired(), "New session should not be expired")

	// Test custom expiration
	session.ExpiresAt = time.Now().Add(time.Hour)
	assert.False(t, session.IsExpired(), "Future expiration should not be expired")

	session.ExpiresAt = time.Now().Add(-time.Hour)
	assert.True(t, session.IsExpired(), "Past expiration should be expired")
}

func TestSession_ConcurrentAccess(t *testing.T) {
	store := newMockStore[string]()
	session := NewSession(store, "test")

	// Test concurrent writes
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()
			session.Set(fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i))
		}(i)
	}

	wg.Wait()

	// Verify all values were set
	for i := 0; i < 10; i++ {
		value, ok := session.GetOk(fmt.Sprintf("key%d", i))
		assert.True(t, ok, "Value should exist")
		assert.Equal(t, fmt.Sprintf("value%d", i), value)
	}
}

func TestSession_ErrorHandling(t *testing.T) {
	store := newMockStore[string]()
	session := NewSession(store, "test")

	// Test Save error
	store.saveErr = ErrSaveFailed
	err := session.Save(nil)
	assert.Error(t, err)
	assert.Equal(t, "save error", err.Error())

	// Test Load error
	store.loadErr = ErrLoadFailed
	_, err = store.Load(nil, "test")
	assert.Error(t, err)
	assert.Equal(t, "load error", err.Error())
}

func TestConfigOptions(t *testing.T) {
	store := newMockStore[string]()
	cookieCfg := &CookieConfig{
		Name:     "test",
		Domain:   "example.com",
		MaxAge:   123,
		Secure:   false,
		HTTPOnly: false,
		SameSite: 2, // Strict
	}
	cfg := NewConfig[string](store,
		WithCookieConfig[string](cookieCfg),
		WithMaxAge[string](456),
		WithSecure[string](true),
		WithHTTPOnly[string](true),
		WithSameSite[string](1), // Lax
		WithDomain[string]("test.com"),
	)
	assert.Equal(t, store, cfg.Store)
	assert.Equal(t, "test", cfg.CookieConfig.Name)
	assert.Equal(t, "test.com", cfg.CookieConfig.Domain)
	assert.Equal(t, 456, cfg.CookieConfig.MaxAge)
	assert.True(t, cfg.CookieConfig.Secure)
	assert.True(t, cfg.CookieConfig.HTTPOnly)
	assert.Equal(t, 1, int(cfg.CookieConfig.SameSite))
}

func TestContextHelpers(t *testing.T) {
	sess := &Session[string]{ID: "abc"}
	ctx := context.Background()
	ctx = WithSession[string](ctx, sess)
	// FromSession
	got, ok := FromSession[string](ctx)
	assert.True(t, ok)
	assert.Equal(t, sess, got)
	// MustFromSession
	assert.Panics(t, func() { MustFromSession[string](context.Background()) })
	assert.Equal(t, sess, MustFromSession[string](ctx))
	// FromSessionOr
	alt := &Session[string]{ID: "alt"}
	assert.Equal(t, alt, FromSessionOr[string](context.Background(), alt))
	assert.Equal(t, sess, FromSessionOr[string](ctx, alt))
	// FromSessionOrFunc
	assert.Equal(t, alt, FromSessionOrFunc[string](context.Background(), func() *Session[string] { return alt }))
	assert.Equal(t, sess, FromSessionOrFunc[string](ctx, func() *Session[string] { return alt }))
}

func TestSession_SetGetName(t *testing.T) {
	s := NewSession(newMockStore[string](), "foo")
	assert.Equal(t, "foo", s.GetName())
	s.SetName("bar")
	assert.Equal(t, "bar", s.GetName())
}

func TestSession_Destroy(t *testing.T) {
	store := &mockStore[string]{
		session: NewSession[string](nil, "test"),
	}
	session := NewSession(store, "test")
	session.Destroy(nil, nil)
	assert.True(t, store.destroyed, "session should be destroyed")
}

func TestEncryptDecrypt(t *testing.T) {
	key := "12345678901234567890123456789012"
	plain := []byte("hello world")
	enc, err := encrypt(plain, key)
	assert.NoError(t, err)
	dec, err := decrypt(enc, key)
	assert.NoError(t, err)
	assert.Equal(t, plain, dec)
}

func TestEncryptDecrypt_InvalidKey(t *testing.T) {
	key := "short"
	_, err := encrypt([]byte("data"), key)
	assert.Error(t, err)
	_, err = decrypt([]byte("data"), key)
	assert.Error(t, err)
}

func TestEncryptDecrypt_CiphertextTooShort(t *testing.T) {
	key := "12345678901234567890123456789012"
	_, err := decrypt([]byte("short"), key)
	assert.ErrorIs(t, err, ErrCiphertextTooShort)
}

func TestEncodeDecodeSession(t *testing.T) {
	key := "12345678901234567890123456789012"
	s := NewSession(newMockStore[string](), "foo")
	s.Set("bar", "baz")
	enc, err := EncodeSession(s, key)
	assert.NoError(t, err)
	dec, err := DecodeSession[string](enc, key)
	assert.NoError(t, err)
	assert.Equal(t, s.ID, dec.ID)
	assert.Equal(t, s.Name, dec.Name)
	assert.Equal(t, "baz", dec.Get("bar"))
}

func TestEncodeSession_MarshalError(t *testing.T) {
	key := "12345678901234567890123456789012"
	s := NewSession(newMockStore[any](), "foo")
	s.Set("bad", make(chan int)) // Channels sind nicht serialisierbar
	_, err := EncodeSession(s, key)
	assert.Error(t, err)
}

func TestDecodeSession_DecodeError(t *testing.T) {
	key := "12345678901234567890123456789012"
	_, err := DecodeSession[string]("!!!", key)
	assert.Error(t, err)
}

func TestDecodeSession_DecryptError(t *testing.T) {
	key := "12345678901234567890123456789012"
	// Manipuliertes base64, das zu ungültigem Ciphertext führt
	enc := "aGVsbG8="
	_, err := DecodeSession[string](enc, key)
	assert.Error(t, err)
}

func TestDecodeSession_UnmarshalError(t *testing.T) {
	key := "12345678901234567890123456789012"
	// Verschlüsselte, aber keine gültige Session-Struktur
	plain := []byte("not a session struct")
	enc, err := encrypt(plain, key)
	assert.NoError(t, err)

	b64 := base64.URLEncoding.EncodeToString(enc)
	_, err = DecodeSession[string](b64, key)
	assert.Error(t, err)
}
