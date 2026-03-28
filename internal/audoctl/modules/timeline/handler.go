package timeline

import (
	"github.com/audoctl/audoctl/pkg/errs"
	"github.com/audoctl/audoctl/pkg/fiberserver/response"
	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) Handler {
	return Handler{
		service: service,
	}
}

func (h Handler) RegisterRoutes(r fiber.Router) {
	events := r.Group("/v1/timeline")

	events.Get("/:session_id", h.GetEventsBySessionID)
}

func (h Handler) GetEventsBySessionID(c fiber.Ctx) error {
	var req GetEventsRequest
	if err := c.Bind().Query(&req); err != nil {
		return errs.BadRequest("Invalid request query")
	}
	if err := c.Bind().URI(&req); err != nil {
		return errs.BadRequest("Invalid request URI")
	}

	if err := req.Validate(); err != nil {
		return err
	}

	events, err := h.service.GetEventsBySessionID(c.Context(), req)
	if err != nil {
		return err
	}

	return response.Success(c, events)
}
