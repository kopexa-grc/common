// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package sessions

import (
	"net/http"

	"github.com/rs/zerolog"
)

// contextKey is a private type to avoid context key collisions
// Only used internally for session context
// Comments in English for clarity

type contextKey string

const sessionContextKey contextKey = "session"

// responseWriter wraps http.ResponseWriter to capture the response and hijack the request
// This allows the handler to update the request (e.g. with a new context)
type responseWriter struct {
	http.ResponseWriter
	status int
	req    *http.Request
}

// HijackRequest allows the handler to set the current request (with updated context)
func (rw *responseWriter) HijackRequest(r *http.Request) {
	rw.req = r
}

// WriteHeader captures the status code
func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

// SessionMiddleware returns a middleware that loads the session from the store and
// puts it into the request context for downstream handlers.
func SessionMiddleware[T any](store Store[T], sessionName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Load session from store
			session, err := store.Load(r, sessionName)
			if err == nil && session != nil {
				// Store session in context using type-safe context functions
				ctx := WithSession(r.Context(), session)
				r = r.WithContext(ctx)
			}

			// Create a response wrapper to capture the response and request
			rw := &responseWriter{ResponseWriter: w, req: r}

			// Call the next handler
			next.ServeHTTP(rw, rw.req)

			// Nach dem Handler: aktuelles Request-Objekt verwenden (ggf. vom Handler gehijackt)
			currentReq := rw.req
			if currentReq == nil {
				currentReq = r
			}
			session = GetSessionFromContext[T](currentReq)
			if session != nil {
				if err := store.Save(w, session); err != nil {
					zerolog.Ctx(currentReq.Context()).Error().
						Err(err).
						Str("session_id", session.ID).
						Msg("failed to save session")
				}
			}
		})
	}
}

// GetSessionFromContext extracts the session from the request context using type-safe context functions
func GetSessionFromContext[T any](r *http.Request) *Session[T] {
	sess, _ := FromSession[T](r.Context())
	return sess
}
