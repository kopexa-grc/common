// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package errors

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	err := New(BadRequest, "test error")
	if err.Code != BadRequest {
		t.Errorf("New() code = %v, want %v", err.Code, BadRequest)
	}
	if err.Message != "test error" {
		t.Errorf("New() message = %v, want %v", err.Message, "test error")
	}
}

func TestNewf(t *testing.T) {
	underlying := errors.New("underlying error")
	err := Newf(BadRequest, underlying, "test error: %v", "formatted")
	if err.Code != BadRequest {
		t.Errorf("Newf() code = %v, want %v", err.Code, BadRequest)
	}
	if err.Message != "test error: formatted" {
		t.Errorf("Newf() message = %v, want %v", err.Message, "test error: formatted")
	}
	if err.Err != underlying {
		t.Errorf("Newf() err = %v, want %v", err.Err, underlying)
	}
}

func TestWithStatus(t *testing.T) {
	err := New(BadRequest, "test error").WithStatus(http.StatusBadRequest)
	if err.Status != http.StatusBadRequest {
		t.Errorf("WithStatus() = %v, want %v", err.Status, http.StatusBadRequest)
	}
}

func TestWithCode(t *testing.T) {
	err := New(BadRequest, "test error").WithCode(NotFound)
	if err.Code != NotFound {
		t.Errorf("WithCode() = %v, want %v", err.Code, NotFound)
	}
}

func TestWithMessage(t *testing.T) {
	err := New(BadRequest, "test error").WithMessage("new message")
	if err.Message != "new message" {
		t.Errorf("WithMessage() = %v, want %v", err.Message, "new message")
	}
}

func TestWithEntity(t *testing.T) {
	err := New(BadRequest, "test error").WithEntity("user")
	if err.Entity != "user" {
		t.Errorf("WithEntity() = %v, want %v", err.Entity, "user")
	}
}

func TestError(t *testing.T) {
	err := New(BadRequest, "test error")
	if err.Error() != "test error" {
		t.Errorf("Error() = %v, want %v", err.Error(), "test error")
	}
}

func TestUnwrap(t *testing.T) {
	underlying := errors.New("underlying error")
	err := New(BadRequest, "test error").With(underlying)
	if err.Unwrap() != underlying {
		t.Errorf("Unwrap() = %v, want %v", err.Unwrap(), underlying)
	}
}

func TestWith(t *testing.T) {
	underlying := errors.New("underlying error")
	err := New(BadRequest, "test error").With(underlying)
	if err.Err != underlying {
		t.Errorf("With() = %v, want %v", err.Err, underlying)
	}
}

func TestWrap(t *testing.T) {
	underlying := errors.New("underlying error")
	err := Wrap(underlying, "wrapped error")
	if err.Code != UnexpectedFailure {
		t.Errorf("Wrap() code = %v, want %v", err.Code, UnexpectedFailure)
	}
	if err.Status != http.StatusInternalServerError {
		t.Errorf("Wrap() status = %v, want %v", err.Status, http.StatusInternalServerError)
	}
	if err.Message != "wrapped error" {
		t.Errorf("Wrap() message = %v, want %v", err.Message, "wrapped error")
	}
	if err.Err != underlying {
		t.Errorf("Wrap() err = %v, want %v", err.Err, underlying)
	}
}

func TestIsError(t *testing.T) {
	err := New(BadRequest, "test error")
	if !IsError(err) {
		t.Error("IsError() = false, want true")
	}
	if IsError(errors.New("test error")) {
		t.Error("IsError() = true, want false")
	}
}

func TestIs(t *testing.T) {
	err := New(BadRequest, "test error")
	if !Is(err, BadRequest) {
		t.Error("Is() = false, want true")
	}
	if Is(err, NotFound) {
		t.Error("Is() = true, want false")
	}
	if Is(errors.New("test error"), BadRequest) {
		t.Error("Is() = true, want false")
	}
}

