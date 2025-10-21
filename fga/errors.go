// Original Licenses under Apache-2.0 by the openlane https://github.com/theopenlane
// SPDX-License-Identifier: Apache-2.0

package fga

import (
	"errors"
	"fmt"
)

// Common errors that can occur when using the FGA client.
var (
	// ErrInvalidArgument is returned when an invalid argument is provided to a function.
	ErrInvalidArgument = errors.New("invalid argument")
	// ErrUnauthorized is returned when the client is not authorized to perform an operation.
	ErrUnauthorized = errors.New("unauthorized")
	// ErrNotFound is returned when a requested resource is not found.
	ErrNotFound      = errors.New("not found")
	ErrEmptyResponse = errors.New("empty response from server")
	ErrInvalidEntity = errors.New("invalid entity")
	// ErrEmptyBatchCheckResponse is returned when a batch check operation returns an empty response.
	// This indicates that the FGA service did not return any results for the batch check request.
	ErrEmptyBatchCheckResponse = errors.New("empty response from batch check")
	// ErrFailedToTransformModel is returned when the model transformation fails
	ErrFailedToTransformModel = errors.New("failed to transform model")
)

// WriteError represents an error that occurred during a write operation to the FGA service.
// It contains details about the operation that failed and the specific error response.
// This error type is returned when a write operation (grant or revoke) fails.
type WriteError struct {
	// User is the user identifier involved in the failed operation
	User string
	// Relation is the relation that was being modified
	Relation string
	// Object is the object identifier involved in the failed operation
	Object string
	// Operation indicates whether this was a write or delete operation
	Operation string
	// ErrorResponse contains the original error from the FGA service
	ErrorResponse error
}

// Error implements the error interface for WriteError.
// Returns a formatted string containing the operation details and error message.
// The format is: "FGA <operation> error: <user> <relation> <object> – <error>"
func (e *WriteError) Error() string {
	return fmt.Sprintf("FGA %s error: %s %s %s – %v",
		e.Operation, e.User, e.Relation, e.Object, e.ErrorResponse)
}

// Unwrap returns the underlying error that caused the WriteError.
// This allows for error type assertions and unwrapping of the original error.
func (e *WriteError) Unwrap() error {
	return e.ErrorResponse
}

// IsWriteError checks if an error is a WriteError.
// This is a convenience function for error type checking.
func IsWriteError(err error) bool {
	var writeErr *WriteError
	return errors.As(err, &writeErr)
}
