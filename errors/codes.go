// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package errors

// ErrorCode represents a type of error that can occur in the system.
type ErrorCode string

// ErrorCategory represents the category of an error.
type ErrorCategory string

const (
	// Error categories
	CategoryClient   ErrorCategory = "client"   // Client-side errors (4xx)
	CategoryServer   ErrorCategory = "server"   // Server-side errors (5xx)
	CategoryResource ErrorCategory = "resource" // Resource-related errors
	CategoryAuth     ErrorCategory = "auth"     // Authentication/Authorization errors
	CategoryNetwork  ErrorCategory = "network"  // Network-related errors
	CategoryTimeout  ErrorCategory = "timeout"  // Timeout-related errors

	// Client errors (4xx)
	BadRequest          ErrorCode = "BAD_REQUEST"          // 400
	Unauthorized        ErrorCode = "UNAUTHORIZED"         // 401
	Forbidden           ErrorCode = "FORBIDDEN"            // 403
	NotFound            ErrorCode = "NOT_FOUND"            // 404
	Conflict            ErrorCode = "CONFLICT"             // 409
	Gone                ErrorCode = "GONE"                 // 410
	UnprocessableEntity ErrorCode = "UNPROCESSABLE_ENTITY" // 422
	TooManyRequests     ErrorCode = "TOO_MANY_REQUESTS"    // 429

	// Server errors (5xx)
	UnexpectedFailure  ErrorCode = "UNEXPECTED_FAILURE"  // 500
	NotImplemented     ErrorCode = "NOT_IMPLEMENTED"     // 501
	ServiceUnavailable ErrorCode = "SERVICE_UNAVAILABLE" // 503
	GatewayTimeout     ErrorCode = "GATEWAY_TIMEOUT"     // 504

	// Resource errors
	ResourceExhausted ErrorCode = "RESOURCE_EXHAUSTED"
	QuotaExceeded     ErrorCode = "QUOTA_EXCEEDED"
	SpaceNotFound     ErrorCode = "SPACE_NOT_FOUND"

	// Auth errors
	NoAuthorization    ErrorCode = "no_authorization"
	InvalidCredentials ErrorCode = "INVALID_CREDENTIALS"
	TokenExpired       ErrorCode = "TOKEN_EXPIRED"

	// Network errors
	ConnectionFailed  ErrorCode = "CONNECTION_FAILED"
	ConnectionTimeout ErrorCode = "CONNECTION_TIMEOUT"
	ConnectionRefused ErrorCode = "CONNECTION_REFUSED"

	// Timeout errors
	DeadlineExceeded ErrorCode = "DEADLINE_EXCEEDED"
	RequestTimeout   ErrorCode = "REQUEST_TIMEOUT"

	// Validation errors
	InvalidArgument    ErrorCode = "INVALID_ARGUMENT"
	FailedPrecondition ErrorCode = "FAILED_PRECONDITION"
	OutOfRange         ErrorCode = "OUT_OF_RANGE"
)