func TestNewGone(t *testing.T) {
	err := NewGone("")
	if err.Code != Gone {
		t.Errorf("NewGone() code = %v, want %v", err.Code, Gone)
	}
	if err.Status != http.StatusGone {
		t.Errorf("NewGone() status = %v, want %v", err.Status, http.StatusGone)
	}
	if err.Message != msgGone {
		t.Errorf("NewGone() message = %v, want %v", err.Message, msgGone)
	}
}

func TestNewUnexpectedFailure(t *testing.T) {
	err := NewUnexpectedFailure("")
	if err.Code != UnexpectedFailure {
		t.Errorf("NewUnexpectedFailure() code = %v, want %v", err.Code, UnexpectedFailure)
	}
	if err.Status != http.StatusInternalServerError {
		t.Errorf("NewUnexpectedFailure() status = %v, want %v", err.Status, http.StatusInternalServerError)
	}
	if err.Message != msgUnexpectedFailure {
		t.Errorf("NewUnexpectedFailure() message = %v, want %v", err.Message, msgUnexpectedFailure)
	}
}

func TestIsUnexpectedFailure(t *testing.T) {
	err := NewUnexpectedFailure("test error")
	if !IsUnexpectedFailure(err) {
		t.Error("IsUnexpectedFailure() = false, want true")
	}
	if IsUnexpectedFailure(errors.New("test error")) {
		t.Error("IsUnexpectedFailure() = true, want false")
	}
}

func TestNewUnauthorized(t *testing.T) {
	err := NewUnauthorized("")
	if err.Code != Unauthorized {
		t.Errorf("NewUnauthorized() code = %v, want %v", err.Code, Unauthorized)
	}
	if err.Status != http.StatusUnauthorized {
		t.Errorf("NewUnauthorized() status = %v, want %v", err.Status, http.StatusUnauthorized)
	}
	if err.Message != msgUnauthorized {
		t.Errorf("NewUnauthorized() message = %v, want %v", err.Message, msgUnauthorized)
	}
}

func TestIsUnauthorized(t *testing.T) {
	err := NewUnauthorized("test error")
	if !IsUnauthorized(err) {
		t.Error("IsUnauthorized() = false, want true")
	}
	if IsUnauthorized(errors.New("test error")) {
		t.Error("IsUnauthorized() = true, want false")
	}
}

func TestNewUnprocessableEntity(t *testing.T) {
	err := NewUnprocessableEntity("")
	if err.Code != UnprocessableEntity {
		t.Errorf("NewUnprocessableEntity() code = %v, want %v", err.Code, UnprocessableEntity)
	}
	if err.Status != http.StatusUnprocessableEntity {
		t.Errorf("NewUnprocessableEntity() status = %v, want %v", err.Status, http.StatusUnprocessableEntity)
	}
	if err.Message != msgUnprocessableEntity {
		t.Errorf("NewUnprocessableEntity() message = %v, want %v", err.Message, msgUnprocessableEntity)
	}
}

func TestNewBadRequest(t *testing.T) {
	err := NewBadRequest("")
	if err.Code != BadRequest {
		t.Errorf("NewBadRequest() code = %v, want %v", err.Code, BadRequest)
	}
	if err.Status != http.StatusBadRequest {
		t.Errorf("NewBadRequest() status = %v, want %v", err.Status, http.StatusBadRequest)
	}
	if err.Message != msgBadRequest {
		t.Errorf("NewBadRequest() message = %v, want %v", err.Message, msgBadRequest)
	}
}

func TestNewConflict(t *testing.T) {
	err := NewConflict("")
	if err.Code != Conflict {
		t.Errorf("NewConflict() code = %v, want %v", err.Code, Conflict)
	}
	if err.Status != http.StatusConflict {
		t.Errorf("NewConflict() status = %v, want %v", err.Status, http.StatusConflict)
	}
	if err.Message != msgConflict {
		t.Errorf("NewConflict() message = %v, want %v", err.Message, msgConflict)
	}
}

func TestIsBadRequest(t *testing.T) {
	err := NewBadRequest("test error")
	if !IsBadRequest(err) {
		t.Error("IsBadRequest() = false, want true")
	}
	if IsBadRequest(errors.New("test error")) {
		t.Error("IsBadRequest() = true, want false")
	}
}

