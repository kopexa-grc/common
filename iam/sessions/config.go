// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package sessions

import (
	"net/http"
)

// Config is used to configure session management
type Config[T any] struct {
	// Store is responsible for managing the session cookies
	Store Store[T]
	// CookieConfig contains the cookie settings for sessions
	CookieConfig *CookieConfig
}

// CookieConfig contains the cookie settings for sessions
type CookieConfig struct {
	// Name is the name of the session cookie
	Name string
	// Domain is the domain for the cookie
	Domain string
	// MaxAge is the maximum age of the session cookie in seconds
	MaxAge int
	// Secure determines if the cookie should only be sent over HTTPS
	Secure bool
	// HTTPOnly determines if the cookie should be accessible via JavaScript
	HTTPOnly bool
	// SameSite determines the SameSite attribute of the cookie
	SameSite http.SameSite
}

// Option allows users to optionally supply configuration to the session middleware
type Option[T any] func(*Config[T])

// NewConfig creates a new session config with options
func NewConfig[T any](store Store[T], opts ...Option[T]) Config[T] {
	c := Config[T]{
		Store: store,
	}

	for _, opt := range opts {
		opt(&c)
	}

	return c
}

// WithCookieConfig allows the user to specify cookie settings
func WithCookieConfig[T any](config *CookieConfig) Option[T] {
	return func(c *Config[T]) {
		c.CookieConfig = config
	}
}

// WithMaxAge allows the user to specify the maximum age for the session cookie
func WithMaxAge[T any](maxAge int) Option[T] {
	return func(c *Config[T]) {
		if c.CookieConfig == nil {
			c.CookieConfig = &CookieConfig{}
		}

		c.CookieConfig.MaxAge = maxAge
	}
}

// WithSecure allows the user to specify if the cookie should only be sent over HTTPS
func WithSecure[T any](secure bool) Option[T] {
	return func(c *Config[T]) {
		if c.CookieConfig == nil {
			c.CookieConfig = &CookieConfig{}
		}

		c.CookieConfig.Secure = secure
	}
}

// WithHTTPOnly allows the user to specify if the cookie should be accessible via JavaScript
func WithHTTPOnly[T any](httpOnly bool) Option[T] {
	return func(c *Config[T]) {
		if c.CookieConfig == nil {
			c.CookieConfig = &CookieConfig{}
		}

		c.CookieConfig.HTTPOnly = httpOnly
	}
}

// WithSameSite allows the user to specify the SameSite attribute of the cookie
func WithSameSite[T any](sameSite http.SameSite) Option[T] {
	return func(c *Config[T]) {
		if c.CookieConfig == nil {
			c.CookieConfig = &CookieConfig{}
		}

		c.CookieConfig.SameSite = sameSite
	}
}

// WithDomain allows the user to specify the domain for the cookie
func WithDomain[T any](domain string) Option[T] {
	return func(c *Config[T]) {
		if c.CookieConfig == nil {
			c.CookieConfig = &CookieConfig{}
		}

		c.CookieConfig.Domain = domain
	}
}
