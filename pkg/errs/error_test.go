package errs

import (
	"encoding/json"
	"errors"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	err := New(
		"test.error",
		"Test error message",
		400,
		CategoryValidation,
	)

	if err.Key != "test.error" {
		t.Errorf("Expected key 'test.error', got '%s'", err.Key)
	}

	if err.Message != "Test error message" {
		t.Errorf("Expected message 'Test error message', got '%s'", err.Message)
	}

	if err.Status != 400 {
		t.Errorf("Expected status 400, got %d", err.Status)
	}

	if err.Category != CategoryValidation {
		t.Errorf("Expected category 'validation', got '%s'", err.Category)
	}

	if err.ID == "" {
		t.Error("Expected non-empty error ID")
	}

	if err.Timestamp.IsZero() {
		t.Error("Expected non-zero timestamp")
	}
}

func TestWrap(t *testing.T) {
	rootErr := errors.New("root cause")
	err := Wrap(rootErr, "wrapped.error", "Wrapped error", 500, CategoryInternal)

	if err.Cause() != rootErr {
		t.Error("Expected root cause to be preserved")
	}

	if err.Error() != "Wrapped error: root cause" {
		t.Errorf("Expected 'Wrapped error: root cause', got '%s'", err.Error())
	}

	unwrapped := errors.Unwrap(err)
	if unwrapped != rootErr {
		t.Error("Expected Unwrap to return root error")
	}
}

func TestAddValidation(t *testing.T) {
	err := New("validation.failed", "Validation failed", 400, CategoryValidation)

	validation := NewValidationError(
		"field.required",
		"name",
		"Name is required",
		nil,
	)

	err.AddValidation(validation)

	if len(err.Validations) != 1 {
		t.Errorf("Expected 1 validation error, got %d", len(err.Validations))
	}

	if err.Validations[0].Field != "name" {
		t.Errorf("Expected field 'name', got '%s'", err.Validations[0].Field)
	}
}

func TestSetParams(t *testing.T) {
	err := New("test.error", "Test", 400, CategoryValidation)

	err.SetParams(map[string]any{
		"key1": "value1",
		"key2": 123,
	})

	if err.Params["key1"] != "value1" {
		t.Errorf("Expected param 'key1' to be 'value1', got '%v'", err.Params["key1"])
	}

	if err.Params["key2"] != 123 {
		t.Errorf("Expected param 'key2' to be 123, got '%v'", err.Params["key2"])
	}
}

func TestSetMetadata(t *testing.T) {
	err := New("test.error", "Test", 400, CategoryValidation)

	err.SetMetadata("trace", "abc123")
	err.SetMetadata("duration_ms", 150)

	if err.Metadata["trace"] != "abc123" {
		t.Errorf("Expected metadata 'trace' to be 'abc123', got '%v'", err.Metadata["trace"])
	}

	if err.Metadata["duration_ms"] != 150 {
		t.Errorf("Expected metadata 'duration_ms' to be 150, got '%v'", err.Metadata["duration_ms"])
	}
}

func TestWithOptions(t *testing.T) {
	err := New("test.error", "Test", 400, CategoryValidation).
		WithOptions(
			WithRequestID("req_123"),
			WithTraceID("trace_456"),
			WithSessionID("sess_789"),
			WithRetryable(true),
		)

	if err.RequestID != "req_123" {
		t.Errorf("Expected RequestID 'req_123', got '%s'", err.RequestID)
	}

	if err.TraceID != "trace_456" {
		t.Errorf("Expected TraceID 'trace_456', got '%s'", err.TraceID)
	}

	if err.SessionID != "sess_789" {
		t.Errorf("Expected SessionID 'sess_789', got '%s'", err.SessionID)
	}

	if !err.Retryable {
		t.Error("Expected error to be retryable")
	}
}

func TestWithRetryAfter(t *testing.T) {
	duration := 30 * time.Second
	err := New("test.error", "Test", 429, CategoryRateLimit).
		WithOptions(WithRetryAfter(duration))

	if !err.Retryable {
		t.Error("Expected error to be retryable when RetryAfter is set")
	}

	if err.RetryAfter == nil {
		t.Fatal("Expected RetryAfter to be set")
	}

	if *err.RetryAfter != duration {
		t.Errorf("Expected RetryAfter to be %v, got %v", duration, *err.RetryAfter)
	}
}

