package audoctl

import (
	"context"
	"log"

	"github.com/audoctl/audoctl/configs"
	"github.com/audoctl/audoctl/internal/shared/helper"
	"github.com/audoctl/audoctl/internal/shared/tools/gormdb"
	"github.com/audoctl/audoctl/pkg/graceful"
	"github.com/audoctl/audoctl/pkg/zerologwrapper"
)

func StartCtl(ctx context.Context, cfg *configs.Config) graceful.StopFn {
	// Initialize logger wrapper
	loggerWrapper, err := zerologwrapper.NewJSONLogger(zerologwrapper.Config{Level: cfg.Log.Level})
	if err != nil {
		log.Fatalf("[ERROR] could not initialize logger %s", err.Error())
	}

	// Get underlying zerolog.Logger
	logger := loggerWrapper.GetZerolog()

	logger.Info().
		Str("app", cfg.Application.Name).
		Str("version", cfg.Application.Version).
		Str("env", cfg.Log.Env).
		Msg("Starting application")

	// Initialize database
	gormLogger := gormdb.NewGormLogger(loggerWrapper)

	db, err := gormdb.Connect(&cfg.Database, gormLogger)
	if err != nil {
		logger.Fatal().
			Err(err).
			Str("driver", cfg.Database.Driver).
			Str("host", cfg.Database.Host).
			Int("port", cfg.Database.Port).
			Msg("Database connection failed")
	}

	logger.Info().
		Str("driver", cfg.Database.Driver).
		Str("host", cfg.Database.Host).
		Int("port", cfg.Database.Port).
		Str("database", cfg.Database.Name).
		Msg("Database connected successfully")

	// Create fiber server
	remoteServer := newFiberServer(cfg, db, logger)

	// Print server info
	if !cfg.HTTPServer.DisableStartupMessage {
		printServerInfo(cfg)
	}

	// Start server in goroutine
	go func() {
		logger.Info().
			Str("host", cfg.HTTPServer.Host).
			Int("port", cfg.HTTPServer.Port).
			Str("base_path", "/api").
			Msg("Starting HTTP server")

		if e := remoteServer.Listen(); e != nil {
			logger.Fatal().
				Err(e).
				Str("correlation_id", helper.GetCorrelationId(ctx)).
				Msg("HTTP server listen error")
		}
	}()

	// Return graceful shutdown function
	return func(ctx context.Context) error {
		logger.Info().Msg("Shutting down application")

		var lastErr error

		// Shutdown HTTP server
		if e := remoteServer.Shutdown(ctx); e != nil {
			logger.Error().
				Err(e).
				Str("correlation_id", helper.GetCorrelationId(ctx)).
				Msg("HTTP server shutdown error")
			lastErr = e
		} else {
			logger.Info().Msg("HTTP server shutdown successfully")
		}

		// Disconnect database
		if e := db.Disconnect(); e != nil {
			logger.Error().
				Err(e).
				Str("correlation_id", helper.GetCorrelationId(ctx)).
				Msg("Database disconnect error")
			lastErr = e
		} else {
			logger.Info().Msg("Database disconnected successfully")
		}

		logger.Info().Msg("Application shutdown complete")
		return lastErr
	}
}
