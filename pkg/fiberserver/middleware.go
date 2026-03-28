package fiberserver

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"
	"strings"
	"time"

	"github.com/audoctl/audoctl/pkg/errs"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

// DetectLanguageMiddleware extracts Accept-Language header and stores it in context
func DetectLanguageMiddleware(c fiber.Ctx) error {
	lang := c.Get(fiber.HeaderAcceptLanguage)
	if lang == "" {
		lang = "en-US"
	}
	c.SetContext(context.WithValue(c.Context(), "lang", lang))
	return c.Next()
}

// RecoverMiddleware recovers from panics and returns structured error
func RecoverMiddleware(includeDebug bool) fiber.Handler {
	return func(c fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				var err error

				switch v := r.(type) {
				case error:
					err = v
				case string:
					err = errors.New(v)
				default:
					err = fmt.Errorf("panic: %v", v)
				}

				// Create error with stack trace
				e := errs.Internal("Panic recovered", err).
					WithOptions(errs.ExtractRequestContext(c)...).
					WithDebugInfo(includeDebug)

				// Add stack trace to metadata if in debug mode
				if includeDebug {
					e.SetMetadata("stack_trace", string(debug.Stack()))
				}

				_ = c.Status(e.Status).JSON(e)
			}
		}()

		return c.Next()
	}
}

// RequestIDMiddleware generates and attaches request ID to each request
func RequestIDMiddleware(header string) fiber.Handler {
	if header == "" {
		header = "X-Request-ID"
	}

	return func(c fiber.Ctx) error {
		requestID := c.Get(header)

		// Generate new request ID if not present
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Store in locals and set response header
		c.Locals("requestid", requestID)
		c.Set(header, requestID)

		return c.Next()
	}
}

// TraceIDMiddleware extracts or generates trace ID for distributed tracing
func TraceIDMiddleware(header string) fiber.Handler {
	if header == "" {
		header = "X-Trace-ID"
	}

	return func(c fiber.Ctx) error {
		traceID := c.Get(header)

		// Generate new trace ID if not present
		if traceID == "" {
			traceID = uuid.New().String()
		}

		// Store in locals and set response header
		c.Locals("traceid", traceID)
		c.Set(header, traceID)

		return c.Next()
	}
}

// SecurityHeadersMiddleware adds security headers to responses
func SecurityHeadersMiddleware(cfg SecurityHeadersConfig) fiber.Handler {
	return func(c fiber.Ctx) error {
		// XSS Protection
		if cfg.XSSProtection != "" {
			c.Set("X-XSS-Protection", cfg.XSSProtection)
		}

		// Content Type Nosniff
		if cfg.ContentTypeNosniff != "" {
			c.Set("X-Content-Type-Options", cfg.ContentTypeNosniff)
		}

		// X-Frame-Options
		if cfg.XFrameOptions != "" {
			c.Set("X-Frame-Options", cfg.XFrameOptions)
		}

		// HSTS
		if cfg.HSTSMaxAge > 0 {
			hsts := fmt.Sprintf("max-age=%d", cfg.HSTSMaxAge)
			if cfg.HSTSIncludeSubdomains {
				hsts += "; includeSubDomains"
			}
			c.Set("Strict-Transport-Security", hsts)
		}

		// Content Security Policy
		if cfg.ContentSecurityPolicy != "" {
			c.Set("Content-Security-Policy", cfg.ContentSecurityPolicy)
		}

		// Referrer Policy
		if cfg.ReferrerPolicy != "" {
			c.Set("Referrer-Policy", cfg.ReferrerPolicy)
		}

		return c.Next()
	}
}

// SecurityHeadersConfig contains security headers configuration
type SecurityHeadersConfig struct {
	XSSProtection         string
	ContentTypeNosniff    string
	XFrameOptions         string
	HSTSMaxAge            int
	HSTSIncludeSubdomains bool
	ContentSecurityPolicy string
	ReferrerPolicy        string
}

