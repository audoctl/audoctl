package fiberserver

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v3"
)

// Server represents the HTTP server
type Server struct {
	router     fiber.Router
	app        *fiber.App
	onShutdown func()
	cfg        Config
}

// New initializes a new Server with the given config and options
func New(cfg Config, errorHandler fiber.ErrorHandler, opts ...ServerOpt) *Server {
	// Apply defaults if needed
	if cfg.Port == 0 {
		cfg.Port = 3000
	}
	if cfg.BodyLimit == 0 {
		cfg.BodyLimit = 10 * 1024 * 1024 // 10MB
	}
	if cfg.ReadTimeout == 0 {
		cfg.ReadTimeout = 10 * time.Second
	}
	if cfg.WriteTimeout == 0 {
		cfg.WriteTimeout = 10 * time.Second
	}
	if cfg.IdleTimeout == 0 {
		cfg.IdleTimeout = 120 * time.Second
	}

	// Create Fiber config
	fiberCfg := cfg.ToFiberConfig()

	// Set error handler
	if errorHandler != nil {
		fiberCfg.ErrorHandler = errorHandler
	}

	// Create Fiber app
	app := fiber.New(fiberCfg)

	// Create server
	srv := &Server{
		cfg:    cfg,
		app:    app,
		router: app.Group(""),
	}

	// Apply options
	for _, opt := range opts {
		opt(srv)
	}

	return srv
}

// App returns the underlying Fiber app
func (s *Server) App() *fiber.App {
	return s.app
}

// Router returns the main router
func (s *Server) Router() fiber.Router {
	return s.router
}

// Config returns the server configuration
func (s *Server) Config() Config {
	return s.cfg
}

// RegisterHandlers registers handler groups with the main router
func (s *Server) RegisterHandlers(handlerGroups ...HandlerGroup) {
	for _, h := range handlerGroups {
		h.RegisterRoutes(s.router)
	}
}

// Listen starts the HTTP server
func (s *Server) Listen() error {
	addr := s.cfg.GetAddress()

	if !s.cfg.DisableStartupMessage {
		fmt.Printf("🚀 audoctl server starting on %s\n", addr)
	}

	return s.app.Listen(addr)
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	if !s.cfg.DisableStartupMessage {
		fmt.Println("🛑 Shutting down server...")
	}

	// Create shutdown context with timeout
	shutdownCtx, cancel := context.WithTimeout(ctx, s.cfg.ShutdownTimeout)
	defer cancel()

	// Call custom shutdown handler if provided
	if s.onShutdown != nil {
		s.onShutdown()
	}

	// Shutdown Fiber app
	return s.app.ShutdownWithContext(shutdownCtx)
}

// OnShutdown registers a function to be called on shutdown
func (s *Server) OnShutdown(fn func()) {
	s.onShutdown = fn
}

// RegisterOptions applies additional server options
func (s *Server) RegisterOptions(opts ...ServerOpt) {
	for _, opt := range opts {
		opt(s)
	}
}