func TestWithDebugInfo(t *testing.T) {
	rootErr := errors.New("root cause")
	err := Wrap(rootErr, "test.error", "Test error", 500, CategoryInternal)

	err.WithDebugInfo(true)

	if err.Debug == nil {
		t.Fatal("Expected debug info to be present")
	}

	if err.Debug.Cause != "root cause" {
		t.Errorf("Expected cause 'root cause', got '%s'", err.Debug.Cause)
	}

	if len(err.Debug.StackTrace) == 0 {
		t.Error("Expected stack trace to be present")
	}

	if len(err.Debug.ErrorChain) == 0 {
		t.Error("Expected error chain to be present")
	}
}

func TestToErr(t *testing.T) {
	originalErr := New("test.error", "Test", 400, CategoryValidation)

	// Test with Err
	convertedErr, ok := ToErr(originalErr)
	if !ok {
		t.Error("Expected conversion to succeed for *Err")
	}
	if convertedErr.Key != "test.error" {
		t.Errorf("Expected key 'test.error', got '%s'", convertedErr.Key)
	}

	// Test with standard error
	stdErr := errors.New("standard error")
	_, ok = ToErr(stdErr)
	if ok {
		t.Error("Expected conversion to fail for standard error")
	}

	// Test with nil
	_, ok = ToErr(nil)
	if ok {
		t.Error("Expected conversion to fail for nil error")
	}
}

func TestErrorIs(t *testing.T) {
	err1 := New("test.error", "Test", 400, CategoryValidation)
	err2 := New("test.error", "Test", 400, CategoryValidation)
	err3 := New("other.error", "Other", 500, CategoryInternal)

	if !errors.Is(err1, err2) {
		t.Error("Expected errors with same key and status to be equal")
	}

	if errors.Is(err1, err3) {
		t.Error("Expected errors with different key/status to not be equal")
	}
}

func TestErrorAs(t *testing.T) {
	err := New("test.error", "Test", 400, CategoryValidation)

	var target *Err
	if !errors.As(err, &target) {
		t.Error("Expected errors.As to succeed")
	}

	if target.Key != "test.error" {
		t.Errorf("Expected key 'test.error', got '%s'", target.Key)
	}
}

func TestMarshalJSON(t *testing.T) {
	err := New("test.error", "Test error", 400, CategoryValidation).
		WithOptions(
			WithRequestID("req_123"),
			WithRetryable(true),
		)

	err.SetParam("key", "value")

	data, jsonErr := json.Marshal(err)
	if jsonErr != nil {
		t.Fatalf("Failed to marshal error: %v", jsonErr)
	}

	var result map[string]any
	if jsonErr := json.Unmarshal(data, &result); jsonErr != nil {
		t.Fatalf("Failed to unmarshal error: %v", jsonErr)
	}

	if result["key"] != "test.error" {
		t.Errorf("Expected key 'test.error', got '%v'", result["key"])
	}

	if result["request_id"] != "req_123" {
		t.Errorf("Expected request_id 'req_123', got '%v'", result["request_id"])
	}

	if result["retryable"] != true {
		t.Errorf("Expected retryable to be true, got '%v'", result["retryable"])
	}
}

func TestToMap(t *testing.T) {
	err := New("test.error", "Test", 400, CategoryValidation).
		WithOptions(WithRequestID("req_123"))

	m := err.ToMap()

	if m["id"] == nil {
		t.Error("Expected 'id' in map")
	}

	if m["key"] != "test.error" {
		t.Errorf("Expected key 'test.error', got '%v'", m["key"])
	}

	if m["request_id"] != "req_123" {
		t.Errorf("Expected request_id 'req_123', got '%v'", m["request_id"])
	}
}

func TestCause(t *testing.T) {
	rootErr := errors.New("root")
	middleErr := Wrap(rootErr, "middle", "Middle error", 500, CategoryInternal)
	topErr := Wrap(middleErr, "top", "Top error", 500, CategoryInternal)

	cause := topErr.Cause()
	if cause != rootErr {
		t.Error("Expected Cause to return root error")
	}
}

func TestStackTrace(t *testing.T) {
	err := New("test.error", "Test", 500, CategoryInternal)

	if len(err.stackTrace) == 0 {
		t.Error("Expected stack trace to be captured")
	}

	// Stack trace should contain function names
	found := false
	for _, frame := range err.stackTrace {
		if frame != "" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected non-empty stack trace frames")
	}
}

func BenchmarkNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = New("test.error", "Test error", 400, CategoryValidation)
	}
}

func BenchmarkWithOptions(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = New("test.error", "Test", 400, CategoryValidation).
			WithOptions(
				WithRequestID("req_123"),
				WithTraceID("trace_456"),
				WithRetryable(true),
			)
	}
}

func BenchmarkMarshalJSON(b *testing.B) {
	err := New("test.error", "Test", 400, CategoryValidation).
		WithOptions(WithRequestID("req_123"))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = json.Marshal(err)
	}
}
