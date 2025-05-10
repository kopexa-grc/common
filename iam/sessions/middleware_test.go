// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package sessions

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testStore implements Store for testing
type testStore[T any] struct {
	sessions map[string]*Session[T]
	saved    bool
	loaded   bool
	saveErr  error
}

func newTestStore[T any]() *testStore[T] {
	return &testStore[T]{
		sessions: make(map[string]*Session[T]),
	}
}

func (s *testStore[T]) Save(_ http.ResponseWriter, session *Session[T]) error {
	if s.saveErr != nil {
		return s.saveErr
	}

	s.sessions[session.ID] = session
	s.saved = true

	return nil
}

func (s *testStore[T]) Load(_ *http.Request, _ string) (*Session[T], error) {
	s.loaded = true
	// In a real test, we would parse the cookie and return the session
	// For simplicity, we just return the first session we find
	for _, session := range s.sessions {
		return session, nil
	}

	return nil, nil
}

func (s *testStore[T]) Destroy(_ http.ResponseWriter, _ *http.Request, name string) {
	delete(s.sessions, name)
}

func TestSessionMiddleware(t *testing.T) {
	type testValue struct {
		Key string
	}

	tests := []struct {
		name           string
		setupStore     func() *testStore[testValue]
		setupRequest   func() *http.Request
		handler        http.HandlerFunc
		expectedStatus int
		checkSession   func(t *testing.T, store *testStore[testValue])
		checkLogs      func(t *testing.T, logs *bytes.Buffer)
	}{
		{
			name: "new session is created and saved",
			setupStore: func() *testStore[testValue] {
				return newTestStore[testValue]()
			},
			setupRequest: func() *http.Request {
				return httptest.NewRequest("GET", "/", nil)
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				session := GetSessionFromContext[testValue](r)
				if session == nil {
					session = &Session[testValue]{
						ID:        "test-session",
						Name:      "test",
						Values:    make(map[string]testValue),
						CreatedAt: time.Now(),
					}
					ctx := WithSession(r.Context(), session)
					r = r.WithContext(ctx)
					if rw, ok := w.(*responseWriter); ok {
						rw.HijackRequest(r)
					}
				}
				session.Values["test"] = testValue{Key: "value"}
				w.WriteHeader(http.StatusOK)
			},
			expectedStatus: http.StatusOK,
			checkSession: func(t *testing.T, store *testStore[testValue]) {
				assert.True(t, store.saved, "session should be saved")
				assert.True(t, store.loaded, "session should be loaded")
			},
			checkLogs: func(t *testing.T, logs *bytes.Buffer) {
				assert.Empty(t, logs.String(), "no error logs expected")
			},
		},
		{
			name: "existing session is loaded and updated",
			setupStore: func() *testStore[testValue] {
				store := newTestStore[testValue]()
				store.sessions["test-session"] = &Session[testValue]{
					ID:        "test-session",
					Name:      "test",
					Values:    make(map[string]testValue),
					CreatedAt: time.Now(),
				}
				return store
			},
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("GET", "/", nil)
				req.AddCookie(&http.Cookie{
					Name:  "test",
					Value: "test-session",
				})
				return req
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				session := GetSessionFromContext[testValue](r)
				require.NotNil(t, session, "session should be loaded")
				session.Values["test"] = testValue{Key: "updated"}
				w.WriteHeader(http.StatusOK)
			},
			expectedStatus: http.StatusOK,
			checkSession: func(t *testing.T, store *testStore[testValue]) {
				assert.True(t, store.saved, "session should be saved")
				assert.True(t, store.loaded, "session should be loaded")
				session, ok := store.sessions["test-session"]
				assert.True(t, ok, "session should exist in store")
				value, ok := session.Values["test"]
				assert.True(t, ok, "value should exist")
				assert.Equal(t, "updated", value.Key)
			},
			checkLogs: func(t *testing.T, logs *bytes.Buffer) {
				assert.Empty(t, logs.String(), "no error logs expected")
			},
		},
		{
			name: "error handling in middleware",
			setupStore: func() *testStore[testValue] {
				return newTestStore[testValue]()
			},
			setupRequest: func() *http.Request {
				return httptest.NewRequest("GET", "/", nil)
			},
			handler: func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			expectedStatus: http.StatusInternalServerError,
			checkSession: func(t *testing.T, store *testStore[testValue]) {
				assert.False(t, store.saved, "session should not be saved on error")
				assert.True(t, store.loaded, "session should be loaded")
			},
			checkLogs: func(t *testing.T, logs *bytes.Buffer) {
				assert.Empty(t, logs.String(), "no error logs expected")
			},
		},
		{
			name: "session save error is logged",
			setupStore: func() *testStore[testValue] {
				store := newTestStore[testValue]()
				store.saveErr = assert.AnError
				return store
			},
			setupRequest: func() *http.Request {
				return httptest.NewRequest("GET", "/", nil)
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				session := GetSessionFromContext[testValue](r)
				if session == nil {
					session = &Session[testValue]{
						ID:        "test-session",
						Name:      "test",
						Values:    make(map[string]testValue),
						CreatedAt: time.Now(),
					}
					ctx := WithSession(r.Context(), session)
					r = r.WithContext(ctx)
					if rw, ok := w.(*responseWriter); ok {
						rw.HijackRequest(r)
					}
				}
				session.Values["test"] = testValue{Key: "value"}
				w.WriteHeader(http.StatusOK)
			},
			expectedStatus: http.StatusOK,
			checkSession: func(t *testing.T, store *testStore[testValue]) {
				assert.False(t, store.saved, "session should not be saved due to error")
				assert.True(t, store.loaded, "session should be loaded")
			},
			checkLogs: func(t *testing.T, logs *bytes.Buffer) {
				assert.Contains(t, logs.String(), "failed to save session")
				assert.Contains(t, logs.String(), "test-session")
				assert.Contains(t, logs.String(), assert.AnError.Error())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			store := tt.setupStore()
			req := tt.setupRequest()
			w := httptest.NewRecorder()

			// Setup logging
			var logs bytes.Buffer
			logger := zerolog.New(&logs).With().Timestamp().Logger()
			ctx := logger.WithContext(context.Background())
			req = req.WithContext(ctx)

			// Create middleware
			middleware := SessionMiddleware(store, "test")

			// Execute request
			// Wir müssen den echten responseWriter verwenden, um das Hijack-Pattern zu testen
			rw := &responseWriter{ResponseWriter: w, req: req}
			middleware(tt.handler).ServeHTTP(rw, req)

			// Check response
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Check session state
			tt.checkSession(t, store)

			// Check logs
			tt.checkLogs(t, &logs)
		})
	}
}
