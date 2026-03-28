package audoctl

import (
	"fmt"
	"strings"
	"time"

	"github.com/audoctl/audoctl/configs"
	"github.com/audoctl/audoctl/internal/audoctl"
	"github.com/audoctl/audoctl/internal/shared/tools/gormdb"
	"github.com/audoctl/audoctl/pkg/fiberserver"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/gofiber/fiber/v3/middleware/pprof"
	"github.com/rs/zerolog"
)

// newFiberServer creates and configures a production-ready Fiber server
func newFiberServer(cfg *configs.Config, db *gormdb.GormDB, logger zerolog.Logger) *fiberserver.Server {
	// Convert config to fiberserver config
	serverConfig := cfg.HTTPServer.ToFiberServerConfig()

	// Build server options
	options := []fiberserver.Option{
		// Base router configuration
		fiberserver.WithBaseRouter("api"),
		fiberserver.WithPing(),
	}

	// Request ID middleware
	if cfg.HTTPServer.EnableRequestID {
		options = append(options, fiberserver.WithRequestIDMiddleware(cfg.HTTPServer.RequestIDHeader))
	}

	// Trace ID middleware for distributed tracing
	options = append(options, fiberserver.WithTraceIDMiddleware("X-Trace-ID"))

	// Language detection middleware
	options = append(options, fiberserver.WithLanguageMiddleware())

	// Structured logging middleware (production-ready)
	if cfg.Log.HTTP.Enabled {
		options = append(options, fiberserver.WithMiddleware(createStructuredLoggingMiddleware(logger, cfg)))
	}

	// Security headers middleware
	if cfg.Security.EnableSecurityHeaders {
		options = append(options, fiberserver.WithSecurityHeadersMiddleware(fiberserver.SecurityHeadersConfig{
			XSSProtection:         cfg.Security.XSSProtection,
			ContentTypeNosniff:    cfg.Security.ContentTypeNosniff,
			XFrameOptions:         cfg.Security.XFrameOptions,
			HSTSMaxAge:            cfg.Security.HSTSMaxAge,
			HSTSIncludeSubdomains: cfg.Security.HSTSIncludeSubdomains,
			ContentSecurityPolicy: cfg.Security.ContentSecurityPolicy,
			ReferrerPolicy:        cfg.Security.ReferrerPolicy,
		}))
	}

	// CORS middleware
	if cfg.CORS.Enabled {
		corsConfig := cors.Config{
			AllowOrigins:     strings.Split(cfg.CORS.AllowOrigins, ","),
			AllowMethods:     strings.Split(cfg.CORS.AllowMethods, ","),
			AllowHeaders:     strings.Split(cfg.CORS.AllowHeaders, ","),
			AllowCredentials: cfg.CORS.AllowCredentials,
			ExposeHeaders:    strings.Split(cfg.CORS.ExposeHeaders, ","),
			MaxAge:           cfg.CORS.MaxAge,
		}
		options = append(options, fiberserver.WithCorsMiddleware(corsConfig))
	}

	// Compression middleware
	if cfg.HTTPServer.EnableCompression {
		compressionLevel := compress.Level(cfg.HTTPServer.CompressLevel)
		options = append(options, fiberserver.WithMiddleware(compress.New(compress.Config{
			Level: compressionLevel,
		})))
	}

	// Rate limiting middleware
	if cfg.HTTPServer.EnableRateLimit {
		options = append(options, fiberserver.WithMiddleware(limiter.New(limiter.Config{
			Max:        cfg.HTTPServer.RateLimitMax,
			Expiration: cfg.HTTPServer.RateLimitDuration,
			KeyGenerator: func(c fiber.Ctx) string {
				return c.IP()
			},
			LimitReached: func(c fiber.Ctx) error {
				return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
					"error":       "Rate limit exceeded",
					"retry_after": cfg.HTTPServer.RateLimitDuration.Seconds(),
				})
			},
		})))
	}

	// Recover middleware (should be one of the first)
	options = append(options, fiberserver.WithRecoverMiddleware(cfg.HTTPServer.EnableStackTrace))

	// Timeout middleware
	options = append(options, fiberserver.WithTimeoutMiddleware(cfg.HTTPServer.ReadTimeout))

	// Health check endpoint
	if cfg.HTTPServer.EnableHealthCheck {
		options = append(options, fiberserver.WithMiddleware(func(c fiber.Ctx) error {
			if c.Path() == cfg.HTTPServer.HealthCheckPath {
				return c.Status(fiber.StatusOK).JSON(fiber.Map{
					"status":    "healthy",
					"timestamp": time.Now().Unix(),
					"service":   cfg.Application.Name,
					"version":   cfg.Application.Version,
				})
			}
			return c.Next()
		}))
	}

	// Metrics endpoint (Prometheus)
	if cfg.HTTPServer.EnableMetrics {
		options = append(options, fiberserver.WithMiddleware(func(c fiber.Ctx) error {
			if c.Path() == cfg.HTTPServer.MetricsPath {
				// TODO: Implement Prometheus metrics
				return c.Status(fiber.StatusOK).SendString("# Metrics endpoint\n")
			}
			return c.Next()
		}))
	}

	// Pprof endpoints for profiling (development only)
	if cfg.HTTPServer.EnablePprofEndpoints {
		options = append(options, fiberserver.WithMiddleware(func(c fiber.Ctx) error {
			if strings.HasPrefix(c.Path(), "/debug/pprof") {
				return pprof.New()(c)
			}
			return c.Next()
		}))
	}

	if cfg.HTTPServer.EnableSwagger {
		options = append(options, fiberserver.WithSwagger("./cmd/audoctl/docs/swagger.json", "/api"))
		options = append(options, fiberserver.WithOpenAPI("./cmd/audoctl/docs/swagger.json", "/api"))
	}

	// Create error handler
	errorHandler := createErrorHandler(logger, cfg.HTTPServer.EnableStackTrace)

	// Create server
	server := fiberserver.New(serverConfig, errorHandler, options...)

	// Register application handlers
	// TODO: Register your domain-specific handlers here
	// Example:
	server.RegisterHandlers(
		audoctl.InitHandlers(
			db.Db,
			cfg,
			logger,
		).Handlers()...)

	return server
}

