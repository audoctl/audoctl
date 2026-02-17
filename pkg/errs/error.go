package errs

import (
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
	"time"
)

// ErrorCategory represents the category of error for better handling
type ErrorCategory string

const (
	CategoryValidation   ErrorCategory = "validation"
	CategoryNotFound     ErrorCategory = "not_found"
	CategoryUnauthorized ErrorCategory = "unauthorized"
	CategoryForbidden    ErrorCategory = "forbidden"
	CategoryConflict     ErrorCategory = "conflict"
	CategoryInternal     ErrorCategory = "internal"
	CategoryExternal     ErrorCategory = "external"
	CategoryTimeout      ErrorCategory = "timeout"
	CategoryRateLimit    ErrorCategory = "rate_limit"
)

// Err is production-grade error type that implements error interface
// It follows RFC 7807 (Problem Details) standard with extensions for observability
type Err struct {
	// Core fields
	ID        string        `json:"id"`        // Unique error instance ID
	Key       string        `json:"key"`       // Error key for i18n/lookup (e.g., "session.not_found")
	Category  ErrorCategory `json:"category"`  // Error category for handling
	Message   string        `json:"message"`   // Human-readable error message
	Status    int           `json:"status"`    // HTTP status code
	Timestamp time.Time     `json:"timestamp"` // When the error occurred

	// Context & tracing
	RequestID string `json:"request_id,omitempty"` // Request ID for distributed tracing
	TraceID   string `json:"trace_id,omitempty"`   // Trace ID for observability
	SessionID string `json:"session_id,omitempty"` // AI session ID (specific to audoctl)

	// Error details
	Params      map[string]any    `json:"params,omitempty"`      // Dynamic parameters for message interpolation
	Validations []ValidationError `json:"validations,omitempty"` // Validation errors
	Metadata    map[string]any    `json:"metadata,omitempty"`    // Additional context metadata

	// Debugging & support
	Debug   *DebugInfo `json:"debug,omitempty"`    // Debug information (only in dev mode)
	HelpURL string     `json:"help_url,omitempty"` // Documentation/help URL

	// Error handling hints
	Retryable  bool           `json:"retryable"`             // Whether the operation can be retried
	RetryAfter *time.Duration `json:"retry_after,omitempty"` // Suggested retry delay

	// Internal fields
	cause      error    `json:"-"` // Wrapped/underlying error
	stackTrace []string `json:"-"` // Stack trace for internal logging
}

// DebugInfo contains debugging information (only exposed in development)
type DebugInfo struct {
	StackTrace []string `json:"stack_trace,omitempty"`
	Cause      string   `json:"cause,omitempty"`
	ErrorChain []string `json:"error_chain,omitempty"`
	File       string   `json:"file,omitempty"`
	Line       int      `json:"line,omitempty"`
}

// ErrorOption is a functional option for error construction
type ErrorOption func(*Err)

// New creates a new production-grade error
func New(key string, message string, status int, category ErrorCategory) *Err {
	e := &Err{
		ID:          generateErrorID(),
		Key:         key,
		Category:    category,
		Message:     message,
		Status:      status,
		Timestamp:   time.Now().UTC(),
		Validations: []ValidationError{},
		Params:      make(map[string]any),
		Metadata:    make(map[string]any),
		Retryable:   false,
		stackTrace:  captureStackTrace(3),
	}
	return e
}

// Wrap wraps an existing error with additional context
func Wrap(err error, key string, message string, status int, category ErrorCategory) *Err {
	if err == nil {
		return nil
	}

	e := New(key, message, status, category)
	e.cause = err

	// If wrapping another Err, preserve some context
	if existingErr, ok := err.(*Err); ok {
		if e.RequestID == "" {
			e.RequestID = existingErr.RequestID
		}
		if e.TraceID == "" {
			e.TraceID = existingErr.TraceID
		}
		if e.SessionID == "" {
			e.SessionID = existingErr.SessionID
		}
	}

	return e
}

// WithOptions applies functional options to the error
func (e *Err) WithOptions(opts ...ErrorOption) *Err {
	for _, opt := range opts {
		opt(e)
	}
	return e
}

// Error implements the error interface
func (e *Err) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.cause)
	}
	return e.Message
}

// Unwrap returns the underlying error for errors.Unwrap
func (e *Err) Unwrap() error {
	return e.cause
}

// Cause returns the root cause error
func (e *Err) Cause() error {
	err := error(e)
	for {
		unwrapped := errors.Unwrap(err)
		if unwrapped == nil {
			return err
		}
		err = unwrapped
	}
}

// AddValidation adds a validation error
func (e *Err) AddValidation(validation ValidationError) *Err {
	e.Validations = append(e.Validations, validation)
	if e.Category != CategoryValidation {
		e.Category = CategoryValidation
	}
	return e
}

// SetParam sets a parameter for message interpolation
func (e *Err) SetParam(key string, value any) *Err {
	if e.Params == nil {
		e.Params = make(map[string]any)
	}
	e.Params[key] = value
	return e
}

