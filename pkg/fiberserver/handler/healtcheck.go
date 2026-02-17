package handler

import (
	"github.com/audoctl/audoctl/pkg/errs"
	"github.com/audoctl/audoctl/pkg/fiberserver/response"
	"github.com/gofiber/fiber/v3"
)

// HealthCheck interface for health check dependencies
type HealthCheck interface {
	Ping() error
}

// HealthHandler handles health check requests
type HealthHandler struct {
	checker HealthCheck
	name    string
}

// HealthResponse represents a health check response
type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}

// NewHealthCheckHandler creates a new health check handler
func NewHealthCheckHandler(checker HealthCheck, name string) *HealthHandler {
	return &HealthHandler{
		checker: checker,
		name:    name,
	}
}

// Serve handles health check requests
func (h *HealthHandler) Serve(c fiber.Ctx) error {
	if err := h.checker.Ping(); err != nil {
		// Service unavailable error
		return errs.External(h.name, err).
			SetMetadata("health_check", true).
			WithOptions(errs.ExtractRequestContext(c)...)
	}

	return response.Success(c, &HealthResponse{
		Status:  "healthy",
		Service: h.name,
	})
}
