// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package cookie

import (
	"net/http"
	"strings"
	"time"

	"errors"

	"github.com/kopexa-grc/common/iam/sessions"
)

// Store implements the sessions.Store interface using HTTP cookies
type Store[T any] struct {
	config Config
}

// Config contains the configuration for the cookie store
type Config struct {
	// SigningKey must be a 16, 32, or 64 character string used to sign the cookie
	// This key should be kept secret and should be at least 32 characters long
	SigningKey string

	// EncryptionKey must be a 16, 32, or 64 character string used to encrypt the cookie
	// This key should be kept secret and should be at least 32 characters long
	EncryptionKey string

	// Domain is the domain for the cookie, leave empty to use the default value of the server
	Domain string

	// MaxAge is the maximum age of the session cookie in seconds.
	// After this time, the cookie will be invalidated.
	MaxAge int

	// Secure determines if the cookie should only be sent over HTTPS
	Secure bool

	// HTTPOnly determines if the cookie should be accessible via JavaScript
	HTTPOnly bool

	// SameSite determines the SameSite attribute of the cookie
	SameSite string

	// DevMode enables development mode with relaxed security settings
	// WARNING: Never use in production!
	DevMode bool
}

// Option is a function that configures a Store
type Option func(*Config)

// WithSigningKey sets the signing key
func WithSigningKey(key string) Option {
	return func(c *Config) {
		c.SigningKey = key
	}
}

// WithEncryptionKey sets the encryption key
func WithEncryptionKey(key string) Option {
	return func(c *Config) {
		c.EncryptionKey = key
	}
}

// WithDomain sets the cookie domain
func WithDomain(domain string) Option {
	return func(c *Config) {
		if domain != "" && !strings.HasPrefix(domain, ".") {
			c.Domain = "." + domain
		} else {
			c.Domain = domain
		}
	}
}

// WithMaxAge sets the cookie max age
func WithMaxAge(maxAge int) Option {
	return func(c *Config) {
		c.MaxAge = maxAge
	}
}

// WithSecure sets the secure flag
func WithSecure(secure bool) Option {
	return func(c *Config) {
		c.Secure = secure
	}
}

// WithHTTPOnly sets the HTTPOnly flag
func WithHTTPOnly(httpOnly bool) Option {
	return func(c *Config) {
		c.HTTPOnly = httpOnly
	}
}

// WithSameSite sets the SameSite attribute
func WithSameSite(sameSite string) Option {
	return func(c *Config) {
		c.SameSite = sameSite
	}
}

// WithDevMode enables development mode with relaxed security settings
// WARNING: Never use in production!
func WithDevMode(devMode bool) Option {
	return func(c *Config) {
		c.DevMode = devMode
	}
}

// Validate prüft die Sicherheit und Gültigkeit der Konfiguration
func (c *Config) Validate() error {
	if len(c.SigningKey) < 32 {
		return errors.New("signing key must be at least 32 bytes")
	}
	if len(c.EncryptionKey) < 32 {
		return errors.New("encryption key must be at least 32 bytes")
	}
	if c.MaxAge <= 0 {
		return errors.New("max age must be positive")
	}
	switch c.SameSite {
	case sessions.CookieSameSiteLax, sessions.CookieSameSiteStrict, sessions.CookieSameSiteNone:
		// ok
	default:
		return errors.New("invalid SameSite value")
	}

	// Im Entwicklungsmodus sind Secure und HTTPOnly optional
	if !c.DevMode {
		if !c.Secure {
			return errors.New("secure must be true for production")
		}
		if !c.HTTPOnly {
			return errors.New("httpOnly must be true for security")
		}
		if c.SameSite == sessions.CookieSameSiteNone && !c.Secure {
			return errors.New("SameSite=None requires Secure=true")
		}
		// Für Subdomains muss die Domain gesetzt sein, aber der Punkt wird automatisch ergänzt
	}

	return nil
}

// NewStore creates a new cookie store with the given options and validates the config
func NewStore[T any](opts ...Option) (*Store[T], error) {
	config := Config{
		SigningKey:    sessions.DefaultSigningKey,
		EncryptionKey: sessions.DefaultEncryptionKey,
		MaxAge:        sessions.DefaultMaxAge,
		Secure:        true,                       // Default: true für maximale Sicherheit
		HTTPOnly:      true,                       // Default: true für maximale Sicherheit
		SameSite:      sessions.CookieSameSiteLax, // Default: Lax für bessere Subdomain-Kompatibilität
		DevMode:       false,                      // Default: Produktionsmodus
	}

	for _, opt := range opts {
		opt(&config)
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return &Store[T]{
		config: config,
	}, nil
}

// Save persists the session data in a cookie
func (s *Store[T]) Save(w http.ResponseWriter, session *sessions.Session[T]) error {
	encoded, err := sessions.EncodeSession(session, s.config.EncryptionKey)
	if err != nil {
		return err
	}

	cookie := &http.Cookie{
		Name:     session.Name,
		Value:    encoded,
		Path:     sessions.CookiePath,
		Domain:   s.config.Domain,
		MaxAge:   s.config.MaxAge,
		Secure:   s.config.Secure,
		HttpOnly: s.config.HTTPOnly,
		SameSite: getSameSite(s.config.SameSite),
	}

	http.SetCookie(w, cookie)
	return nil
}

// Load retrieves the session data from a cookie
func (s *Store[T]) Load(r *http.Request, name string) (*sessions.Session[T], error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		if err == http.ErrNoCookie {
			return nil, sessions.ErrInvalidSession
		}
		return nil, err
	}

	session, err := sessions.DecodeSession[T](cookie.Value, s.config.EncryptionKey)
	if err != nil {
		return nil, err
	}

	if session.IsExpired() {
		return nil, sessions.ErrSessionExpired
	}

	return session, nil
}

// Destroy removes the session by setting an expired cookie
func (s *Store[T]) Destroy(w http.ResponseWriter, r *http.Request, name string) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     sessions.CookiePath,
		Domain:   s.config.Domain,
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
		Secure:   s.config.Secure,
		HttpOnly: s.config.HTTPOnly,
		SameSite: getSameSite(s.config.SameSite),
	}

	http.SetCookie(w, cookie)
}

// getSameSite converts the SameSite string to http.SameSite
func getSameSite(sameSite string) http.SameSite {
	switch sameSite {
	case sessions.CookieSameSiteStrict:
		return http.SameSiteStrictMode
	case sessions.CookieSameSiteNone:
		return http.SameSiteNoneMode
	default:
		return http.SameSiteLaxMode
	}
}