func TestNewNotFound(t *testing.T) {
	err := NewNotFound("")
	if err.Code != NotFound {
		t.Errorf("NewNotFound() code = %v, want %v", err.Code, NotFound)
	}
	if err.Status != http.StatusNotFound {
		t.Errorf("NewNotFound() status = %v, want %v", err.Status, http.StatusNotFound)
	}
	if err.Message != msgNotFound {
		t.Errorf("NewNotFound() message = %v, want %v", err.Message, msgNotFound)
	}
}

func TestIsNotFound(t *testing.T) {
	err := NewNotFound("test error")
	if !IsNotFound(err) {
		t.Error("IsNotFound() = false, want true")
	}
	if IsNotFound(errors.New("test error")) {
		t.Error("IsNotFound() = true, want false")
	}
}

func TestNewForbidden(t *testing.T) {
	err := NewForbidden("")
	if err.Code != Forbidden {
		t.Errorf("NewForbidden() code = %v, want %v", err.Code, Forbidden)
	}
	if err.Status != http.StatusForbidden {
		t.Errorf("NewForbidden() status = %v, want %v", err.Status, http.StatusForbidden)
	}
	if err.Message != msgForbidden {
		t.Errorf("NewForbidden() message = %v, want %v", err.Message, msgForbidden)
	}
}

func TestIsForbidden(t *testing.T) {
	err := NewForbidden("test error")
	if !IsForbidden(err) {
		t.Error("IsForbidden() = false, want true")
	}
	if IsForbidden(errors.New("test error")) {
		t.Error("IsForbidden() = true, want false")
	}
}

func TestNewInvalidArgument(t *testing.T) {
	err := NewInvalidArgument("")
	if err.Code != InvalidArgument {
		t.Errorf("NewInvalidArgument() code = %v, want %v", err.Code, InvalidArgument)
	}
	if err.Status != http.StatusBadRequest {
		t.Errorf("NewInvalidArgument() status = %v, want %v", err.Status, http.StatusBadRequest)
	}
	if err.Message != msgInvalidArgument {
		t.Errorf("NewInvalidArgument() message = %v, want %v", err.Message, msgInvalidArgument)
	}
}

func TestIsInvalidArgument(t *testing.T) {
	err := NewInvalidArgument("test error")
	if !IsInvalidArgument(err) {
		t.Error("IsInvalidArgument() = false, want true")
	}
	if IsInvalidArgument(errors.New("test error")) {
		t.Error("IsInvalidArgument() = true, want false")
	}
}

func TestNewFailedPrecondition(t *testing.T) {
	err := NewFailedPrecondition("")
	if err.Code != FailedPrecondition {
		t.Errorf("NewFailedPrecondition() code = %v, want %v", err.Code, FailedPrecondition)
	}
	if err.Status != http.StatusPreconditionFailed {
		t.Errorf("NewFailedPrecondition() status = %v, want %v", err.Status, http.StatusPreconditionFailed)
	}
	if err.Message != msgFailedPrecondition {
		t.Errorf("NewFailedPrecondition() message = %v, want %v", err.Message, msgFailedPrecondition)
	}
}

func TestIsFailedPrecondition(t *testing.T) {
	err := NewFailedPrecondition("test error")
	if !IsFailedPrecondition(err) {
		t.Error("IsFailedPrecondition() = false, want true")
	}
	if IsFailedPrecondition(errors.New("test error")) {
		t.Error("IsFailedPrecondition() = true, want false")
	}
}

func TestNewTooManyRequests(t *testing.T) {
	err := NewTooManyRequests("")
	assert.Equal(t, TooManyRequests, Code(err))
	assert.Equal(t, http.StatusTooManyRequests, Status(err))
	assert.Equal(t, msgTooManyRequests, err.Error())

	customMsg := "custom message"
	err = NewTooManyRequests(customMsg)
	assert.Equal(t, TooManyRequests, Code(err))
	assert.Equal(t, http.StatusTooManyRequests, Status(err))
	assert.Equal(t, customMsg, err.Error())
}

