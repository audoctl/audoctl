package errs

import (
	"errors"
	"net/http"
	"testing"
	"time"
)

func TestNotFound(t *testing.T) {
	err := NotFound("session", "sess_123")

	if err.Status != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", err.Status)
	}

	if err.Category != CategoryNotFound {
		t.Errorf("Expected category 'not_found', got '%s'", err.Category)
	}

	if err.Params["resource"] != "session" {
		t.Errorf("Expected resource 'session', got '%v'", err.Params["resource"])
	}

	if err.Params["id"] != "sess_123" {
		t.Errorf("Expected id 'sess_123', got '%v'", err.Params["id"])
	}
}

func TestValidation(t *testing.T) {
	v1 := NewValidationError("field.required", "name", "Name is required", nil)
	v2 := NewValidationError("field.invalid", "email", "Email is invalid", nil)

	err := Validation("Validation failed", v1, v2)

	if err.Status != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", err.Status)
	}

	if err.Category != CategoryValidation {
		t.Errorf("Expected category 'validation', got '%s'", err.Category)
	}

	if len(err.Validations) != 2 {
		t.Errorf("Expected 2 validations, got %d", len(err.Validations))
	}
}

func TestUnauthorized(t *testing.T) {
	err := Unauthorized("Invalid token")

	if err.Status != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", err.Status)
	}

	if err.Category != CategoryUnauthorized {
		t.Errorf("Expected category 'unauthorized', got '%s'", err.Category)
	}

	if err.Message != "Invalid token" {
		t.Errorf("Expected message 'Invalid token', got '%s'", err.Message)
	}

	// Test with empty message
	err2 := Unauthorized("")
	if err2.Message != "Unauthorized access" {
		t.Errorf("Expected default message, got '%s'", err2.Message)
	}
}

func TestForbidden(t *testing.T) {
	err := Forbidden("Insufficient permissions")

	if err.Status != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", err.Status)
	}

	if err.Category != CategoryForbidden {
		t.Errorf("Expected category 'forbidden', got '%s'", err.Category)
	}
}

func TestConflict(t *testing.T) {
	err := Conflict("session", "already exists")

	if err.Status != http.StatusConflict {
		t.Errorf("Expected status 409, got %d", err.Status)
	}

	if err.Category != CategoryConflict {
		t.Errorf("Expected category 'conflict', got '%s'", err.Category)
	}

	if err.Params["resource"] != "session" {
		t.Errorf("Expected resource 'session', got '%v'", err.Params["resource"])
	}
}

func TestInternal(t *testing.T) {
	cause := errors.New("database connection failed")
	err := Internal("Operation failed", cause)

	if err.Status != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", err.Status)
	}

	if err.Category != CategoryInternal {
		t.Errorf("Expected category 'internal', got '%s'", err.Category)
	}

	if errors.Unwrap(err) != cause {
		t.Error("Expected cause to be wrapped")
	}
}

func TestExternal(t *testing.T) {
	cause := errors.New("connection timeout")
	err := External("OpenAI", cause)

	if err.Status != http.StatusBadGateway {
		t.Errorf("Expected status 502, got %d", err.Status)
	}

	if err.Category != CategoryExternal {
		t.Errorf("Expected category 'external', got '%s'", err.Category)
	}

	if !err.Retryable {
		t.Error("Expected external error to be retryable")
	}

	if err.Params["service"] != "OpenAI" {
		t.Errorf("Expected service 'OpenAI', got '%v'", err.Params["service"])
	}
}

func TestTimeout(t *testing.T) {
	duration := 30 * time.Second
	err := Timeout("llm_call", duration)

	if err.Status != http.StatusRequestTimeout {
		t.Errorf("Expected status 408, got %d", err.Status)
	}

	if err.Category != CategoryTimeout {
		t.Errorf("Expected category 'timeout', got '%s'", err.Category)
	}

	if !err.Retryable {
		t.Error("Expected timeout error to be retryable")
	}

	if err.Params["operation"] != "llm_call" {
		t.Errorf("Expected operation 'llm_call', got '%v'", err.Params["operation"])
	}
}