// APIKeyAuthMiddleware validates API key from header
func APIKeyAuthMiddleware(apiKeys []string, header string) fiber.Handler {
	if header == "" {
		header = "X-API-Key"
	}

	// Create map for O(1) lookup
	keyMap := make(map[string]bool)
	for _, key := range apiKeys {
		keyMap[key] = true
	}

	return func(c fiber.Ctx) error {
		apiKey := c.Get(header)

		if apiKey == "" {
			return errs.Unauthorized("API key is required").
				WithOptions(errs.ExtractRequestContext(c)...)
		}

		if !keyMap[apiKey] {
			return errs.Unauthorized("Invalid API key").
				WithOptions(errs.ExtractRequestContext(c)...)
		}

		return c.Next()
	}
}

// SessionIDMiddleware extracts session ID from path or query and stores it in locals
func SessionIDMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		// Try path parameter first
		sessionID := c.Params("session_id")

		// Try query parameter if not in path
		if sessionID == "" {
			sessionID = c.Query("session_id")
		}

		// Store in locals if present
		if sessionID != "" {
			c.Locals("sessionid", sessionID)
		}

		return c.Next()
	}
}

// TimeoutMiddleware adds request timeout
func TimeoutMiddleware(timeout time.Duration) fiber.Handler {
	return func(c fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(c.Context(), timeout)
		defer cancel()

		c.SetContext(ctx)

		// Channel to signal completion
		done := make(chan error, 1)

		go func() {
			done <- c.Next()
		}()

		select {
		case err := <-done:
			return err
		case <-ctx.Done():
			return errs.Timeout("request", timeout).
				WithOptions(errs.ExtractRequestContext(c)...)
		}
	}
}

// LoggingMiddleware logs HTTP requests (basic implementation - DEPRECATED, use StructuredLoggingMiddleware)
func LoggingMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		start := time.Now()

		// Process request
		err := c.Next()

		// Log after request
		duration := time.Since(start)
		status := c.Response().StatusCode()

		// You can integrate with your logger here
		fmt.Printf("[HTTP] %s %s - %d (%v)\n",
			c.Method(),
			c.Path(),
			status,
			duration,
		)

		return err
	}
}

// SkipPathMiddleware skips middleware for specific paths
func SkipPathMiddleware(paths []string, middleware fiber.Handler) fiber.Handler {
	// Create map for O(1) lookup
	skipMap := make(map[string]bool)
	for _, path := range paths {
		skipMap[path] = true
	}

	return func(c fiber.Ctx) error {
		path := c.Path()

		// Check exact match
		if skipMap[path] {
			return c.Next()
		}

		// Check prefix match for wildcard paths
		for skipPath := range skipMap {
			if strings.HasSuffix(skipPath, "*") {
				prefix := strings.TrimSuffix(skipPath, "*")
				if strings.HasPrefix(path, prefix) {
					return c.Next()
				}
			}
		}

		return middleware(c)
	}
}

// StructuredLoggingConfig contains configuration for structured HTTP logging
type StructuredLoggingConfig struct {
	Logger interface {
		Info(msg string, args map[string]any)
		Warn(msg string, args map[string]any)
		Error(msg string, args map[string]any)
		Debug(msg string, args map[string]any)
	}
	LogRequestBody  bool
	LogResponseBody bool
	SkipPaths       []string
	SlowThreshold   time.Duration
	MaxBodyLogSize  int
}

// StructuredLoggingMiddleware provides production-ready HTTP request/response logging
func StructuredLoggingMiddleware(cfg StructuredLoggingConfig) fiber.Handler {
	// Build skip path map for O(1) lookup
	skipMap := make(map[string]bool)
	for _, path := range cfg.SkipPaths {
		skipMap[path] = true
	}

	return func(c fiber.Ctx) error {
		// Check if path should be skipped
		path := c.Path()
		if shouldSkipPath(path, skipMap) {
			return c.Next()
		}

		start := time.Now()

		// Capture request information
		fields := buildRequestFields(c, cfg)

		// Process request
		err := c.Next()

		// Calculate duration
		duration := time.Since(start)
		fields["duration_ms"] = duration.Milliseconds()
		fields["duration"] = duration.String()

		// Capture response information
		buildResponseFields(c, fields, cfg)

		// Determine log level and message
		status := c.Response().StatusCode()
		logLevel, message := determineLogLevel(status, duration, cfg.SlowThreshold, err)

		// Add error if present
		if err != nil {
			fields["error"] = err.Error()
			fields["error_type"] = fmt.Sprintf("%T", err)
		}

		// Log based on level
		switch logLevel {
		case "debug":
			cfg.Logger.Debug(message, fields)
		case "error":
			cfg.Logger.Error(message, fields)
		case "warn":
			cfg.Logger.Warn(message, fields)
		default:
			cfg.Logger.Info(message, fields)
		}

		return err
	}
}

