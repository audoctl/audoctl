package errs

import (
	"errors"

	"github.com/gofiber/fiber/v3"
)

// FiberHandler handles error responses in Fiber framework
type FiberHandler struct {
	includeDebug bool
}

// NewFiberHandler creates a new Fiber error handler
func NewFiberHandler(includeDebug bool) *FiberHandler {
	return &FiberHandler{
		includeDebug: includeDebug,
	}
}

// Handle handles error and sends appropriate response
func (h *FiberHandler) Handle(c fiber.Ctx, err error) error {
	if err == nil {
		return nil
	}

	// Try to convert to *Err
	if e, ok := ToErr(err); ok {
		// Add request context if not already set
		if e.RequestID == "" {
			reqID := c.Get("X-Request-ID")
			if reqID == "" {
				if localReqID, ok := c.Locals("requestid").(string); ok {
					reqID = localReqID
				}
			}
			e.RequestID = reqID
		}

		// Add debug info if enabled
		e.WithDebugInfo(h.includeDebug)

		return c.Status(e.Status).JSON(e)
	}

	// Handle Fiber errors
	if fiberErr, ok := err.(*fiber.Error); ok {
		e := New(
			"fiber.error",
			fiberErr.Message,
			fiberErr.Code,
			CategoryInternal,
		).WithOptions(WithCause(err))

		return c.Status(fiberErr.Code).JSON(e)
	}

	// Handle unknown errors
	e := Internal("Internal server error", err)
	e.WithDebugInfo(h.includeDebug)

	return c.Status(e.Status).JSON(e)
}

// Middleware returns a Fiber error handler middleware
func (h *FiberHandler) Middleware() fiber.ErrorHandler {
	return func(c fiber.Ctx, err error) error {
		return h.Handle(c, err)
	}
}

// HandleSuccess sends a success response
func HandleSuccess(c fiber.Ctx, data any) error {
	return c.JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}

// HandleCreated sends a created response
func HandleCreated(c fiber.Ctx, data any) error {
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}

// HandleNoContent sends a no content response
func HandleNoContent(c fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}

// ExtractRequestContext extracts common request context for errors
func ExtractRequestContext(c fiber.Ctx) []ErrorOption {
	opts := []ErrorOption{}

	// Request ID
	reqID := c.Get("X-Request-ID")
	if reqID == "" {
		if localReqID, ok := c.Locals("requestid").(string); ok {
			reqID = localReqID
		}
	}
	if reqID != "" {
		opts = append(opts, WithRequestID(reqID))
	}

	// Trace ID
	if traceID := c.Get("X-Trace-ID"); traceID != "" {
		opts = append(opts, WithTraceID(traceID))
	}

	// Session ID (specific to audoctl)
	sessionID := c.Params("session_id")
	if sessionID == "" {
		sessionID = c.Query("session_id")
	}
	if sessionID != "" {
		opts = append(opts, WithSessionID(sessionID))
	}

	return opts
}

// RecoverMiddleware is a panic recovery middleware that converts panics to errors
func RecoverMiddleware() fiber.Handler {
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
					err = errors.New("unknown panic")
				}

				e := Internal("Panic recovered", err).
					WithOptions(ExtractRequestContext(c)...).
					WithDebugInfo(true)

				_ = c.Status(e.Status).JSON(e)
			}
		}()

		return c.Next()
	}
}
