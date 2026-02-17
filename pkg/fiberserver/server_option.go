package fiberserver

import (
	"strings"
	"time"

	"github.com/audoctl/audoctl/pkg/fiberserver/handler"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/basicauth"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/gofiber/fiber/v3/middleware/pprof"
)

type ServerOpt func(s *Server)

// Option is an alias for ServerOpt for better API consistency
type Option = ServerOpt

// WithOnShutdown sets a function to be called on server shutdown
func WithOnShutdown(onShutdown func()) ServerOpt {
	return func(s *Server) {
		s.onShutdown = onShutdown
	}
}

// WithBaseRouter sets the base route prefix
// Should be added before other routes
func WithBaseRouter(prefix string, middlewares ...any) ServerOpt {
	return func(s *Server) {
		s.router = s.app.Group(prefix, middlewares...)
	}
}

// WithPing adds a ping/liveness handler at the specified path
func WithPing(path ...string) ServerOpt {
	pingPath := "/ping"
	if len(path) > 0 && path[0] != "" {
		pingPath = path[0]
	}

	pingHandler := handler.NewPingHandler()
	return func(s *Server) {
		s.app.Get(pingPath, pingHandler.Serve)
	}
}

// WithHealthCheck adds a health check handler
func WithHealthCheck(checker handler.HealthCheck, name string, path ...string) ServerOpt {
	healthHandler := handler.NewHealthCheckHandler(checker, name)

	return func(s *Server) {
		healthPath := s.cfg.HealthCheckPath
		if healthPath == "" {
			healthPath = "/health"
		}
		if len(path) > 0 && path[0] != "" {
			healthPath = path[0]
		}

		s.app.Get(healthPath, healthHandler.Serve)
	}
}

// WithHandlerGroups adds handler groups to the server
func WithHandlerGroups(hs ...HandlerGroup) ServerOpt {
	return func(s *Server) {
		for _, h := range hs {
			h.RegisterRoutes(s.router)
		}
	}
}

// WithRecoverMiddleware adds panic recovery middleware
func WithRecoverMiddleware(includeDebug bool) ServerOpt {
	return func(s *Server) {
		s.app.Use(RecoverMiddleware(includeDebug))
	}
}

// WithRequestIDMiddleware adds request ID middleware
func WithRequestIDMiddleware(header ...string) ServerOpt {
	h := "X-Request-ID"
	if len(header) > 0 && header[0] != "" {
		h = header[0]
	}

	return func(s *Server) {
		s.app.Use(RequestIDMiddleware(h))
	}
}

// WithTraceIDMiddleware adds trace ID middleware for distributed tracing
func WithTraceIDMiddleware(header ...string) ServerOpt {
	h := "X-Trace-ID"
	if len(header) > 0 && header[0] != "" {
		h = header[0]
	}

	return func(s *Server) {
		s.app.Use(TraceIDMiddleware(h))
	}
}

// WithSessionIDMiddleware adds session ID extraction middleware
func WithSessionIDMiddleware() ServerOpt {
	return func(s *Server) {
		s.app.Use(SessionIDMiddleware())
	}
}

// WithSecurityHeaders adds security headers to responses
func WithSecurityHeaders(cfg SecurityHeadersConfig) ServerOpt {
	return func(s *Server) {
		s.app.Use(SecurityHeadersMiddleware(cfg))
	}
}

// WithSecurityHeadersMiddleware is an alias for WithSecurityHeaders
func WithSecurityHeadersMiddleware(cfg SecurityHeadersConfig) ServerOpt {
	return WithSecurityHeaders(cfg)
}

// WithCompression adds compression middleware
func WithCompression(level ...compress.Level) ServerOpt {
	compressionLevel := compress.LevelDefault
	if len(level) > 0 {
		compressionLevel = level[0]
	}

	return func(s *Server) {
		s.app.Use(compress.New(compress.Config{
			Level: compressionLevel,
		}))
	}
}

// WithCORS adds CORS middleware
func WithCORS(allowOrigins, allowMethods, allowHeaders string, allowCredentials bool) ServerOpt {
	return func(s *Server) {
		s.app.Use(cors.New(cors.Config{
			AllowOrigins:     strings.Split(allowOrigins, ","),
			AllowMethods:     strings.Split(allowMethods, ","),
			AllowHeaders:     strings.Split(allowHeaders, ","),
			AllowCredentials: allowCredentials,
		}))
	}
}

// WithCORSConfig adds CORS middleware with custom config
func WithCORSConfig(cfg cors.Config) ServerOpt {
	return func(s *Server) {
		s.app.Use(cors.New(cfg))
	}
}