func TestIsTooManyRequests(t *testing.T) {
	err := NewTooManyRequests("")
	assert.True(t, IsTooManyRequests(err))
	assert.False(t, IsTooManyRequests(NewBadRequest("")))
}

func TestNewNotImplemented(t *testing.T) {
	err := NewNotImplemented("")
	assert.Equal(t, NotImplemented, Code(err))
	assert.Equal(t, http.StatusNotImplemented, Status(err))
	assert.Equal(t, msgNotImplemented, err.Error())

	customMsg := "custom message"
	err = NewNotImplemented(customMsg)
	assert.Equal(t, NotImplemented, Code(err))
	assert.Equal(t, http.StatusNotImplemented, Status(err))
	assert.Equal(t, customMsg, err.Error())
}

func TestIsNotImplemented(t *testing.T) {
	err := NewNotImplemented("")
	assert.True(t, IsNotImplemented(err))
	assert.False(t, IsNotImplemented(NewBadRequest("")))
}

func TestNewServiceUnavailable(t *testing.T) {
	err := NewServiceUnavailable("")
	assert.Equal(t, ServiceUnavailable, Code(err))
	assert.Equal(t, http.StatusServiceUnavailable, Status(err))
	assert.Equal(t, msgServiceUnavailable, err.Error())

	customMsg := "custom message"
	err = NewServiceUnavailable(customMsg)
	assert.Equal(t, ServiceUnavailable, Code(err))
	assert.Equal(t, http.StatusServiceUnavailable, Status(err))
	assert.Equal(t, customMsg, err.Error())
}

func TestIsServiceUnavailable(t *testing.T) {
	err := NewServiceUnavailable("")
	assert.True(t, IsServiceUnavailable(err))
	assert.False(t, IsServiceUnavailable(NewBadRequest("")))
}

func TestNewGatewayTimeout(t *testing.T) {
	err := NewGatewayTimeout("")
	assert.Equal(t, GatewayTimeout, Code(err))
	assert.Equal(t, http.StatusGatewayTimeout, Status(err))
	assert.Equal(t, msgGatewayTimeout, err.Error())

	customMsg := "custom message"
	err = NewGatewayTimeout(customMsg)
	assert.Equal(t, GatewayTimeout, Code(err))
	assert.Equal(t, http.StatusGatewayTimeout, Status(err))
	assert.Equal(t, customMsg, err.Error())
}

func TestIsGatewayTimeout(t *testing.T) {
	err := NewGatewayTimeout("")
	assert.True(t, IsGatewayTimeout(err))
	assert.False(t, IsGatewayTimeout(NewBadRequest("")))
}

func TestNewResourceExhausted(t *testing.T) {
	err := NewResourceExhausted("")
	assert.Equal(t, ResourceExhausted, Code(err))
	assert.Equal(t, http.StatusInsufficientStorage, Status(err))
	assert.Equal(t, msgResourceExhausted, err.Error())

	customMsg := "custom message"
	err = NewResourceExhausted(customMsg)
	assert.Equal(t, ResourceExhausted, Code(err))
	assert.Equal(t, http.StatusInsufficientStorage, Status(err))
	assert.Equal(t, customMsg, err.Error())
}

func TestIsResourceExhausted(t *testing.T) {
	err := NewResourceExhausted("")
	assert.True(t, IsResourceExhausted(err))
	assert.False(t, IsResourceExhausted(NewBadRequest("")))
}

func TestNewQuotaExceeded(t *testing.T) {
	err := NewQuotaExceeded("")
	assert.Equal(t, QuotaExceeded, Code(err))
	assert.Equal(t, http.StatusTooManyRequests, Status(err))
	assert.Equal(t, msgQuotaExceeded, err.Error())

	customMsg := "custom message"
	err = NewQuotaExceeded(customMsg)
	assert.Equal(t, QuotaExceeded, Code(err))
	assert.Equal(t, http.StatusTooManyRequests, Status(err))
	assert.Equal(t, customMsg, err.Error())
}

