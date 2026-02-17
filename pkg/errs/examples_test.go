package errs_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/audoctl/audoctl/pkg/errs"
)

// Example_basicError demonstrates basic error creation
func Example_basicError() {
	err := errs.NotFound("session", "sess_123")
	fmt.Println(err.Error())
	// Output: session not found
}

// Example_validationError demonstrates validation errors
func Example_validationError() {
	err := errs.Validation(
		"Request validation failed",
		errs.NewValidationError(
			"agent.required",
			"agent",
			"Agent name is required",
			nil,
		),
		errs.NewValidationError(
			"metadata.invalid",
			"metadata.timeout",
			"Timeout must be positive",
			map[string]any{"provided": -5},
		),
	)

	data, _ := json.MarshalIndent(err, "", "  ")
	fmt.Println(string(data))
}

// Example_wrappedError demonstrates error wrapping
func TestWrappedError(t *testing.T) {
	// Simulate a database error
	dbErr := errors.New("connection timeout")

	// Wrap it with context
	err := errs.DatabaseError("fetch session", dbErr)

	// Add request context
	err.WithOptions(
		errs.WithRequestID("req_abc123"),
		errs.WithTraceID("trace_xyz789"),
	)

	if err.Error() != "Database error during fetch session: connection timeout" {
		t.Errorf("Expected error message 'database error during fetch session: connection timeout', got '%s'", err.Error())
	}

	if errors.Unwrap(err) != dbErr {
		t.Errorf("Expected Unwrap to return root error, got '%v'", errors.Unwrap(err))
	}
}

// Example_sessionError demonstrates audoctl-specific errors
func TestSessionError(t *testing.T) {
	sessionID := "sess_abc123"

	err := errs.SessionNotFound(sessionID).
		WithOptions(
			errs.WithRequestID("req_xyz"),
			errs.WithHelpURL("https://docs.audoctl.dev/errors/session-not-found"),
		)

	data, _ := json.MarshalIndent(err, "", "  ")
	if string(data) != `{"id":"err_1708123456789","key":"session.not_found","category":"not_found","message":"Session not found","status":404,"timestamp":"2026-02-16T00:00:00Z","request_id":"req_xyz","help_url":"https://docs.audoctl.dev/errors/session-not-found"}` {
		t.Errorf("Expected error message '{\"id\":\"err_1708123456789\",\"key\":\"session.not_found\",\"category\":\"not_found\",\"message\":\"Session not found\",\"status\":404,\"timestamp\":\"2026-02-16T00:00:00Z\",\"request_id\":\"req_xyz\",\"help_url\":\"https://docs.audoctl.dev/errors/session-not-found\"}', got '%s'", string(data))
	}
}

// Example_retryableError demonstrates retryable errors
func TestRetryableError(t *testing.T) {
	err := errs.RateLimit(30 * time.Second)

	fmt.Println("Error:", err.Error())
	fmt.Println("Retryable:", err.Retryable)
	if err.RetryAfter != nil {
		t.Errorf("Expected RetryAfter to be set, got '%v'", err.RetryAfter)
	}
}

// Example_errorChaining demonstrates error chain handling
func TestErrorChaining(t *testing.T) {
	// Root cause
	rootErr := errors.New("network unreachable")

	// Wrap with context
	storageErr := errs.StorageError("persist event", rootErr)

	// Wrap again with business context
	finalErr := errs.Wrap(
		storageErr,
		"event.persist_failed",
		"Failed to persist event to timeline",
		500,
		errs.CategoryInternal,
	).WithOptions(
		errs.WithSessionID("sess_123"),
		errs.WithRequestID("req_456"),
	)

	// Unwrap to get original
	if finalErr.Error() != "Failed to persist event to timeline: Storage error during persist event: network unreachable" {
		t.Errorf("Expected error message 'Failed to persist event to timeline: network unreachable', got '%s'", finalErr.Error())
	}
	if finalErr.Cause() != rootErr {
		t.Errorf("Expected Cause to return root error, got '%v'", finalErr.Cause())
	}

	// Check error type
	if storageErr, ok := errs.ToErr(errors.Unwrap(finalErr)); ok {
		fmt.Println("Storage layer error detected:", storageErr.Key)
	}
}

// Example_withDebugInfo demonstrates debug information
func TestWithDebugInfo(t *testing.T) {
	err := errs.Internal("Something went wrong", errors.New("null pointer"))

	// In production: includeDebug = false
	// In development: includeDebug = true
	err.WithDebugInfo(true)

	if err.Debug != nil {
		fmt.Println("Has debug info: true")
		fmt.Println("Has stack trace:", len(err.Debug.StackTrace) > 0)
		fmt.Println("Cause:", err.Debug.Cause)
	}
	if err.Debug.StackTrace == nil {
		t.Errorf("Expected stack trace to be present, got '%v'", err.Debug.StackTrace)
	}
	if err.Debug.Cause == "" {
		t.Errorf("Expected cause to be present, got '%s'", err.Debug.Cause)
	}
}

// Example_customError demonstrates building custom errors
func TestCustomError(t *testing.T) {
	err := errs.New(
		"agent.execution_failed",
		"AI agent execution failed",
		500,
		errs.CategoryInternal,
	).SetParams(map[string]any{
		"agent": "refund_agent",
		"step":  "tool_execution",
		"tool":  "process_refund",
	}).SetMetadata("cost_usd", 0.05).
		SetMetadata("tokens_used", 1234).
		WithOptions(
			errs.WithSessionID("sess_789"),
			errs.WithRetryable(true),
			errs.WithHelpURL("https://docs.audoctl.dev/troubleshooting/agent-execution"),
		)

	data, _ := json.MarshalIndent(err, "", "  ")
	fmt.Println(string(data))
}

// Example_errorComparison demonstrates error comparison
func TestErrorComparison(t *testing.T) {
	err1 := errs.NotFound("session", "123")
	err2 := errs.NotFound("session", "456")
	err3 := errs.Unauthorized("Invalid token")

	// Using errors.Is
	fmt.Println("err1 is err2:", errors.Is(err1, err2)) // true (same key and status)
	if !errors.Is(err1, err2) {
		t.Errorf("Expected err1 to be equal to err2, got '%v'", errors.Is(err1, err2))
	}
	if errors.Is(err1, err3) {
		t.Errorf("Expected err1 to be not equal to err3, got '%v'", errors.Is(err1, err3))
	}
}

// Example_httpMiddleware demonstrates how to use in HTTP handler
func TestHttpMiddleware(t *testing.T) {
	// Simulated HTTP handler
	handleRequest := func(sessionID string) error {
		// Business logic
		if sessionID == "" {
			return errs.BadRequest("Session ID is required").
				AddValidation(errs.NewValidationError(
					"session_id.required",
					"session_id",
					"Session ID cannot be empty",
					nil,
				))
		}

		// Check if session exists
		sessionExists := false // Simulated check
		if !sessionExists {
			return errs.SessionNotFound(sessionID)
		}

		return nil
	}

	// Middleware/handler error response
	err := handleRequest("")
	if err != nil {
		if e, ok := errs.ToErr(err); ok {
			// Set response status
			statusCode := e.Status

			// Include debug info based on environment
			isDevelopment := true
			e.WithDebugInfo(isDevelopment)

			// Serialize to JSON
			data, _ := json.MarshalIndent(e, "", "  ")

			fmt.Printf("HTTP %d\n", statusCode)
			fmt.Println(string(data))
		}
	}
}
