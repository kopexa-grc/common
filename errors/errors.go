// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package errors

import (
	"fmt"
	"net/http"
	"time"
)

// Error is our internal error type.
// swagger:model
type Error struct {
	// Code is the error code that identifies the type of error.
	Code ErrorCode `json:"code" validate:"required"`
	// Category is the category of the error (e.g., client, server, auth).
	Category ErrorCategory `json:"category" validate:"required"`
	// Status is the HTTP status code associated with the error.
	Status int `json:"status" validate:"required"`
	// Message is a human-readable error message.
	Message string `json:"message" validate:"required"`
	// Entity is the entity that the error is related to (e.g., "user", "document").
	Entity string `json:"entity,omitempty"`
	// RequestID is the ID of the request that caused the error.
	RequestID string `json:"request_id,omitempty"`
	// Timestamp is when the error occurred.
	Timestamp time.Time `json:"timestamp"`
	// Details contains additional error details.
	Details map[string]interface{} `json:"details,omitempty"`
	// Err is the underlying error.
	Err error `json:"-"`
}

// New creates a new Error.
func New(code ErrorCode, message string) *Error {
	return &Error{
		Code:      code,
		Category:  getCategoryForCode(code),
		Message:   message,
		Timestamp: time.Now(),
		Details:   make(map[string]interface{}),
	}
}

// getCategoryForCode returns the appropriate category for a given error code.
func getCategoryForCode(code ErrorCode) ErrorCategory {
	switch code {
	case BadRequest, Unauthorized, Forbidden, NotFound, Conflict, Gone,
		UnprocessableEntity, TooManyRequests:
		return CategoryClient
	case UnexpectedFailure, NotImplemented, ServiceUnavailable, GatewayTimeout:
		return CategoryServer
	case ResourceExhausted, QuotaExceeded, SpaceNotFound:
		return CategoryResource
	case NoAuthorization, InvalidCredentials, TokenExpired:
		return CategoryAuth
	case ConnectionFailed, ConnectionTimeout, ConnectionRefused:
		return CategoryNetwork
	case DeadlineExceeded, RequestTimeout:
		return CategoryTimeout
	default:
		return CategoryClient
	}
}

// Newf creates a new Error with formatted message and underlying error.
func Newf(code ErrorCode, err error, format string, args ...any) *Error {
	return New(code, fmt.Sprintf(format, args...)).With(err)
}

// WithStatus sets the HTTP status code for the Error.
func (e *Error) WithStatus(status int) *Error {
	e.Status = status
	return e
}

// WithCode sets the error code for the Error.
func (e *Error) WithCode(code ErrorCode) *Error {
	e.Code = code
	return e
}

// WithMessage sets the error message for the Error.
func (e *Error) WithMessage(message string) *Error {
	e.Message = message
	return e
}

// WithEntity sets the entity for the Error.
func (e *Error) WithEntity(entity string) *Error {
	e.Entity = entity
	return e
}

// WithRequestID sets the request ID for the Error.
func (e *Error) WithRequestID(requestID string) *Error {
	e.RequestID = requestID
	return e
}

// WithDetails adds additional details to the Error.
func (e *Error) WithDetails(key string, value interface{}) *Error {
	if e.Details == nil {
		e.Details = make(map[string]interface{})
	}
	e.Details[key] = value
	return e
}

// Error returns the error message.
func (e *Error) Error() string {
	return e.Message
}

// Unwrap returns the underlying error.
func (e *Error) Unwrap() error {
	return e.Err
}

// With adds an underlying error to the Error.
func (e *Error) With(err error) *Error {
	e.Err = err
	return e
}

// Wrap wraps an error with a message.
func Wrap(err error, message string) *Error {
	return &Error{
		Code:    UnexpectedFailure,
		Status:  http.StatusInternalServerError,
		Message: message,
		Err:     err,
	}
}

// IsError checks if the error is an Error.
func IsError(err error) bool {
	_, ok := err.(*Error)
	return ok
}

// Is checks if the error is an Error and if the code matches.
func Is(err error, code ErrorCode) bool {
	if e, ok := err.(*Error); ok {
		return e.Code == code
	}
	return false
}

// NewGone creates a new Error with the Gone code.
func NewGone(message string) *Error {
	if message == "" {
		message = msgGone
	}
	return New(Gone, message).WithStatus(http.StatusGone)
}

// NewUnexpectedFailure creates a new Error with the UnexpectedFailure code.
func NewUnexpectedFailure(message string) *Error {
	if message == "" {
		message = msgUnexpectedFailure
	}
	return New(UnexpectedFailure, message).WithStatus(http.StatusInternalServerError)
}

// IsUnexpectedFailure checks if the error is an Error with the UnexpectedFailure code.
func IsUnexpectedFailure(err error) bool {
	return Is(err, UnexpectedFailure)
}