func TestIsQuotaExceeded(t *testing.T) {
	err := NewQuotaExceeded("")
	assert.True(t, IsQuotaExceeded(err))
	assert.False(t, IsQuotaExceeded(NewBadRequest("")))
}

func TestNewInvalidCredentials(t *testing.T) {
	err := NewInvalidCredentials("")
	assert.Equal(t, InvalidCredentials, Code(err))
	assert.Equal(t, http.StatusUnauthorized, Status(err))
	assert.Equal(t, msgInvalidCredentials, err.Error())

	customMsg := "custom message"
	err = NewInvalidCredentials(customMsg)
	assert.Equal(t, InvalidCredentials, Code(err))
	assert.Equal(t, http.StatusUnauthorized, Status(err))
	assert.Equal(t, customMsg, err.Error())
}

func TestIsInvalidCredentials(t *testing.T) {
	err := NewInvalidCredentials("")
	assert.True(t, IsInvalidCredentials(err))
	assert.False(t, IsInvalidCredentials(NewBadRequest("")))
}

func TestNewTokenExpired(t *testing.T) {
	err := NewTokenExpired("")
	assert.Equal(t, TokenExpired, Code(err))
	assert.Equal(t, http.StatusUnauthorized, Status(err))
	assert.Equal(t, msgTokenExpired, err.Error())

	customMsg := "custom message"
	err = NewTokenExpired(customMsg)
	assert.Equal(t, TokenExpired, Code(err))
	assert.Equal(t, http.StatusUnauthorized, Status(err))
	assert.Equal(t, customMsg, err.Error())
}

func TestIsTokenExpired(t *testing.T) {
	err := NewTokenExpired("")
	assert.True(t, IsTokenExpired(err))
	assert.False(t, IsTokenExpired(NewBadRequest("")))
}

func TestNewConnectionFailed(t *testing.T) {
	err := NewConnectionFailed("")
	assert.Equal(t, ConnectionFailed, Code(err))
	assert.Equal(t, http.StatusServiceUnavailable, Status(err))
	assert.Equal(t, msgConnectionFailed, err.Error())

	customMsg := "custom message"
	err = NewConnectionFailed(customMsg)
	assert.Equal(t, ConnectionFailed, Code(err))
	assert.Equal(t, http.StatusServiceUnavailable, Status(err))
	assert.Equal(t, customMsg, err.Error())
}

func TestIsConnectionFailed(t *testing.T) {
	err := NewConnectionFailed("")
	assert.True(t, IsConnectionFailed(err))
	assert.False(t, IsConnectionFailed(NewBadRequest("")))
}

func TestNewConnectionTimeout(t *testing.T) {
	err := NewConnectionTimeout("")
	assert.Equal(t, ConnectionTimeout, Code(err))
	assert.Equal(t, http.StatusGatewayTimeout, Status(err))
	assert.Equal(t, msgConnectionTimeout, err.Error())

	customMsg := "custom message"
	err = NewConnectionTimeout(customMsg)
	assert.Equal(t, ConnectionTimeout, Code(err))
	assert.Equal(t, http.StatusGatewayTimeout, Status(err))
	assert.Equal(t, customMsg, err.Error())
}

func TestIsConnectionTimeout(t *testing.T) {
	err := NewConnectionTimeout("")
	assert.True(t, IsConnectionTimeout(err))
	assert.False(t, IsConnectionTimeout(NewBadRequest("")))
}

func TestNewConnectionRefused(t *testing.T) {
	err := NewConnectionRefused("")
	assert.Equal(t, ConnectionRefused, Code(err))
	assert.Equal(t, http.StatusServiceUnavailable, Status(err))
	assert.Equal(t, msgConnectionRefused, err.Error())

	customMsg := "custom message"
	err = NewConnectionRefused(customMsg)
	assert.Equal(t, ConnectionRefused, Code(err))
	assert.Equal(t, http.StatusServiceUnavailable, Status(err))
	assert.Equal(t, customMsg, err.Error())
}