// WithCorsMiddleware is an alias for WithCORSConfig
func WithCorsMiddleware(cfg cors.Config) ServerOpt {
	return WithCORSConfig(cfg)
}

// WithRateLimiter adds rate limiting middleware
func WithRateLimiter(max int, duration time.Duration, keyGenerator ...func(fiber.Ctx) string) ServerOpt {
	cfg := limiter.Config{
		Max:        max,
		Expiration: duration,
		KeyGenerator: func(c fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Rate limit exceeded",
			})
		},
	}

	if len(keyGenerator) > 0 {
		cfg.KeyGenerator = keyGenerator[0]
	}

	return func(s *Server) {
		s.app.Use(limiter.New(cfg))
	}
}

// WithAPIKeyAuth adds API key authentication middleware
func WithAPIKeyAuth(apiKeys []string, header ...string) ServerOpt {
	h := "X-API-Key"
	if len(header) > 0 && header[0] != "" {
		h = header[0]
	}

	return func(s *Server) {
		s.app.Use(APIKeyAuthMiddleware(apiKeys, h))
	}
}

// WithTimeout adds request timeout middleware
func WithTimeout(timeout time.Duration) ServerOpt {
	return func(s *Server) {
		s.app.Use(TimeoutMiddleware(timeout))
	}
}

// WithTimeoutMiddleware is an alias for WithTimeout
func WithTimeoutMiddleware(timeout time.Duration) ServerOpt {
	return WithTimeout(timeout)
}

// WithLogging adds basic logging middleware
func WithLogging() ServerOpt {
	return func(s *Server) {
		s.app.Use(LoggingMiddleware())
	}
}

// WithLanguageDetection adds language detection middleware
func WithLanguageDetection() ServerOpt {
	return func(s *Server) {
		s.app.Use(DetectLanguageMiddleware)
	}
}

// WithLanguageMiddleware is an alias for WithLanguageDetection
func WithLanguageMiddleware() ServerOpt {
	return WithLanguageDetection()
}

// WithPprof adds pprof endpoints for profiling (development only)
func WithPprof() ServerOpt {
	return func(s *Server) {
		s.app.Use(pprof.New())
	}
}

// WithBasicAuth adds HTTP basic authentication
func WithBasicAuth(username, password string) ServerOpt {
	return func(s *Server) {
		s.app.Use(basicauth.New(basicauth.Config{
			Users: map[string]string{
				username: password,
			},
		}))
	}
}

// WithMiddleware adds custom middleware functions
func WithMiddleware(middlewares ...fiber.Handler) ServerOpt {
	return func(s *Server) {
		for _, m := range middlewares {
			s.app.Use(m)
		}
	}
}

// WithStaticFiles serves static files from a directory
func WithStaticFiles(prefix, root string) ServerOpt {
	return func(s *Server) {
		s.app.Get(prefix+"/*", func(c fiber.Ctx) error {
			return c.SendFile(root + c.Params("*"))
		})
	}
}

// ApplyProductionMiddleware applies recommended production middlewares
func ApplyProductionMiddleware() ServerOpt {
	return func(s *Server) {
		// Request ID
		if s.cfg.EnableRequestID {
			s.app.Use(RequestIDMiddleware(s.cfg.RequestIDHeader))
		}

		// Compression
		if s.cfg.EnableCompression {
			level := compress.Level(s.cfg.CompressLevel)
			s.app.Use(compress.New(compress.Config{
				Level: level,
			}))
		}

		// Rate limiting
		if s.cfg.EnableRateLimit {
			s.app.Use(limiter.New(limiter.Config{
				Max:        s.cfg.RateLimitMax,
				Expiration: s.cfg.RateLimitDuration,
				KeyGenerator: func(c fiber.Ctx) string {
					return c.IP()
				},
			}))
		}
	}
}

// ApplyCommonMiddleware applies commonly used middlewares based on config
func ApplyCommonMiddleware(includeDebug bool, corsAllowOrigins string) ServerOpt {
	return func(s *Server) {
		// Recovery
		s.app.Use(RecoverMiddleware(includeDebug))

		// Request ID
		if s.cfg.EnableRequestID {
			s.app.Use(RequestIDMiddleware(s.cfg.RequestIDHeader))
		}

		// Session ID extraction
		s.app.Use(SessionIDMiddleware())

		// CORS
		if corsAllowOrigins != "" {
			s.app.Use(cors.New(cors.Config{
				AllowOrigins: strings.Split(corsAllowOrigins, ","),
				AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
				AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Request-ID", "X-Trace-ID"},
			}))
		}

		// Compression
		if s.cfg.EnableCompression {
			s.app.Use(compress.New(compress.Config{
				Level: compress.Level(s.cfg.CompressLevel),
			}))
		}
	}
}