// createErrorHandler creates a custom error handler for Fiber
func createErrorHandler(logger zerolog.Logger, includeStackTrace bool) fiber.ErrorHandler {
	return func(c fiber.Ctx, err error) error {
		// Default to 500
		code := fiber.StatusInternalServerError

		// Retrieve status code from Fiber error
		if e, ok := err.(*fiber.Error); ok {
			code = e.Code
		}

		// Get request context
		requestID := getLocalString(c, "requestid")
		traceID := getLocalString(c, "traceid")

		// Log the error
		logEvent := logger.Error().
			Err(err).
			Int("status", code).
			Str("method", c.Method()).
			Str("path", c.Path()).
			Str("ip", c.IP()).
			Str("request_id", requestID).
			Str("trace_id", traceID)

		if includeStackTrace {
			logEvent.Str("stack", fmt.Sprintf("%+v", err))
		}

		logEvent.Msg("Request error")

		// Build error response
		response := fiber.Map{
			"error": err.Error(),
			"code":  code,
		}

		// Add request ID to response
		if requestID != "" {
			response["request_id"] = requestID
		}

		// Add trace ID to response
		if traceID != "" {
			response["trace_id"] = traceID
		}

		// Add stack trace in debug mode
		if includeStackTrace {
			response["stack"] = fmt.Sprintf("%+v", err)
		}

		// Send response
		return c.Status(code).JSON(response)
	}
}

// createStructuredLoggingMiddleware creates a production-ready structured logging middleware
func createStructuredLoggingMiddleware(logger zerolog.Logger, cfg *configs.Config) fiber.Handler {
	// Create a logger adapter that implements the required interface
	loggerAdapter := &zerologAdapter{logger: logger}

	// Build skip paths from config
	skipPaths := cfg.Log.HTTP.SkipPaths
	if skipPaths == nil {
		skipPaths = []string{}
	}

	// Add default skip paths if health/metrics are enabled
	if cfg.HTTPServer.EnableHealthCheck {
		skipPaths = append(skipPaths, cfg.HTTPServer.HealthCheckPath)
	}
	if cfg.HTTPServer.EnableMetrics {
		skipPaths = append(skipPaths, cfg.HTTPServer.MetricsPath)
	}

	return fiberserver.StructuredLoggingMiddleware(fiberserver.StructuredLoggingConfig{
		Logger:          loggerAdapter,
		LogRequestBody:  cfg.Log.HTTP.LogRequestBody,
		LogResponseBody: cfg.Log.HTTP.LogResponseBody,
		SkipPaths:       skipPaths,
		SlowThreshold:   cfg.Log.HTTP.SlowThreshold,
		MaxBodyLogSize:  cfg.Log.HTTP.MaxBodyLogSize,
	})
}