// NewUnauthorized creates a new Error with the Unauthorized code.
func NewUnauthorized(message string) *Error {
	if message == "" {
		message = msgUnauthorized
	}
	return New(Unauthorized, message).WithStatus(http.StatusUnauthorized)
}

// IsUnauthorized checks if the error is an Error with the Unauthorized code.
func IsUnauthorized(err error) bool {
	return Is(err, Unauthorized)
}

// NewUnprocessableEntity creates a new Error with the UnprocessableEntity code.
func NewUnprocessableEntity(message string) *Error {
	if message == "" {
		message = msgUnprocessableEntity
	}
	return New(UnprocessableEntity, message).WithStatus(http.StatusUnprocessableEntity)
}

// NewBadRequest creates a new Error with the BadRequest code.
func NewBadRequest(message string) *Error {
	if message == "" {
		message = msgBadRequest
	}
	return New(BadRequest, message).WithStatus(http.StatusBadRequest)
}

// NewConflict creates a new Error with the Conflict code.
func NewConflict(message string) *Error {
	if message == "" {
		message = msgConflict
	}
	return New(Conflict, message).WithStatus(http.StatusConflict)
}

// IsBadRequest checks if the error is an Error with the BadRequest code.
func IsBadRequest(err error) bool {
	return Is(err, BadRequest)
}

// NewNotFound creates a new Error with the NotFound code.
func NewNotFound(message string) *Error {
	if message == "" {
		message = msgNotFound
	}
	return New(NotFound, message).WithStatus(http.StatusNotFound)
}

// IsNotFound checks if the error is an Error with the NotFound code.
func IsNotFound(err error) bool {
	return Is(err, NotFound)
}

// NewForbidden creates a new Error with the Forbidden code.
func NewForbidden(message string) *Error {
	if message == "" {
		message = msgForbidden
	}
	return New(Forbidden, message).WithStatus(http.StatusForbidden)
}

// IsForbidden checks if the error is an Error with the Forbidden code.
func IsForbidden(err error) bool {
	return Is(err, Forbidden)
}

// NewInvalidArgument creates a new Error with the InvalidArgument code.
func NewInvalidArgument(message string) *Error {
	if message == "" {
		message = msgInvalidArgument
	}
	return New(InvalidArgument, message).WithStatus(http.StatusBadRequest)
}

// IsInvalidArgument checks if the error is an Error with the InvalidArgument code.
func IsInvalidArgument(err error) bool {
	return Is(err, InvalidArgument)
}

// NewFailedPrecondition creates a new Error with the FailedPrecondition code.
func NewFailedPrecondition(message string) *Error {
	if message == "" {
		message = msgFailedPrecondition
	}
	return New(FailedPrecondition, message).WithStatus(http.StatusPreconditionFailed)
}

// IsFailedPrecondition checks if the error is an Error with the FailedPrecondition code.
func IsFailedPrecondition(err error) bool {
	return Is(err, FailedPrecondition)
}

// NewTooManyRequests creates a new error with the TooManyRequests code
func NewTooManyRequests(msg string) error {
	if msg == "" {
		msg = msgTooManyRequests
	}
	return New(TooManyRequests, msg).WithStatus(http.StatusTooManyRequests)
}

// IsTooManyRequests checks if the error is a TooManyRequests error
func IsTooManyRequests(err error) bool {
	return Is(err, TooManyRequests)
}

// NewNotImplemented creates a new error with the NotImplemented code
func NewNotImplemented(msg string) error {
	if msg == "" {
		msg = msgNotImplemented
	}
	return New(NotImplemented, msg).WithStatus(http.StatusNotImplemented)
}

// IsNotImplemented checks if the error is a NotImplemented error
func IsNotImplemented(err error) bool {
	return Is(err, NotImplemented)
}

// NewServiceUnavailable creates a new error with the ServiceUnavailable code
func NewServiceUnavailable(msg string) error {
	if msg == "" {
		msg = msgServiceUnavailable
	}
	return New(ServiceUnavailable, msg).WithStatus(http.StatusServiceUnavailable)
}

// IsServiceUnavailable checks if the error is a ServiceUnavailable error
func IsServiceUnavailable(err error) bool {
	return Is(err, ServiceUnavailable)
}

// NewGatewayTimeout creates a new error with the GatewayTimeout code
func NewGatewayTimeout(msg string) error {
	if msg == "" {
		msg = msgGatewayTimeout
	}
	return New(GatewayTimeout, msg).WithStatus(http.StatusGatewayTimeout)
}

// IsGatewayTimeout checks if the error is a GatewayTimeout error
func IsGatewayTimeout(err error) bool {
	return Is(err, GatewayTimeout)
}

