// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

// Package sessions provides a type-safe session management system for web applications.
// It supports generic types for session values and includes features like:
// - Secure session ID generation using UUID v4
// - Configurable session storage
// - Cookie-based session management
// - Thread-safe operations
package sessions

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// Common errors that can occur during session operations
var (
	ErrInvalidKeyLength   = errors.New("key must be 16, 32, or 64 bytes")
	ErrInvalidSession     = errors.New("invalid session")
	ErrSessionExpired     = errors.New("session has expired")
	ErrCiphertextTooShort = errors.New("ciphertext too short")
)

// Store defines the interface for session storage implementations
type Store[T any] interface {
	// Save persists the session data
	Save(w http.ResponseWriter, session *Session[T]) error

	// Load retrieves the session data
	Load(r *http.Request, name string) (*Session[T], error)

	// Destroy removes the session
	Destroy(w http.ResponseWriter, r *http.Request, name string)
}

// Session represents a user session with type-safe values
type Session[T any] struct {
	// ID is the unique identifier for this session
	ID string `json:"id"`

	// Name is the name of the session
	Name string `json:"name"`

	// Values contains the session data
	Values map[string]T `json:"values"`

	// CreatedAt is the timestamp when the session was created
	CreatedAt time.Time `json:"createdAt"`

	// ExpiresAt is the timestamp when the session will expire
	ExpiresAt time.Time `json:"expiresAt"`

	mu    sync.RWMutex
	store Store[T]
}

// NewSession creates a new session with the given store and name
func NewSession[T any](store Store[T], name string) *Session[T] {
	now := time.Now()
	return &Session[T]{
		ID:        GenerateSessionID(),
		Name:      name,
		Values:    make(map[string]T),
		CreatedAt: now,
		ExpiresAt: now.Add(time.Hour), // Default 1 hour expiration
		store:     store,
	}
}

// SetName sets the name of the session
func (s *Session[T]) SetName(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Name = name
}

// GetName returns the name of the session
func (s *Session[T]) GetName() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Name
}

// Set stores a value in the session
func (s *Session[T]) Set(key string, value T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Values[key] = value
}

// Get retrieves a value from the session
func (s *Session[T]) Get(key string) T {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Values[key]
}

// GetOk retrieves a value from the session and indicates if it exists
func (s *Session[T]) GetOk(key string) (T, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, ok := s.Values[key]
	return value, ok
}

// Delete removes a value from the session
func (s *Session[T]) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.Values, key)
}

// Clear removes all values from the session
func (s *Session[T]) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Values = make(map[string]T)
}

// Save persists the session to the store
func (s *Session[T]) Save(w http.ResponseWriter) error {
	return s.store.Save(w, s)
}

// Destroy removes the session from the store
func (s *Session[T]) Destroy(w http.ResponseWriter, r *http.Request) {
	s.store.Destroy(w, r, s.Name)
}

// IsExpired checks if the session has expired
func (s *Session[T]) IsExpired() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return time.Now().After(s.ExpiresAt)
}

// Rotate generates a new session ID while preserving the session data
// This helps prevent session fixation attacks
func (s *Session[T]) Rotate() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ID = GenerateSessionID()
	s.CreatedAt = time.Now()
}

// GenerateSessionID generates a new cryptographically secure random session ID (256 Bit)
func GenerateSessionID() string {
	b := make([]byte, 32) // 256 Bit
	_, err := rand.Read(b)
	if err != nil {
		panic("failed to generate secure session ID: " + err.Error())
	}
	return fmt.Sprintf("%x", b)
}

// validateKeyLength checks if the key length is valid for AES encryption
func validateKeyLength(key string) error {
	length := len(key)
	if length != MinKeyLength && length != DefaultKeyLength && length != MaxKeyLength {
		return fmt.Errorf("%w: got %d bytes, want %d, %d, or %d bytes", ErrInvalidKeyLength, length, MinKeyLength, DefaultKeyLength, MaxKeyLength)
	}
	return nil
}

// encrypt encrypts the session data using AES-GCM
func encrypt(data []byte, key string) ([]byte, error) {
	if err := validateKeyLength(key); err != nil {
		return nil, err
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	return gcm.Seal(nonce, nonce, data, nil), nil
}

// decrypt decrypts the session data using AES-GCM
func decrypt(data []byte, key string) ([]byte, error) {
	if err := validateKeyLength(key); err != nil {
		return nil, err
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	if len(data) < gcm.NonceSize() {
		return nil, ErrCiphertextTooShort
	}

	nonce := data[:gcm.NonceSize()]
	ciphertext := data[gcm.NonceSize():]

	return gcm.Open(nil, nonce, ciphertext, nil)
}

// EncodeSession encodes the session data to a base64 string
func EncodeSession[T any](session *Session[T], key string) (string, error) {
	data, err := json.Marshal(session)
	if err != nil {
		return "", fmt.Errorf("failed to marshal session: %w", err)
	}

	encrypted, err := encrypt(data, key)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt session: %w", err)
	}

	return base64.URLEncoding.EncodeToString(encrypted), nil
}

// DecodeSession decodes the session data from a base64 string
func DecodeSession[T any](data string, key string) (*Session[T], error) {
	decoded, err := base64.URLEncoding.DecodeString(data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode session: %w", err)
	}

	decrypted, err := decrypt(decoded, key)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt session: %w", err)
	}

	var session Session[T]
	if err := json.Unmarshal(decrypted, &session); err != nil {
		return nil, fmt.Errorf("failed to unmarshal session: %w", err)
	}

	return &session, nil
}
