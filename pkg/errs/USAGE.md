# Quick Start Guide

## Basic Usage

### 1. Simple Error Creation

```go
// Not found error
err := errs.SessionNotFound("sess_123")

// Validation error
err := errs.Validation("Invalid request",
    errs.NewValidationError("agent.required", "agent", "Agent name is required", nil),
)

// Internal error
err := errs.Internal("Database error", dbErr)
```

### 2. Fiber Handler'da Kullanım

```go
package handler

import (
    "github.com/audoctl/audoctl/pkg/errs"
    "github.com/audoctl/audoctl/pkg/fiberserver/response"
    "github.com/gofiber/fiber/v3"
)

func GetSession(c fiber.Ctx) error {
    sessionID := c.Params("id")

    session, err := sessionService.GetByID(sessionID)
    if err != nil {
        // Database error
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

### 3. Fiber App Setup

```go
package main

import (
    "os"
    "github.com/audoctl/audoctl/pkg/errs"
    "github.com/gofiber/fiber/v3"
)

func main() {
    isDevelopment := os.Getenv("ENV") != "production"

    app := fiber.New(fiber.Config{
        ErrorHandler: errs.NewFiberHandler(isDevelopment).Middleware(),
    })

    // Recovery middleware
    app.Use(errs.RecoverMiddleware())

    // Your routes
    app.Get("/v1/sessions/:id", GetSession)

    app.Listen(":3000")
}
```

### 4. Validation Errors

```go
func CreateSession(c fiber.Ctx) error {
    var req CreateSessionRequest
    if err := c.BodyParser(&req); err != nil {
        return errs.BadRequest("Invalid JSON")
    }

    // Validate
    validationErr := errs.Validation("Request validation failed")
    hasErrors := false

    if req.Agent == "" {
        validationErr.AddValidation(errs.NewValidationError(
            "agent.required",
            "agent",
            "Agent name is required",
            nil,
        ))
        hasErrors = true
    }

    if req.Timeout < 0 {
        validationErr.AddValidation(errs.NewValidationError(
            "timeout.invalid",
            "timeout",
            "Timeout must be positive",
            map[string]any{"provided": req.Timeout},
        ))
        hasErrors = true
    }

    if hasErrors {
        return validationErr.WithOptions(errs.ExtractRequestContext(c)...)
    }

    // Create session
    session, err := sessionService.Create(req)
    if err != nil {
        return errs.Internal("Failed to create session", err)
    }

    return response.Created(c, session)
}
```

### 5. On Service Layer Error Wrapping

```go
package service

import (
    "github.com/audoctl/audoctl/pkg/errs"
)

func (s *SessionService) GetByID(id string) (*Session, error) {
    session, err := s.repo.FindByID(id)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, errs.NotFound("session", id)
        }

        // Wrap database error with context
        return nil, errs.DatabaseError("fetch session", err)
    }

    return session, nil
}

func (s *SessionService) FinishSession(id string) error {
    session, err := s.GetByID(id)
    if err != nil {
        return err // Already wrapped
    }

    if session.Status == "finished" {
        return errs.SessionAlreadyFinished(id)
    }

    session.Status = "finished"
    if err := s.repo.Update(session); err != nil {
        return errs.StorageError("update session", err)
    }

    return nil
}
```

### 6. External Service Errors

```go
func (s *LLMService) Call(prompt string) (string, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    response, err := s.client.Complete(ctx, prompt)
    if err != nil {
        // Timeout
        if errors.Is(err, context.DeadlineExceeded) {
            return "", errs.Timeout("llm_call", 30*time.Second)
        }

        // External service error (retryable)
        return "", errs.External("OpenAI", err)
    }

    return response, nil
}
```

### 7. Response Samples

#### Success Response
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

#### Error Response
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
    "params": {
      "resource": "session",
      "id": "sess_123"
    },
    "retryable": false
  }
}
```

#### Validation Error Response
```json
{
  "success": false,
  "error": {
    "id": "err_1708123456790",
    "key": "validation.failed",
    "category": "validation",
    "message": "Request validation failed",
    "status": 400,
    "timestamp": "2024-02-16T12:35:00Z",
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

## Middleware Example

### Request ID Middleware

```go
func RequestIDMiddleware() fiber.Handler {
    return func(c fiber.Ctx) error {
        requestID := c.Get("X-Request-ID")
        if requestID == "" {
            requestID = fmt.Sprintf("req_%d", time.Now().UnixNano())
        }
        c.Locals("requestid", requestID)
        c.Set("X-Request-ID", requestID)
        return c.Next()
    }
}
```

### Logging Middleware

```go
func LoggingMiddleware() fiber.Handler {
    return func(c fiber.Ctx) error {
        start := time.Now()

        err := c.Next()

        duration := time.Since(start)
        status := c.Response().StatusCode()

        log := logger.Info()
        if status >= 400 {
            log = logger.Error()

            // Log error details if available
            if e, ok := errs.ToErr(err); ok {
                log = log.Fields(e.ToMap())
            }
        }

        log.
            Str("method", c.Method()).
            Str("path", c.Path()).
            Int("status", status).
            Dur("duration", duration).
            Msg("HTTP request")

        return err
    }
}
```

## Best Practices

### 1. Always Add Context

```go
// GOOD ✓
return errs.SessionNotFound(id).
    WithOptions(errs.ExtractRequestContext(c)...)

// BAD ✗
return errs.SessionNotFound(id)
```

### 2. Use Error Wrapping

```go
// GOOD ✓
if err := repo.Save(session); err != nil {
    return errs.DatabaseError("save session", err)
}

// BAD ✗
if err := repo.Save(session); err != nil {
    return err
}
```

### 3. Mark Retryable

```go
// GOOD ✓
return errs.External("OpenAI", err).
    WithOptions(errs.WithRetryable(true))

// FOR RATELIMIT
return errs.RateLimit(30 * time.Second)
```

### 4. In Development Add Debug Info

```go
isDevelopment := os.Getenv("ENV") != "production"
err.WithDebugInfo(isDevelopment)
```

### 5. Use Factory Functions

```go
// GOOD ✓
err := errs.NotFound("session", id)

// BAD ✗
err := errs.New("session.not_found", "not found", 404, errs.CategoryNotFound)
```

## Testing

```go
func TestHandler(t *testing.T) {
    app := fiber.New()
    app.Get("/sessions/:id", GetSession)

    req := httptest.NewRequest("GET", "/sessions/invalid", nil)
    resp, _ := app.Test(req)

    assert.Equal(t, 404, resp.StatusCode)

    var body response.Response
    json.NewDecoder(resp.Body).Decode(&body)

    assert.False(t, body.Success)
    assert.Equal(t, "session.not_found", body.Error.Key)
}
```