// NewResourceExhausted creates a new error with the ResourceExhausted code
func NewResourceExhausted(msg string) error {
	if msg == "" {
		msg = msgResourceExhausted
	}
	return New(ResourceExhausted, msg).WithStatus(http.StatusInsufficientStorage)
}

// IsResourceExhausted checks if the error is a ResourceExhausted error
func IsResourceExhausted(err error) bool {
	return Is(err, ResourceExhausted)
}

// NewQuotaExceeded creates a new error with the QuotaExceeded code
func NewQuotaExceeded(msg string) error {
	if msg == "" {
		msg = msgQuotaExceeded
	}
	return New(QuotaExceeded, msg).WithStatus(http.StatusTooManyRequests)
}

// IsQuotaExceeded checks if the error is a QuotaExceeded error
func IsQuotaExceeded(err error) bool {
	return Is(err, QuotaExceeded)
}

// NewInvalidCredentials creates a new error with the InvalidCredentials code
func NewInvalidCredentials(msg string) error {
	if msg == "" {
		msg = msgInvalidCredentials
	}
	return New(InvalidCredentials, msg).WithStatus(http.StatusUnauthorized)
}

// IsInvalidCredentials checks if the error is an InvalidCredentials error
func IsInvalidCredentials(err error) bool {
	return Is(err, InvalidCredentials)
}

// NewTokenExpired creates a new error with the TokenExpired code
func NewTokenExpired(msg string) error {
	if msg == "" {
		msg = msgTokenExpired
	}
	return New(TokenExpired, msg).WithStatus(http.StatusUnauthorized)
}

// IsTokenExpired checks if the error is a TokenExpired error
func IsTokenExpired(err error) bool {
	return Is(err, TokenExpired)
}

// NewConnectionFailed creates a new error with the ConnectionFailed code
func NewConnectionFailed(msg string) error {
	if msg == "" {
		msg = msgConnectionFailed
	}
	return New(ConnectionFailed, msg).WithStatus(http.StatusServiceUnavailable)
}

// IsConnectionFailed checks if the error is a ConnectionFailed error
func IsConnectionFailed(err error) bool {
	return Is(err, ConnectionFailed)
}

// NewConnectionTimeout creates a new error with the ConnectionTimeout code
func NewConnectionTimeout(msg string) error {
	if msg == "" {
		msg = msgConnectionTimeout
	}
	return New(ConnectionTimeout, msg).WithStatus(http.StatusGatewayTimeout)
}

// IsConnectionTimeout checks if the error is a ConnectionTimeout error
func IsConnectionTimeout(err error) bool {
	return Is(err, ConnectionTimeout)
}

// NewConnectionRefused creates a new error with the ConnectionRefused code
func NewConnectionRefused(msg string) error {
	if msg == "" {
		msg = msgConnectionRefused
	}
	return New(ConnectionRefused, msg).WithStatus(http.StatusServiceUnavailable)
}

// IsConnectionRefused checks if the error is a ConnectionRefused error
func IsConnectionRefused(err error) bool {
	return Is(err, ConnectionRefused)
}

// NewDeadlineExceeded creates a new error with the DeadlineExceeded code
func NewDeadlineExceeded(msg string) error {
	if msg == "" {
		msg = msgDeadlineExceeded
	}
	return New(DeadlineExceeded, msg).WithStatus(http.StatusGatewayTimeout)
}

// IsDeadlineExceeded checks if the error is a DeadlineExceeded error
func IsDeadlineExceeded(err error) bool {
	return Is(err, DeadlineExceeded)
}

// NewRequestTimeout creates a new error with the RequestTimeout code
func NewRequestTimeout(msg string) error {
	if msg == "" {
		msg = msgRequestTimeout
	}
	return New(RequestTimeout, msg).WithStatus(http.StatusRequestTimeout)
}

// IsRequestTimeout checks if the error is a RequestTimeout error
func IsRequestTimeout(err error) bool {
	return Is(err, RequestTimeout)
}

// NewOutOfRange creates a new error with the OutOfRange code
func NewOutOfRange(msg string) error {
	if msg == "" {
		msg = msgOutOfRange
	}
	return New(OutOfRange, msg).WithStatus(http.StatusBadRequest)
}

// IsOutOfRange checks if the error is an OutOfRange error
func IsOutOfRange(err error) bool {
	return Is(err, OutOfRange)
}

// Code returns the error code for the given error
func Code(err error) ErrorCode {
	if err == nil {
		return ""
	}
	if e, ok := err.(*Error); ok {
		return e.Code
	}
	return UnexpectedFailure
}

// Status returns the HTTP status code for the given error
func Status(err error) int {
	if err == nil {
		return http.StatusOK
	}
	if e, ok := err.(*Error); ok {
		return e.Status
	}
	return http.StatusInternalServerError
}
