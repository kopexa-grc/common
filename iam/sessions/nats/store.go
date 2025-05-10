// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package nats

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/kopexa-grc/common/iam/sessions"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

// Store implements the sessions.Store interface using NATS JetStream
type Store[T any] struct {
	js     jetstream.JetStream
	kv     jetstream.KeyValue
	config Config
}

// Config contains the configuration for the NATS store
type Config struct {
	// BucketName is the name of the JetStream KV bucket
	BucketName string

	// MaxAge is the maximum age of the session in seconds
	MaxAge int

	// ServerURL is the NATS server URL
	ServerURL string
}

// SessionData represents the stored session data with metadata
type SessionData[T any] struct {
	Session   *sessions.Session[T] `json:"session"`
	IP        string               `json:"ip"`
	UserAgent string               `json:"user_agent"`
	LastSeen  time.Time            `json:"last_seen"`
}

// Option is a function that configures a Store
type Option func(*Config)

// WithBucketName sets the bucket name
func WithBucketName(name string) Option {
	return func(c *Config) {
		c.BucketName = name
	}
}

// WithMaxAge sets the session max age
func WithMaxAge(maxAge int) Option {
	return func(c *Config) {
		c.MaxAge = maxAge
	}
}

// WithServerURL sets the NATS server URL
func WithServerURL(url string) Option {
	return func(c *Config) {
		c.ServerURL = url
	}
}

// Validate checks the configuration
func (c *Config) Validate() error {
	if c.BucketName == "" {
		return errors.New("bucket name is required")
	}
	if c.MaxAge <= 0 {
		return errors.New("max age must be positive")
	}
	if c.ServerURL == "" {
		return errors.New("server URL is required")
	}
	return nil
}

// NewStore creates a new NATS store with the given options
func NewStore[T any](opts ...Option) (*Store[T], error) {
	config := Config{
		BucketName: "sessions",
		MaxAge:     sessions.DefaultMaxAge,
		ServerURL:  nats.DefaultURL,
	}

	for _, opt := range opts {
		opt(&config)
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	// Connect to NATS
	nc, err := nats.Connect(config.ServerURL)
	if err != nil {
		return nil, err
	}

	// Create JetStream context
	js, err := jetstream.New(nc)
	if err != nil {
		return nil, err
	}

	// Create or get KV bucket
	kv, err := js.CreateKeyValue(context.Background(), jetstream.KeyValueConfig{
		Bucket:      config.BucketName,
		Description: "Session store",
		TTL:         time.Duration(config.MaxAge) * time.Second,
	})
	if err != nil {
		if !errors.Is(err, jetstream.ErrBucketExists) {
			return nil, err
		}
		kv, err = js.KeyValue(context.Background(), config.BucketName)
		if err != nil {
			return nil, err
		}
	}

	return &Store[T]{
		js:     js,
		kv:     kv,
		config: config,
	}, nil
}

// Save persists the session data in NATS KV store
func (s *Store[T]) Save(w http.ResponseWriter, session *sessions.Session[T]) error {
	// Get client IP and User-Agent
	ip := w.Header().Get("X-Real-IP")
	if ip == "" {
		ip = w.Header().Get("X-Forwarded-For")
	}
	userAgent := w.Header().Get("User-Agent")

	// Create session data with metadata
	data := SessionData[T]{
		Session:   session,
		IP:        ip,
		UserAgent: userAgent,
		LastSeen:  time.Now(),
	}

	// Marshal session data
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Store in NATS KV
	_, err = s.kv.Put(context.Background(), session.ID, bytes)
	if err != nil {
		return err
	}

	// Set cookie
	http.SetCookie(w, &http.Cookie{
		Name:     session.Name,
		Value:    session.ID,
		Path:     sessions.CookiePath,
		MaxAge:   s.config.MaxAge,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	return nil
}

// Load retrieves the session data from NATS KV store
func (s *Store[T]) Load(r *http.Request, name string) (*sessions.Session[T], error) {
	// Get session ID from cookie
	cookie, err := r.Cookie(name)
	if err != nil {
		if err == http.ErrNoCookie {
			return nil, sessions.ErrInvalidSession
		}
		return nil, err
	}

	// Get session data from NATS KV
	entry, err := s.kv.Get(context.Background(), cookie.Value)
	if err != nil {
		if errors.Is(err, jetstream.ErrKeyNotFound) {
			return nil, sessions.ErrInvalidSession
		}
		return nil, err
	}

	// Unmarshal session data
	var data SessionData[T]
	if err := json.Unmarshal(entry.Value(), &data); err != nil {
		return nil, err
	}

	// Check if session is expired
	if data.Session.IsExpired() {
		return nil, sessions.ErrSessionExpired
	}

	// Update last seen timestamp
	data.LastSeen = time.Now()
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	// Update session data
	_, err = s.kv.Put(context.Background(), cookie.Value, bytes)
	if err != nil {
		return nil, err
	}

	return data.Session, nil
}

// Destroy removes the session data from NATS KV store
func (s *Store[T]) Destroy(w http.ResponseWriter, r *http.Request, name string) {
	// Clear cookie
	http.SetCookie(w, &http.Cookie{
		Name:   name,
		Value:  "",
		Path:   sessions.CookiePath,
		MaxAge: -1,
	})

	// Get session ID from cookie
	cookie, err := r.Cookie(name)
	if err != nil {
		return
	}

	// Delete from NATS KV
	_ = s.kv.Delete(context.Background(), cookie.Value)
}

// GetActiveSessions returns all active sessions
func (s *Store[T]) GetActiveSessions() ([]SessionData[T], error) {
	keys, err := s.kv.Keys(context.Background())
	if err != nil {
		return nil, err
	}

	var sessions []SessionData[T]
	for _, key := range keys {
		entry, err := s.kv.Get(context.Background(), key)
		if err != nil {
			continue
		}

		var data SessionData[T]
		if err := json.Unmarshal(entry.Value(), &data); err != nil {
			continue
		}

		if !data.Session.IsExpired() {
			sessions = append(sessions, data)
		}
	}

	return sessions, nil
}
