package fiberserver

import "github.com/gofiber/fiber/v3"

type HandlerGroup interface {
	RegisterRoutes(r fiber.Router)
}