// SetParams sets multiple parameters
func (e *Err) SetParams(params map[string]any) *Err {
	if e.Params == nil {
		e.Params = make(map[string]any)
	}
	for k, v := range params {
		e.Params[k] = v
	}
	return e
}

// SetMetadata sets a metadata field
func (e *Err) SetMetadata(key string, value any) *Err {
	if e.Metadata == nil {
		e.Metadata = make(map[string]any)
	}
	e.Metadata[key] = value
	return e
}

// WithDebugInfo adds debug information to the error
func (e *Err) WithDebugInfo(includeDebug bool) *Err {
	if !includeDebug {
		return e
	}

	debug := &DebugInfo{
		StackTrace: e.stackTrace,
	}

	if e.cause != nil {
		debug.Cause = e.cause.Error()
		debug.ErrorChain = e.buildErrorChain()
	}

	if len(e.stackTrace) > 0 {
		// Parse file and line from first stack frame
		if frame := e.stackTrace[0]; frame != "" {
			debug.File = frame
		}
	}

	e.Debug = debug
	return e
}

// ToMap converts error to a map for logging/serialization
func (e *Err) ToMap() map[string]any {
	m := map[string]any{
		"id":        e.ID,
		"key":       e.Key,
		"category":  e.Category,
		"message":   e.Message,
		"status":    e.Status,
		"timestamp": e.Timestamp,
		"retryable": e.Retryable,
	}

	if e.RequestID != "" {
		m["request_id"] = e.RequestID
	}
	if e.TraceID != "" {
		m["trace_id"] = e.TraceID
	}
	if e.SessionID != "" {
		m["session_id"] = e.SessionID
	}
	if len(e.Params) > 0 {
		m["params"] = e.Params
	}
	if len(e.Validations) > 0 {
		m["validations"] = e.Validations
	}
	if len(e.Metadata) > 0 {
		m["metadata"] = e.Metadata
	}
	if e.Debug != nil {
		m["debug"] = e.Debug
	}
	if e.HelpURL != "" {
		m["help_url"] = e.HelpURL
	}
	if e.RetryAfter != nil {
		m["retry_after"] = e.RetryAfter.Seconds()
	}

	return m
}

// MarshalJSON implements json.Marshaler
func (e *Err) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.ToMap())
}

// As implements errors.As target
func (e *Err) As(target any) bool {
	if err, ok := target.(**Err); ok {
		*err = e
		return true
	}
	return false
}

// Is implements errors.Is comparison
func (e *Err) Is(target error) bool {
	if t, ok := target.(*Err); ok {
		return e.Key == t.Key && e.Status == t.Status
	}
	return false
}

// Helper functions

// ToErr converts a standard error to *Err if possible
func ToErr(err error) (*Err, bool) {
	if err == nil {
		return nil, false
	}
	var e *Err
	if errors.As(err, &e) {
		return e, true
	}
	return nil, false
}

// buildErrorChain builds a chain of error messages
func (e *Err) buildErrorChain() []string {
	var chain []string
	err := error(e)

	for err != nil {
		chain = append(chain, err.Error())
		err = errors.Unwrap(err)
	}

	return chain
}

// captureStackTrace captures the current stack trace
func captureStackTrace(skip int) []string {
	const maxDepth = 32
	var pcs [maxDepth]uintptr
	n := runtime.Callers(skip, pcs[:])

	frames := runtime.CallersFrames(pcs[:n])
	var trace []string

	for {
		frame, more := frames.Next()
		trace = append(trace, fmt.Sprintf("%s:%d %s", frame.File, frame.Line, frame.Function))
		if !more {
			break
		}
	}

	return trace
}

// generateErrorID generates a unique error instance ID
func generateErrorID() string {
	return fmt.Sprintf("err_%d", time.Now().UnixNano())
}

// Functional options

func WithRequestID(requestID string) ErrorOption {
	return func(e *Err) {
		e.RequestID = requestID
	}
}

func WithTraceID(traceID string) ErrorOption {
	return func(e *Err) {
		e.TraceID = traceID
	}
}

func WithSessionID(sessionID string) ErrorOption {
	return func(e *Err) {
		e.SessionID = sessionID
	}
}

func WithCause(cause error) ErrorOption {
	return func(e *Err) {
		e.cause = cause
	}
}

func WithRetryable(retryable bool) ErrorOption {
	return func(e *Err) {
		e.Retryable = retryable
	}
}

func WithRetryAfter(duration time.Duration) ErrorOption {
	return func(e *Err) {
		e.RetryAfter = &duration
		e.Retryable = true
	}
}

func WithHelpURL(url string) ErrorOption {
	return func(e *Err) {
		e.HelpURL = url
	}
}

func WithParams(params map[string]any) ErrorOption {
	return func(e *Err) {
		e.SetParams(params)
	}
}

func WithMetadata(metadata map[string]any) ErrorOption {
	return func(e *Err) {
		for k, v := range metadata {
			e.SetMetadata(k, v)
		}
	}
}
