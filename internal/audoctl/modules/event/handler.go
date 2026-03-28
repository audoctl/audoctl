package event

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
	events := r.Group("/v1/events")

	events.Post("/", h.CreateEvent)
}

func (h Handler) CreateEvent(c fiber.Ctx) error {
	var req EventCreateRequest
	if err := c.Bind().Body(&req); err != nil {
		return errs.BadRequest("Invalid request body")
	}

	event, err := h.service.CreateEvent(c.Context(), req)
	if err != nil {
		return err
	}

	return response.Created(c, event)
}
