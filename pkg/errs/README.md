# Error Handling Package

Production-grade error handling system for audoctl. This package provides structured, observable, and traceable error responses following RFC 7807 (Problem Details) standard with extensions for distributed systems and AI agent observability.

## Features

✅ **Structured Error Responses** - Consistent error format across the API  
✅ **Error Categories** - Typed error categories for better handling  
✅ **Request Tracing** - Request ID, Trace ID, Session ID support  
✅ **Error Wrapping** - Full error chain support with `errors.Unwrap`  
✅ **Validation Errors** - Built-in validation error handling  
✅ **Debug Information** - Stack traces and error chains in development  
✅ **Retryable Hints** - Indicate whether errors are retryable  
✅ **Metadata Support** - Attach custom metadata to errors  
✅ **Fiber Integration** - Ready-to-use Fiber middleware and handlers  

## Quick Start

### Basic Usage

```go
import "github.com/audoctl/audoctl/pkg/errs"

// Simple error
err := errs.NotFound("session", "sess_123")

// With context
err := errs.SessionNotFound(sessionID).
    WithOptions(
        errs.WithRequestID(requestID),
        errs.WithSessionID(sessionID),
    )
```

### Error Categories

```go
// Validation errors
err := errs.Validation("Invalid request",
    errs.NewValidationError("agent.required", "agent", "Agent name is required", nil),
)

// Not found
err := errs.NotFound("session", sessionID)

// Authorization
err := errs.Unauthorized("Invalid API key")
err := errs.Forbidden("Insufficient permissions")

// Server errors
err := errs.Internal("Database connection failed", dbErr)
err := errs.External("OpenAI API", apiErr)

// Rate limiting
err := errs.RateLimit(30 * time.Second)

// Timeout
err := errs.Timeout("llm_call", 30*time.Second)
```

### Error Wrapping

```go
// Root cause
dbErr := sql.ErrNoRows

// Wrap with context
err := errs.Wrap(
    dbErr,
    "session.fetch_failed",
    "Failed to fetch session from database",
    500,
    errs.CategoryInternal,
).WithOptions(
    errs.WithSessionID(sessionID),
    errs.WithRetryable(true),
)

// Access underlying error
rootCause := err.Cause()
```

### Validation Errors

```go
err := errs.Validation("Request validation failed")

// Add multiple validation errors
err.AddValidation(errs.NewValidationError(
    "agent.required",
    "agent",
    "Agent name is required",
    nil,
))

err.AddValidation(errs.NewValidationError(
    "timeout.invalid",
    "timeout",
    "Timeout must be positive",
    map[string]any{"provided": -5},
))
```

### Custom Errors

```go
err := errs.New(
    "agent.execution_failed",
    "AI agent execution failed",
    500,
    errs.CategoryInternal,
).SetParams(map[string]any{
    "agent": "refund_agent",
    "step":  "tool_execution",
}).SetMetadata("cost_usd", 0.05).
  SetMetadata("tokens_used", 1234).
  WithOptions(
    errs.WithSessionID(sessionID),
    errs.WithRetryable(true),
    errs.WithHelpURL("https://docs.audoctl.dev/errors/agent-execution"),
)
```

## Fiber Integration

### Error Handler Middleware

```go
import (
    "github.com/audoctl/audoctl/pkg/errs"
    "github.com/gofiber/fiber/v3"
)

app := fiber.New(fiber.Config{
    ErrorHandler: errs.NewFiberHandler(isDevelopment).Middleware(),
})
```

### In HTTP Handlers

```go
func GetSession(c fiber.Ctx) error {
    sessionID := c.Params("id")
    
    session, err := repo.FindByID(sessionID)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return errs.SessionNotFound(sessionID).
                WithOptions(errs.ExtractRequestContext(c)...)
        }
        
        return errs.DatabaseError("fetch session", err).
            WithOptions(errs.ExtractRequestContext(c)...)
    }
    
    return response.Success(c, session)
}
```

### Using Response Builder

```go
import "github.com/audoctl/audoctl/pkg/fiberserver/response"

// Success response
return response.Success(c, data)

// Created response
return response.Created(c, newSession)

// Error response
return response.Error(c, err, isDevelopment)

// With pagination
return response.WithPagination(c, sessions, page, perPage, total)
```

### Recovery Middleware

```go
app.Use(errs.RecoverMiddleware())
```

## Error Response Format

### Success Response

```json
{
  "success": true,
  "data": {
    "id": "sess_123",
    "agent": "refund_agent",
    "status": "active"
  }
}
```

### Error Response

```json
{
  "success": false,
  "error": {
    "id": "err_1708123456789",
    "key": "session.not_found",
    "category": "not_found",
    "message": "session not found",
    "status": 404,
    "timestamp": "2024-02-16T12:34:56Z",
    "request_id": "req_abc123",
    "session_id": "sess_123",
    "params": {
      "resource": "session",
      "id": "sess_123"
    },
    "retryable": false
  }
}
```

### Validation Error Response

