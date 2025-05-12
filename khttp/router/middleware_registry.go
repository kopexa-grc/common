// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package router

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/go-chi/chi/v5"
	"github.com/zyedidia/generic/list"
)

// Middleware is a function that wraps a http.Handler.
type Middleware = func(http.Handler) http.Handler

// MiddlewareRegistry is a registry for middlewares.
type MiddlewareRegistry struct {
	middlewares *list.List[Middleware]
}

// NewMiddlewareRegistry creates a new MiddlewareRegistry.
func NewMiddlewareRegistry() MiddlewareRegistry {
	return MiddlewareRegistry{
		middlewares: list.New[Middleware](),
	}
}

// Register adds a middleware to the registry.
func (registry MiddlewareRegistry) Register(middleware Middleware) {
	if registry.findFirst(middleware) != nil {
		return
	}
	registry.middlewares.PushBack(middleware)
}

// UseOnRouter adds all middlewares in the registry to the given router.
func (registry MiddlewareRegistry) UseOnRouter(router chi.Router) {
	registry.middlewares.Front.Each(func(middleware Middleware) {
		router.Use(middleware)
	})
}

// Replace replaces the first occurrence of oldMiddleware with newMiddleware.
func (registry MiddlewareRegistry) Replace(oldMiddleware Middleware, newMiddleware Middleware) error {
	nodeToReplace := registry.findFirst(oldMiddleware)
	if nodeToReplace == nil {
		return fmt.Errorf("middleware not found")
	}
	nodeToReplace.Value = newMiddleware
	return nil
}

// Unregister removes the first occurrence of middleware from the registry.
func (registry MiddlewareRegistry) Unregister(middleware Middleware) {
	nodeToRemove := registry.findFirst(middleware)
	if nodeToRemove != nil {
		registry.middlewares.Remove(nodeToRemove)
	}
}

func (registry MiddlewareRegistry) findFirst(middleware Middleware) *list.Node[Middleware] {
	for node := registry.middlewares.Front; node != nil; node = node.Next {
		if reflect.ValueOf(node.Value).Pointer() == reflect.ValueOf(middleware).Pointer() {
			return node
		}
	}

	return nil
}