func TestRateLimit(t *testing.T) {
	retryAfter := 60 * time.Second
	err := RateLimit(retryAfter)

	if err.Status != http.StatusTooManyRequests {
		t.Errorf("Expected status 429, got %d", err.Status)
	}

	if err.Category != CategoryRateLimit {
		t.Errorf("Expected category 'rate_limit', got '%s'", err.Category)
	}

	if err.RetryAfter == nil {
		t.Fatal("Expected RetryAfter to be set")
	}

	if *err.RetryAfter != retryAfter {
		t.Errorf("Expected RetryAfter %v, got %v", retryAfter, *err.RetryAfter)
	}
}

func TestBadRequest(t *testing.T) {
	err := BadRequest("Invalid JSON")

	if err.Status != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", err.Status)
	}

	if err.Category != CategoryValidation {
		t.Errorf("Expected category 'validation', got '%s'", err.Category)
	}
}

func TestSessionNotFound(t *testing.T) {
	sessionID := "sess_123"
	err := SessionNotFound(sessionID)

	if err.Status != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", err.Status)
	}

	if err.SessionID != sessionID {
		t.Errorf("Expected SessionID '%s', got '%s'", sessionID, err.SessionID)
	}

	if err.Metadata["domain"] != "audoctl" {
		t.Errorf("Expected domain 'audoctl', got '%v'", err.Metadata["domain"])
	}
}

func TestEventNotFound(t *testing.T) {
	eventID := "evt_456"
	err := EventNotFound(eventID)

	if err.Status != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", err.Status)
	}

	if err.Params["id"] != eventID {
		t.Errorf("Expected id '%s', got '%v'", eventID, err.Params["id"])
	}
}

func TestInvalidEventType(t *testing.T) {
	validTypes := []string{"prompt", "llm_call", "tool_execution"}
	err := InvalidEventType("invalid_type", validTypes)

	if err.Status != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", err.Status)
	}

	if err.Category != CategoryValidation {
		t.Errorf("Expected category 'validation', got '%s'", err.Category)
	}

	if len(err.Validations) != 1 {
		t.Errorf("Expected 1 validation error, got %d", len(err.Validations))
	}

	validation := err.Validations[0]
	if validation.Field != "type" {
		t.Errorf("Expected field 'type', got '%s'", validation.Field)
	}
}

func TestSessionAlreadyFinished(t *testing.T) {
	sessionID := "sess_789"
	err := SessionAlreadyFinished(sessionID)

	if err.Status != http.StatusConflict {
		t.Errorf("Expected status 409, got %d", err.Status)
	}

	if err.Category != CategoryConflict {
		t.Errorf("Expected category 'conflict', got '%s'", err.Category)
	}

	if err.SessionID != sessionID {
		t.Errorf("Expected SessionID '%s', got '%s'", sessionID, err.SessionID)
	}

	if err.Params["state"] != "finished" {
		t.Errorf("Expected state 'finished', got '%v'", err.Params["state"])
	}
}

func TestDatabaseError(t *testing.T) {
	cause := errors.New("connection lost")
	err := DatabaseError("insert", cause)

	if err.Status != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", err.Status)
	}

	if err.Metadata["operation"] != "insert" {
		t.Errorf("Expected operation 'insert', got '%v'", err.Metadata["operation"])
	}

	if err.Metadata["layer"] != "database" {
		t.Errorf("Expected layer 'database', got '%v'", err.Metadata["layer"])
	}

	if errors.Unwrap(err) != cause {
		t.Error("Expected cause to be wrapped")
	}
}

func TestStorageError(t *testing.T) {
	cause := errors.New("disk full")
	err := StorageError("write", cause)

	if err.Status != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", err.Status)
	}

	if !err.Retryable {
		t.Error("Expected storage error to be retryable")
	}

	if err.Params["operation"] != "write" {
		t.Errorf("Expected operation 'write', got '%v'", err.Params["operation"])
	}

	if errors.Unwrap(err) != cause {
		t.Error("Expected cause to be wrapped")
	}
}