func TestIsConnectionRefused(t *testing.T) {
	err := NewConnectionRefused("")
	assert.True(t, IsConnectionRefused(err))
	assert.False(t, IsConnectionRefused(NewBadRequest("")))
}

func TestNewDeadlineExceeded(t *testing.T) {
	err := NewDeadlineExceeded("")
	assert.Equal(t, DeadlineExceeded, Code(err))
	assert.Equal(t, http.StatusGatewayTimeout, Status(err))
	assert.Equal(t, msgDeadlineExceeded, err.Error())

	customMsg := "custom message"
	err = NewDeadlineExceeded(customMsg)
	assert.Equal(t, DeadlineExceeded, Code(err))
	assert.Equal(t, http.StatusGatewayTimeout, Status(err))
	assert.Equal(t, customMsg, err.Error())
}

func TestIsDeadlineExceeded(t *testing.T) {
	err := NewDeadlineExceeded("")
	assert.True(t, IsDeadlineExceeded(err))
	assert.False(t, IsDeadlineExceeded(NewBadRequest("")))
}

func TestNewRequestTimeout(t *testing.T) {
	err := NewRequestTimeout("")
	assert.Equal(t, RequestTimeout, Code(err))
	assert.Equal(t, http.StatusRequestTimeout, Status(err))
	assert.Equal(t, msgRequestTimeout, err.Error())

	customMsg := "custom message"
	err = NewRequestTimeout(customMsg)
	assert.Equal(t, RequestTimeout, Code(err))
	assert.Equal(t, http.StatusRequestTimeout, Status(err))
	assert.Equal(t, customMsg, err.Error())
}

func TestIsRequestTimeout(t *testing.T) {
	err := NewRequestTimeout("")
	assert.True(t, IsRequestTimeout(err))
	assert.False(t, IsRequestTimeout(NewBadRequest("")))
}

func TestNewOutOfRange(t *testing.T) {
	err := NewOutOfRange("")
	assert.Equal(t, OutOfRange, Code(err))
	assert.Equal(t, http.StatusBadRequest, Status(err))
	assert.Equal(t, msgOutOfRange, err.Error())

	customMsg := "custom message"
	err = NewOutOfRange(customMsg)
	assert.Equal(t, OutOfRange, Code(err))
	assert.Equal(t, http.StatusBadRequest, Status(err))
	assert.Equal(t, customMsg, err.Error())
}

func TestIsOutOfRange(t *testing.T) {
	err := NewOutOfRange("")
	assert.True(t, IsOutOfRange(err))
	assert.False(t, IsOutOfRange(NewBadRequest("")))
}

func TestWithRequestID(t *testing.T) {
	err := New(BadRequest, "test error").WithRequestID("req-123")
	assert.Equal(t, "req-123", err.RequestID)
}

func TestWithDetails(t *testing.T) {
	err := New(BadRequest, "test error").WithDetails("foo", 42)
	assert.Equal(t, 42, err.Details["foo"])
}

func TestFromHTTPStatus(t *testing.T) {
	assert.Equal(t, BadRequest, FromHTTPStatus(http.StatusBadRequest, "bad").Code)
	assert.Equal(t, Unauthorized, FromHTTPStatus(http.StatusUnauthorized, "unauth").Code)
	assert.Equal(t, Forbidden, FromHTTPStatus(http.StatusForbidden, "forbidden").Code)
	assert.Equal(t, NotFound, FromHTTPStatus(http.StatusNotFound, "notfound").Code)
	assert.Equal(t, Conflict, FromHTTPStatus(http.StatusConflict, "conflict").Code)
	assert.Equal(t, Gone, FromHTTPStatus(http.StatusGone, "gone").Code)
	assert.Equal(t, UnprocessableEntity, FromHTTPStatus(http.StatusUnprocessableEntity, "unprocessable").Code)
	assert.Equal(t, TooManyRequests, FromHTTPStatus(http.StatusTooManyRequests, "toomany").Code)
	assert.Equal(t, UnexpectedFailure, FromHTTPStatus(http.StatusInternalServerError, "fail").Code)
	assert.Equal(t, NotImplemented, FromHTTPStatus(http.StatusNotImplemented, "notimpl").Code)
	assert.Equal(t, ServiceUnavailable, FromHTTPStatus(http.StatusServiceUnavailable, "unavail").Code)
	assert.Equal(t, GatewayTimeout, FromHTTPStatus(http.StatusGatewayTimeout, "timeout").Code)
	// Default fallback
	assert.Equal(t, UnexpectedFailure, FromHTTPStatus(599, "other").Code)
}

