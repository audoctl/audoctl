package handler

import (
	"github.com/audoctl/audoctl/pkg/fiberserver/response"
	"github.com/gofiber/fiber/v3"
)

type PingHandler struct{}

func NewPingHandler() *PingHandler {
	return &PingHandler{}
}

func (ph *PingHandler) Serve(c fiber.Ctx) error {
	return response.Success(c, map[string]bool{"ok": true})
}
