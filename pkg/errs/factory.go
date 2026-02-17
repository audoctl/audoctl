package errs

import (
	"fmt"
	"net/http"
	"time"
)

// Common error factory functions for typical scenarios

// NotFound creates a not found error
func NotFound(resource, id string) *Err {
	return New(
		fmt.Sprintf("%s.not_found", resource),
		fmt.Sprintf("%s not found", resource),
		http.StatusNotFound,
		CategoryNotFound,
	).SetParam("resource", resource).SetParam("id", id)
}

// Validation creates a validation error
func Validation(message string, validations ...ValidationError) *Err {
	err := New(
		"validation.failed",
		message,
		http.StatusBadRequest,
		CategoryValidation,
	)

	for _, v := range validations {
		err.AddValidation(v)
	}

	return err
}

// Unauthorized creates an unauthorized error
func Unauthorized(message string) *Err {
	if message == "" {
		message = "Unauthorized access"
	}
	return New(
		"auth.unauthorized",
		message,
		http.StatusUnauthorized,
		CategoryUnauthorized,
	)
}

// Forbidden creates a forbidden error
func Forbidden(message string) *Err {
	if message == "" {
		message = "Access forbidden"
	}
	return New(
		"auth.forbidden",
		message,
		http.StatusForbidden,
		CategoryForbidden,
	)
}

// Conflict creates a conflict error
func Conflict(resource, reason string) *Err {
	return New(
		fmt.Sprintf("%s.conflict", resource),
		fmt.Sprintf("Conflict: %s", reason),
		http.StatusConflict,
		CategoryConflict,
	).SetParam("resource", resource).SetParam("reason", reason)
}

// Internal creates an internal server error
func Internal(message string, cause error) *Err {
	if message == "" {
		message = "Internal server error"
	}
	return New(
		"internal.error",
		message,
		http.StatusInternalServerError,
		CategoryInternal,
	).WithOptions(WithCause(cause))
}

// External creates an external service error
func External(service string, cause error) *Err {
	return New(
		"external.service_error",
		fmt.Sprintf("External service error: %s", service),
		http.StatusBadGateway,
		CategoryExternal,
	).SetParam("service", service).WithOptions(
		WithCause(cause),
		WithRetryable(true),
	)
}

// Timeout creates a timeout error
func Timeout(operation string, duration time.Duration) *Err {
	return New(
		"timeout",
		fmt.Sprintf("Operation timed out: %s", operation),
		http.StatusRequestTimeout,
		CategoryTimeout,
	).SetParams(map[string]any{
		"operation": operation,
		"duration":  duration.Seconds(),
	}).WithOptions(WithRetryable(true))
}

// RateLimit creates a rate limit error
func RateLimit(retryAfter time.Duration) *Err {
	return New(
		"rate_limit.exceeded",
		"Rate limit exceeded",
		http.StatusTooManyRequests,
		CategoryRateLimit,
	).WithOptions(WithRetryAfter(retryAfter))
}

// BadRequest creates a bad request error
func BadRequest(message string) *Err {
	if message == "" {
		message = "Bad request"
	}
	return New(
		"request.bad_request",
		message,
		http.StatusBadRequest,
		CategoryValidation,
	)
}

// SessionNotFound creates a session not found error (specific to audoctl)
func SessionNotFound(sessionID string) *Err {
	return NotFound("session", sessionID).
		SetMetadata("domain", "audoctl").
		WithOptions(WithSessionID(sessionID))
}

// EventNotFound creates an event not found error (specific to audoctl)
func EventNotFound(eventID string) *Err {
	return NotFound("event", eventID).
		SetMetadata("domain", "audoctl")
}

// InvalidEventType creates an invalid event type error (specific to audoctl)
func InvalidEventType(eventType string, validTypes []string) *Err {
	return Validation(
		"Invalid event type",
		NewValidationError(
			"event.invalid_type",
			"type",
			fmt.Sprintf("Event type '%s' is not valid", eventType),
			map[string]any{
				"provided": eventType,
				"valid":    validTypes,
			},
		),
	)
}

// SessionAlreadyFinished creates a session already finished error (specific to audoctl)
func SessionAlreadyFinished(sessionID string) *Err {
	return Conflict(
		"session",
		"Session has already been finished",
	).SetParams(map[string]any{
		"session_id": sessionID,
		"state":      "finished",
	}).WithOptions(WithSessionID(sessionID))
}

// DatabaseError creates a database error
func DatabaseError(operation string, cause error) *Err {
	return Internal(
		fmt.Sprintf("Database error during %s", operation),
		cause,
	).SetMetadata("operation", operation).SetMetadata("layer", "database")
}

// StorageError creates a storage error (specific to audoctl)
func StorageError(operation string, cause error) *Err {
	return Internal(
		fmt.Sprintf("Storage error during %s", operation),
		cause,
	).SetParams(map[string]any{
		"operation": operation,
		"layer":     "storage",
	}).WithOptions(WithCause(cause), WithRetryable(true))
}
