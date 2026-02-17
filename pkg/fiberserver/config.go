package fiberserver

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v3"
)

// Config contains Fiber server configuration
type Config struct {
	// Server basics
	Host string
	Port int

	// Timeouts
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration

	// Body limits
	BodyLimit      int
	MaxHeaderBytes int

	// Buffer sizes
	ReadBufferSize  int
	WriteBufferSize int

	// Performance
	Prefork                      bool
	Concurrency                  int
	ReduceMemoryUsage            bool
	DisableKeepalive             bool
	StreamRequestBody            bool
	DisablePreParseMultipartForm bool

	// Compression
	EnableCompression bool
	CompressLevel     int

	// Security
	EnableRequestID    bool
	RequestIDHeader    string
	TrustedProxies     []string
	ProxyHeader        string
	EnableIPValidation bool

	// Rate limiting
	EnableRateLimit   bool
	RateLimitMax      int
	RateLimitDuration time.Duration

	// Graceful shutdown
	ShutdownTimeout time.Duration

	// Development/Debug
	DisableStartupMessage bool
	EnablePrintRoutes     bool
	EnableStackTrace      bool

	// Observability
	EnableMetrics     bool
	MetricsPath       string
	EnableHealthCheck bool
	HealthCheckPath   string
	EnablePprof       bool
}

// ToFiberConfig converts Config to fiber.Config
func (c *Config) ToFiberConfig() fiber.Config {
	return fiber.Config{
		// Server
		ServerHeader:  "audoctl",
		StrictRouting: false,
		CaseSensitive: false,
		UnescapePath:  true,

		// Body limits
		BodyLimit: c.BodyLimit,

		// Timeouts
		ReadTimeout:  c.ReadTimeout,
		WriteTimeout: c.WriteTimeout,
		IdleTimeout:  c.IdleTimeout,

		// Buffer sizes
		ReadBufferSize:  c.ReadBufferSize,
		WriteBufferSize: c.WriteBufferSize,

		// Performance
		ReduceMemoryUsage:            c.ReduceMemoryUsage,
		DisableKeepalive:             c.DisableKeepalive,
		StreamRequestBody:            c.StreamRequestBody,
		DisablePreParseMultipartForm: c.DisablePreParseMultipartForm,

		// Security
		ProxyHeader:        c.ProxyHeader,
		EnableIPValidation: c.EnableIPValidation,
	}
}

// GetAddress returns the full server address
func (c *Config) GetAddress() string {
	if c.Host == "" {
		c.Host = "0.0.0.0"
	}
	if c.Port == 0 {
		c.Port = 8080
	}
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// DefaultConfig returns a config with sensible defaults
func DefaultConfig() Config {
	return Config{
		Host:                         "0.0.0.0",
		Port:                         9090,
		ReadTimeout:                  10 * time.Second,
		WriteTimeout:                 10 * time.Second,
		IdleTimeout:                  120 * time.Second,
		BodyLimit:                    10 * 1024 * 1024, // 10MB
		MaxHeaderBytes:               8192,             // 8KB
		ReadBufferSize:               8192,             // 8KB
		WriteBufferSize:              8192,             // 8KB
		Prefork:                      false,
		Concurrency:                  262144,
		ReduceMemoryUsage:            false,
		DisableKeepalive:             false,
		StreamRequestBody:            false,
		DisablePreParseMultipartForm: true,
		EnableCompression:            true,
		CompressLevel:                1,
		EnableRequestID:              true,
		RequestIDHeader:              "X-Request-ID",
		ProxyHeader:                  "X-Forwarded-For",
		EnableIPValidation:           false,
		EnableRateLimit:              false,
		RateLimitMax:                 100,
		RateLimitDuration:            time.Minute,
		ShutdownTimeout:              30 * time.Second,
		DisableStartupMessage:        false,
		EnablePrintRoutes:            false,
		EnableStackTrace:             false,
		EnableMetrics:                false,
		MetricsPath:                  "/metrics",
		EnableHealthCheck:            true,
		HealthCheckPath:              "/health",
		EnablePprof:                  false,
	}
}
