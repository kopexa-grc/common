// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package router_test

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/kopexa-grc/common/khttp/metric"
	"github.com/kopexa-grc/common/khttp/router"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testMiddleware router.Middleware = func(next http.Handler) http.Handler {
	return nil
}

func TestMiddlewareRegistry_Register(t *testing.T) {
	t.Run("register single middleware", func(t *testing.T) {
		testRouter := chi.NewRouter()

		testSubject := router.NewMiddlewareRegistry()

		testSubject.Register(testMiddleware)

		testSubject.UseOnRouter(testRouter)
		assert.True(t, containsMiddleware(testRouter, testMiddleware))
	})

	t.Run("register middleware twice", func(t *testing.T) {
		testRouter := chi.NewRouter()

		testSubject := router.NewMiddlewareRegistry()
		testSubject.Register(testMiddleware)
		testSubject.Register(testMiddleware)

		testSubject.UseOnRouter(testRouter)
		assert.True(t, containsMiddleware(testRouter, testMiddleware))
		assert.Len(t, testRouter.Middlewares(), 1)
	})

	t.Run("register multiple middlewares", func(t *testing.T) {
		testRouter := chi.NewRouter()

		testSubject := router.NewMiddlewareRegistry()

		testSubject.Register(testMiddleware)
		testSubject.Register(metric.Middleware)

		testSubject.UseOnRouter(testRouter)
		assert.True(t, containsMiddleware(testRouter, testMiddleware))
		assert.True(t, containsMiddleware(testRouter, metric.Middleware))
		assert.Len(t, testRouter.Middlewares(), 2)
	})
}

func TestMiddlewareRegistry_Replace(t *testing.T) {
	t.Run("replace existing middleware", func(t *testing.T) {
		testRouter := chi.NewRouter()

		testSubject := router.NewMiddlewareRegistry()
		testSubject.Register(testMiddleware)

		err := testSubject.Replace(testMiddleware, metric.Middleware)
		require.NoError(t, err)

		testSubject.UseOnRouter(testRouter)
		assert.True(t, containsMiddleware(testRouter, metric.Middleware))
		assert.Len(t, testRouter.Middlewares(), 1)
	})

	t.Run("replace non-existing middleware", func(t *testing.T) {
		testRouter := chi.NewRouter()

		testSubject := router.NewMiddlewareRegistry()
		testSubject.Register(testMiddleware)

		err := testSubject.Replace(metric.Middleware, testMiddleware)
		require.Error(t, err)

		testSubject.UseOnRouter(testRouter)
		assert.True(t, containsMiddleware(testRouter, testMiddleware))
		assert.Len(t, testRouter.Middlewares(), 1)
	})

	t.Run("replace middleware with itself", func(t *testing.T) {
		testRouter := chi.NewRouter()

		testSubject := router.NewMiddlewareRegistry()
		testSubject.Register(testMiddleware)

		err := testSubject.Replace(testMiddleware, testMiddleware)
		require.NoError(t, err)

		testSubject.UseOnRouter(testRouter)
		assert.True(t, containsMiddleware(testRouter, testMiddleware))
		assert.Len(t, testRouter.Middlewares(), 1)
	})
}

func TestMiddlewareRegistry_Unregister(t *testing.T) {
	t.Run("unregister existing middleware", func(t *testing.T) {
		testRouter := chi.NewRouter()

		testSubject := router.NewMiddlewareRegistry()
		testSubject.Register(testMiddleware)

		testSubject.Unregister(testMiddleware)

		testSubject.UseOnRouter(testRouter)
		assert.False(t, containsMiddleware(testRouter, testMiddleware))
		assert.Empty(t, testRouter.Middlewares())
	})

	t.Run("unregister non-existing middleware", func(t *testing.T) {
		testRouter := chi.NewRouter()

		testSubject := router.NewMiddlewareRegistry()
		testSubject.Register(testMiddleware)

		testSubject.Unregister(metric.Middleware)

		testSubject.UseOnRouter(testRouter)
		assert.True(t, containsMiddleware(testRouter, testMiddleware))
		assert.Len(t, testRouter.Middlewares(), 1)
	})

	t.Run("unregister middleware from empty registry", func(t *testing.T) {
		testRouter := chi.NewRouter()

		testSubject := router.NewMiddlewareRegistry()

		testSubject.Unregister(testMiddleware)

		testSubject.UseOnRouter(testRouter)

		assert.Empty(t, testRouter.Middlewares())
	})
}

func TestMiddlewareRegistry_UseOnRouter(t *testing.T) {
	t.Run("use on router", func(t *testing.T) {
		testRouter := chi.NewRouter()

		testSubject := router.NewMiddlewareRegistry()
		testSubject.Register(testMiddleware)

		testSubject.UseOnRouter(testRouter)
		assert.True(t, containsMiddleware(testRouter, testMiddleware))
		assert.Len(t, testRouter.Middlewares(), 1)
	})

	t.Run("use on router empty registry", func(t *testing.T) {
		testRouter := chi.NewRouter()

		testSubject := router.NewMiddlewareRegistry()

		testSubject.UseOnRouter(testRouter)

		assert.Empty(t, testRouter.Middlewares())
	})
}

func containsMiddleware(router chi.Router, middleware router.Middleware) bool {
	return lo.ContainsBy(router.Middlewares(), func(routerMiddleware func(http.Handler) http.Handler) bool {
		return reflect.ValueOf(routerMiddleware).Pointer() == reflect.ValueOf(middleware).Pointer()
	})
}
