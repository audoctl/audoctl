// Package response provides standardized HTTP response handling for Fiber
package response

import (
	"github.com/audoctl/audoctl/pkg/errs"
	"github.com/gofiber/fiber/v3"
)

// Response represents a standardized API response structure
type Response struct {
	Success bool      `json:"success"`
	Data    any       `json:"data,omitempty"`
	Error   *errs.Err `json:"error,omitempty"`
	Meta    *Meta     `json:"meta,omitempty"`
}

// Meta contains response metadata (pagination, timing, etc.)
type Meta struct {
	Timestamp   string `json:"timestamp,omitempty"`
	RequestID   string `json:"request_id,omitempty"`
	Page        int    `json:"page,omitempty"`
	PerPage     int    `json:"per_page,omitempty"`
	Total       int64  `json:"total,omitempty"`
	ProcessTime string `json:"process_time,omitempty"`
}

// ResponseBuilder helps build standardized responses
type ResponseBuilder struct {
	ctx          fiber.Ctx
	data         any
	err          *errs.Err
	meta         *Meta
	includeDebug bool
}

// New creates a new response builder
func New(c fiber.Ctx) *ResponseBuilder {
	return &ResponseBuilder{
		ctx:          c,
		includeDebug: false,
	}
}

// WithDebug enables debug information in error responses
func (r *ResponseBuilder) WithDebug(debug bool) *ResponseBuilder {
	r.includeDebug = debug
	return r
}

// Data sets the response data
func (r *ResponseBuilder) Data(data any) *ResponseBuilder {
	r.data = data
	return r
}

// Error sets the error
func (r *ResponseBuilder) Error(e *errs.Err) *ResponseBuilder {
	r.err = e
	return r
}

// Meta sets response metadata
func (r *ResponseBuilder) Meta(meta *Meta) *ResponseBuilder {
	r.meta = meta
	return r
}

// JSON sends the JSON response
func (r *ResponseBuilder) JSON() error {
	resp := Response{
		Success: r.err == nil,
		Data:    r.data,
		Error:   r.err,
		Meta:    r.meta,
	}

	// Add request context to error if present
	if r.err != nil {
		if r.err.RequestID == "" {
			if reqID := r.ctx.Get("X-Request-ID"); reqID != "" {
				r.err.RequestID = reqID
			}
		}

		// Add debug info if enabled
		r.err.WithDebugInfo(r.includeDebug)

		return r.ctx.Status(r.err.Status).JSON(resp)
	}

	return r.ctx.JSON(resp)
}

// Convenience functions

// Success sends a successful response with data
func Success(c fiber.Ctx, data any) error {
	return New(c).Data(data).JSON()
}

// Created sends a 201 Created response
func Created(c fiber.Ctx, data any) error {
	return c.Status(fiber.StatusCreated).JSON(Response{
		Success: true,
		Data:    data,
	})
}

// NoContent sends a 204 No Content response
func NoContent(c fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}

// Error sends an error response
func Error(c fiber.Ctx, err *errs.Err, includeDebug bool) error {
	return New(c).WithDebug(includeDebug).Error(err).JSON()
}

// WithPagination creates a response with pagination metadata
func WithPagination(c fiber.Ctx, data any, page, perPage int, total int64) error {
	return New(c).
		Data(data).
		Meta(&Meta{
			Page:    page,
			PerPage: perPage,
			Total:   total,
		}).
		JSON()
}