func TestFromContextError(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 0)
	defer cancel()
	<-ctx.Done()
	cerr := ctx.Err()
	if cerr == context.DeadlineExceeded {
		assert.Equal(t, DeadlineExceeded, FromContextError(cerr).Code)
	}
	assert.Equal(t, RequestTimeout, FromContextError(context.Canceled).Code)
	assert.Contains(t, FromContextError(errors.New("foo")).Error(), "Context error")
}

func TestFromNetworkError(t *testing.T) {
	err := errors.New("neterr")
	converted := FromNetworkError(err)
	assert.Equal(t, ConnectionFailed, converted.Code)
	assert.Equal(t, err, converted.Err)
	assert.Nil(t, FromNetworkError(nil))
}

func TestIsRetryable(t *testing.T) {
	assert.True(t, IsRetryable(New(ServiceUnavailable, "")))
	assert.True(t, IsRetryable(New(GatewayTimeout, "")))
	assert.True(t, IsRetryable(New(ConnectionFailed, "")))
	assert.True(t, IsRetryable(New(ConnectionTimeout, "")))
	assert.True(t, IsRetryable(New(ConnectionRefused, "")))
	assert.True(t, IsRetryable(New(RequestTimeout, "")))
	assert.False(t, IsRetryable(New(BadRequest, "")))
}

func TestIsTimeout(t *testing.T) {
	assert.True(t, IsTimeout(New(DeadlineExceeded, "")))
	assert.True(t, IsTimeout(New(RequestTimeout, "")))
	assert.True(t, IsTimeout(New(GatewayTimeout, "")))
	assert.True(t, IsTimeout(New(ConnectionTimeout, "")))
	assert.False(t, IsTimeout(New(BadRequest, "")))
}

func TestIsAuthError(t *testing.T) {
	assert.True(t, IsAuthError(New(NoAuthorization, "")))
	assert.True(t, IsAuthError(New(InvalidCredentials, "")))
	assert.True(t, IsAuthError(New(TokenExpired, "")))
	assert.False(t, IsAuthError(New(BadRequest, "")))
}

func TestIsClientError(t *testing.T) {
	assert.True(t, IsClientError(New(BadRequest, "")))
	assert.False(t, IsClientError(New(ServiceUnavailable, "")))
}

func TestIsServerError(t *testing.T) {
	assert.True(t, IsServerError(New(ServiceUnavailable, "")))
	assert.False(t, IsServerError(New(BadRequest, "")))
}

func TestCodeAndStatusFallbacks(t *testing.T) {
	// nil error
	assert.Equal(t, ErrorCode(""), Code(nil))
	assert.Equal(t, http.StatusOK, Status(nil))

	// Standard error (kein *Error)
	stdErr := errors.New("foo")
	assert.Equal(t, UnexpectedFailure, Code(stdErr))
	assert.Equal(t, http.StatusInternalServerError, Status(stdErr))
}

func TestWithDetailsNilMap(t *testing.T) {
	err := New(BadRequest, "test error")
	err.Details = nil // explizit nil setzen
	err = err.WithDetails("bar", 99)
	assert.Equal(t, 99, err.Details["bar"])
}

func TestIsAuthErrorNegative(t *testing.T) {
	assert.False(t, IsAuthError(errors.New("foo")))
}

func TestIsClientErrorNegative(t *testing.T) {
	assert.False(t, IsClientError(errors.New("foo")))
}

func TestIsServerErrorNegative(t *testing.T) {
	assert.False(t, IsServerError(errors.New("foo")))
}