// buildRequestFields creates structured fields from the request
func buildRequestFields(c fiber.Ctx, cfg StructuredLoggingConfig) map[string]any {
	fields := map[string]any{
		"component":  "http",
		"type":       "request",
		"method":     c.Method(),
		"path":       c.Path(),
		"route":      c.Route().Path,
		"ip":         c.IP(),
		"user_agent": c.Get(fiber.HeaderUserAgent),
		"protocol":   c.Protocol(),
	}

	// Add request ID if present
	if reqID, ok := c.Locals("requestid").(string); ok && reqID != "" {
		fields["request_id"] = reqID
	}

	// Add trace ID if present
	if traceID, ok := c.Locals("traceid").(string); ok && traceID != "" {
		fields["trace_id"] = traceID
	}

	// Add session ID if present
	if sessionID, ok := c.Locals("sessionid").(string); ok && sessionID != "" {
		fields["session_id"] = sessionID
	}

	// Add correlation ID from context
	if corrID := c.Context().Value("correlation-id"); corrID != nil {
		if corrIDStr, ok := corrID.(string); ok && corrIDStr != "" {
			fields["correlation_id"] = corrIDStr
		}
	}

	// Add hostname
	fields["hostname"] = c.Hostname()

	// Add query parameters if present
	if queryString := string(c.Request().URI().QueryString()); queryString != "" {
		fields["query"] = queryString
	}

	// Add request headers (selected)
	fields["headers"] = map[string]string{
		"content_type":   c.Get(fiber.HeaderContentType),
		"content_length": c.Get(fiber.HeaderContentLength),
		"accept":         c.Get(fiber.HeaderAccept),
		"referer":        c.Get(fiber.HeaderReferer),
	}

	// Add request body if configured
	if cfg.LogRequestBody && len(c.Body()) > 0 {
		bodySize := len(c.Body())
		fields["request_body_size"] = bodySize

		if bodySize <= cfg.MaxBodyLogSize {
			fields["request_body"] = string(c.Body())
		} else {
			fields["request_body"] = string(c.Body()[:cfg.MaxBodyLogSize]) + "... (truncated)"
			fields["request_body_truncated"] = true
		}
	}

	return fields
}

// buildResponseFields adds response information to fields
func buildResponseFields(c fiber.Ctx, fields map[string]any, cfg StructuredLoggingConfig) {
	status := c.Response().StatusCode()
	fields["status"] = status
	fields["status_class"] = fmt.Sprintf("%dxx", status/100)

	// Add response size
	responseSize := len(c.Response().Body())
	fields["response_size"] = responseSize

	// Add response body if configured
	if cfg.LogResponseBody && responseSize > 0 {
		if responseSize <= cfg.MaxBodyLogSize {
			fields["response_body"] = string(c.Response().Body())
		} else {
			fields["response_body"] = string(c.Response().Body()[:cfg.MaxBodyLogSize]) + "... (truncated)"
			fields["response_body_truncated"] = true
		}
	}
}

// determineLogLevel determines the appropriate log level based on status and duration
func determineLogLevel(status int, duration time.Duration, slowThreshold time.Duration, err error) (string, string) {
	// Error responses (5xx) or errors
	if status >= 500 || err != nil {
		return "error", "HTTP request failed"
	}

	// Client errors (4xx)
	if status >= 400 {
		return "warn", "HTTP request completed with client error"
	}

	// Slow requests
	if duration > slowThreshold {
		return "warn", "Slow HTTP request detected"
	}

	// Success
	return "info", "HTTP request completed"
}

// shouldSkipPath checks if a path should be skipped from logging
func shouldSkipPath(path string, skipMap map[string]bool) bool {
	// Check exact match
	if skipMap[path] {
		return true
	}

	// Check prefix match for wildcard paths
	for skipPath := range skipMap {
		if strings.HasSuffix(skipPath, "*") {
			prefix := strings.TrimSuffix(skipPath, "*")
			if strings.HasPrefix(path, prefix) {
				return true
			}
		}
	}

	return false
}