// zerologAdapter adapts zerolog.Logger to the interface expected by StructuredLoggingMiddleware
type zerologAdapter struct {
	logger zerolog.Logger
}

func (z *zerologAdapter) Info(msg string, args map[string]any) {
	event := z.logger.Info()
	for key, value := range args {
		event = event.Interface(key, value)
	}
	event.Msg(msg)
}

func (z *zerologAdapter) Warn(msg string, args map[string]any) {
	event := z.logger.Warn()
	for key, value := range args {
		event = event.Interface(key, value)
	}
	event.Msg(msg)
}

func (z *zerologAdapter) Error(msg string, args map[string]any) {
	event := z.logger.Error()
	for key, value := range args {
		event = event.Interface(key, value)
	}
	event.Msg(msg)
}

func (z *zerologAdapter) Debug(msg string, args map[string]any) {
	event := z.logger.Debug()
	for key, value := range args {
		event = event.Interface(key, value)
	}
	event.Msg(msg)
}

// getLocalString safely retrieves a string value from fiber.Ctx locals
func getLocalString(c fiber.Ctx, key string) string {
	if val := c.Locals(key); val != nil {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

// printServerInfo prints server startup information
func printServerInfo(cfg *configs.Config) {
	fmt.Printf("\n")
	fmt.Printf("╔═══════════════════════════════════════════════════════════╗\n")
	fmt.Printf("║                   AUDOCTL SERVER                          ║\n")
	fmt.Printf("╠═══════════════════════════════════════════════════════════╣\n")
	fmt.Printf("║ Service:     %-44s ║\n", cfg.Application.Name)
	fmt.Printf("║ Version:     %-44s ║\n", cfg.Application.Version)
	fmt.Printf("║ Environment: %-44s ║\n", cfg.Log.Env)
	fmt.Printf("║ Address:     %-44s ║\n", fmt.Sprintf("http://%s:%d", cfg.HTTPServer.Host, cfg.HTTPServer.Port))
	fmt.Printf("║ API Base:    %-44s ║\n", "/api")
	fmt.Printf("╠═══════════════════════════════════════════════════════════╣\n")
	fmt.Printf("║ Health:      %-44s ║\n", cfg.HTTPServer.HealthCheckPath)
	if cfg.HTTPServer.EnableMetrics {
		fmt.Printf("║ Metrics:     %-44s ║\n", cfg.HTTPServer.MetricsPath)
	}
	if cfg.HTTPServer.EnablePprofEndpoints {
		fmt.Printf("║ Profiling:   %-44s ║\n", "/debug/pprof")
	}
	fmt.Printf("╠═══════════════════════════════════════════════════════════╣\n")
	fmt.Printf("║ Features:                                                 ║\n")
	fmt.Printf("║   • CORS:            %-36s ║\n", boolToStatus(cfg.CORS.Enabled))
	fmt.Printf("║   • Compression:     %-36s ║\n", boolToStatus(cfg.HTTPServer.EnableCompression))
	fmt.Printf("║   • Rate Limiting:   %-36s ║\n", boolToStatus(cfg.HTTPServer.EnableRateLimit))
	fmt.Printf("║   • Security Headers: %-35s ║\n", boolToStatus(cfg.Security.EnableSecurityHeaders))
	fmt.Printf("║   • Request ID:      %-36s ║\n", boolToStatus(cfg.HTTPServer.EnableRequestID))
	fmt.Printf("╚═══════════════════════════════════════════════════════════╝\n")
	fmt.Printf("\n")
}

func boolToStatus(b bool) string {
	if b {
		return "enabled"
	}
	return "disabled"
}
