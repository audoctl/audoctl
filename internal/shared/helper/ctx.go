package helper

import (
	"context"

	"github.com/gofiber/fiber/v3"
)

var CorrelationIdKey = "correlation-id"

func GetCorrelationId(ctx context.Context) string {
	val := ctx.Value(CorrelationIdKey)
	if val == nil {
		return ""
	}
	return val.(string)
}

func GetIP(c fiber.Ctx) string {
	if ips := c.IPs(); len(ips) > 0 {
		return ips[0]
	}
	return c.IP()
}