```json
{
  "success": false,
  "error": {
    "id": "err_1708123456789",
    "key": "validation.failed",
    "category": "validation",
    "message": "Request validation failed",
    "status": 400,
    "timestamp": "2024-02-16T12:34:56Z",
    "validations": [
      {
        "key": "agent.required",
        "field": "agent",
        "message": "Agent name is required"
      },
      {
        "key": "timeout.invalid",
        "field": "timeout",
        "message": "Timeout must be positive",
        "params": {
          "provided": -5
        }
      }
    ],
    "retryable": false
  }
}
```

### Error with Debug Info (Development Only)

```json
{
  "success": false,
  "error": {
    "id": "err_1708123456789",
    "key": "internal.error",
    "message": "Database connection failed: connection timeout",
    "status": 500,
    "debug": {
      "stack_trace": [
        "/app/pkg/errs/error.go:123 New",
        "/app/internal/repository/session.go:45 FindByID",
        "/app/internal/handler/session.go:67 GetSession"
      ],
      "cause": "connection timeout",
      "error_chain": [
        "Database connection failed: connection timeout",
        "connection timeout"
      ]
    }
  }
}
```

## Error Categories

| Category | HTTP Status | Use Case |
|----------|-------------|----------|
| `validation` | 400 | Input validation failures |
| `not_found` | 404 | Resource not found |
| `unauthorized` | 401 | Authentication required |
| `forbidden` | 403 | Insufficient permissions |
| `conflict` | 409 | Resource state conflict |
| `internal` | 500 | Internal server errors |
| `external` | 502 | External service failures |
| `timeout` | 408 | Operation timeout |
| `rate_limit` | 429 | Rate limit exceeded |

## Best Practices

### 1. Use Factory Functions

Prefer factory functions over direct construction:

```go
// Good
err := errs.NotFound("session", id)

// Avoid
err := errs.New("session.not_found", "not found", 404, errs.CategoryNotFound)
```

### 2. Add Context

Always add request context when available:

```go
err.WithOptions(
    errs.WithRequestID(requestID),
    errs.WithTraceID(traceID),
    errs.WithSessionID(sessionID),
)

// Or use helper in Fiber
err.WithOptions(errs.ExtractRequestContext(c)...)
```

### 3. Wrap Errors

Wrap lower-level errors with business context:

```go
if err := repo.Save(session); err != nil {
    return errs.DatabaseError("save session", err)
}
```

### 4. Use Validation Errors

For input validation, use validation errors:

```go
err := errs.Validation("Invalid request")
if agent == "" {
    err.AddValidation(errs.NewValidationError(
        "agent.required",
        "agent",
        "Agent name is required",
        nil,
    ))
}
```

### 5. Handle Retryable Errors

Mark retryable errors appropriately:

```go
return errs.External("OpenAI", err).
    WithOptions(errs.WithRetryable(true))

// Or with retry delay
return errs.RateLimit(30 * time.Second)
```

### 6. Debug in Development Only

```go
err.WithDebugInfo(os.Getenv("ENV") != "production")
```

### 7. Add Help URLs

```go
err.WithOptions(errs.WithHelpURL(
    "https://docs.audoctl.dev/errors/session-not-found",
))
```

## Testing

```go
func TestErrorCreation(t *testing.T) {
    err := errs.NotFound("session", "sess_123")
    
    assert.Equal(t, "session.not_found", err.Key)
    assert.Equal(t, 404, err.Status)
    assert.Equal(t, errs.CategoryNotFound, err.Category)
}

func TestErrorWrapping(t *testing.T) {
    rootErr := errors.New("connection failed")
    err := errs.DatabaseError("fetch", rootErr)
    
    assert.Equal(t, rootErr, errors.Unwrap(err))
    assert.Equal(t, rootErr, err.Cause())
}

func TestErrorComparison(t *testing.T) {
    err1 := errs.NotFound("session", "123")
    err2 := errs.NotFound("session", "456")
    
    assert.True(t, errors.Is(err1, err2))
}
```

## Migration from Old Error System

If you have existing error code:

```go
// Old
e := &errs.ErrorResponse{
    Key:     "session.not_found",
    Message: "not found",
    Status:  404,
}

// New
e := errs.NotFound("session", sessionID)
```

## Advanced Usage

### Custom Error Categories

```go
const (
    CategoryAIFailure errs.ErrorCategory = "ai_failure"
    CategoryTokenLimit errs.ErrorCategory = "token_limit"
)

err := errs.New(
    "llm.context_too_long",
    "Context exceeds token limit",
    400,
    CategoryTokenLimit,
).SetParams(map[string]any{
    "tokens_provided": 50000,
    "tokens_limit":    32000,
})
```

### Structured Logging

```go
log.Error().
    Fields(err.ToMap()).
    Msg("Operation failed")
```

### Error Metrics

```go
if e, ok := errs.ToErr(err); ok {
    metrics.ErrorCount.
        WithLabelValues(string(e.Category), e.Key).
        Inc()
}
```

## License

Part of audoctl project - MIT License
