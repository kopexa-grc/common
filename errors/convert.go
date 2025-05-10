// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package errors

import (
	"context"
	"errors"
	"net/http"
)

// FromHTTPStatus converts an HTTP status code to an appropriate error.
func FromHTTPStatus(status int, message string) *Error {
	switch status {
	case http.StatusBadRequest:
		return NewBadRequest(message)
	case http.StatusUnauthorized:
		return NewUnauthorized(message)
	case http.StatusForbidden:
		return NewForbidden(message)
	case http.StatusNotFound:
		return NewNotFound(message)
	case http.StatusConflict:
		return NewConflict(message)
	case http.StatusGone:
		return NewGone(message)
	case http.StatusUnprocessableEntity:
		return NewUnprocessableEntity(message)
	case http.StatusTooManyRequests:
		return New(TooManyRequests, message).WithStatus(status)
	case http.StatusInternalServerError:
		return NewUnexpectedFailure(message)
	case http.StatusNotImplemented:
		return New(NotImplemented, message).WithStatus(status)
	case http.StatusServiceUnavailable:
		return New(ServiceUnavailable, message).WithStatus(status)
	case http.StatusGatewayTimeout:
		return New(GatewayTimeout, message).WithStatus(status)
	default:
		return NewUnexpectedFailure(message)
	}
}

// FromContextError converts a context error to an appropriate error.
func FromContextError(err error) *Error {
	if err == nil {
		return nil
	}

	if errors.Is(err, context.DeadlineExceeded) {
		return New(DeadlineExceeded, "Operation timed out").WithStatus(http.StatusGatewayTimeout)
	}

	if errors.Is(err, context.Canceled) {
		return New(RequestTimeout, "Operation was canceled").WithStatus(http.StatusRequestTimeout)
	}

	return Wrap(err, "Context error")
}

// FromNetworkError converts a network error to an appropriate error.
func FromNetworkError(err error) *Error {
	if err == nil {
		return nil
	}
	// TODO: Implement network error detection logic
	return New(ConnectionFailed, "Network error occurred").With(err)
}

// IsRetryable determines if an error is retryable.
func IsRetryable(err error) bool {
	if e, ok := err.(*Error); ok {
		switch e.Code {
		case ServiceUnavailable, GatewayTimeout, ConnectionFailed,
			ConnectionTimeout, ConnectionRefused, RequestTimeout:
			return true
		}
	}

	return false
}

// IsTimeout determines if an error is a timeout error.
func IsTimeout(err error) bool {
	if e, ok := err.(*Error); ok {
		switch e.Code {
		case DeadlineExceeded, RequestTimeout, GatewayTimeout,
			ConnectionTimeout:
			return true
		}
	}

	return false
}

// IsAuthError determines if an error is an authentication/authorization error.
func IsAuthError(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.Category == CategoryAuth
	}

	return false
}

// IsClientError determines if an error is a client error.
func IsClientError(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.Category == CategoryClient
	}

	return false
}

// IsServerError determines if an error is a server error.
func IsServerError(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.Category == CategoryServer
	}

	return false
}
