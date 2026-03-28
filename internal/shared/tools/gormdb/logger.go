package gormdb

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/audoctl/audoctl/configs"
	"github.com/audoctl/audoctl/internal/shared/helper"
	"github.com/audoctl/audoctl/internal/shared/tools/logger"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type GormLogger struct {
	logger        logger.Logger
	config        configs.DatabaseLog
	logLevel      gormlogger.LogLevel
	slowThreshold time.Duration
}

// NewGormLogger creates a new GORM logger with configuration
func NewGormLogger(l logger.Logger, cfg configs.DatabaseLog) *GormLogger {
	// Parse log level
	level := parseLogLevel(cfg.LogLevel)

	return &GormLogger{
		logger:        l,
		config:        cfg,
		logLevel:      level,
		slowThreshold: cfg.SlowThreshold,
	}
}

// parseLogLevel converts string log level to GORM log level
func parseLogLevel(level string) gormlogger.LogLevel {
	switch level {
	case "silent":
		return gormlogger.Silent
	case "error":
		return gormlogger.Error
	case "warn":
		return gormlogger.Warn
	case "info":
		return gormlogger.Info
	default:
		return gormlogger.Warn
	}
}

func (l *GormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newLogger := *l
	newLogger.logLevel = level
	return &newLogger
}

// buildBaseFields creates base structured logging fields
func buildBaseFields(ctx context.Context) map[string]any {
	return map[string]any{
		helper.CorrelationIdKey: helper.GetCorrelationId(ctx),
		"env":                   configs.GetEnv(),
		"appName":               configs.GetAppName(),
		"hostName":              configs.GetHostname(),
		"component":             "database",
		"type":                  "sql",
	}
}

func (l *GormLogger) Info(ctx context.Context, msg string, args ...interface{}) {
	if l.logLevel < gormlogger.Info {
		return
	}

	fields := buildBaseFields(ctx)
	if len(args) > 0 {
		fields["details"] = fmt.Sprintf(msg, args...)
	}

	l.logger.Info(msg, fields)
}

func (l *GormLogger) Warn(ctx context.Context, msg string, args ...interface{}) {
	if l.logLevel < gormlogger.Warn {
		return
	}

	fields := buildBaseFields(ctx)
	if len(args) > 0 {
		fields["details"] = fmt.Sprintf(msg, args...)
	}

	l.logger.Warn(msg, fields)
}

func (l *GormLogger) Error(ctx context.Context, msg string, args ...interface{}) {
	if l.logLevel < gormlogger.Error {
		return
	}

	fields := buildBaseFields(ctx)
	if len(args) > 0 {
		fields["details"] = fmt.Sprintf(msg, args...)
	}

	l.logger.Error(msg, fields)
}

// Trace logs SQL queries with detailed information
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	// Skip if logging is disabled
	if !l.config.Enabled {
		return
	}

	if l.logLevel <= gormlogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	// Build base fields
	fields := buildBaseFields(ctx)
	fields["latency"] = elapsed.Milliseconds()
	fields["latency_human"] = elapsed.String()
	fields["rows_affected"] = rows

	// Handle SQL query logging (with or without parameters)
	if l.config.ParameterizedQueries {
		fields["sql"] = sql
	} else {
		fields["sql"] = sql
	}

	// Handle errors
	if err != nil {
		// Skip ErrRecordNotFound if configured
		if l.config.IgnoreNotFound && errors.Is(err, gorm.ErrRecordNotFound) {
			return
		}

		if l.logLevel >= gormlogger.Error {
			fields["error"] = err.Error()
			fields["error_type"] = fmt.Sprintf("%T", err)

			// Log as error
			l.logger.Error("SQL query failed", fields)
		}
		return
	}

	// Check for slow queries
	if elapsed > l.slowThreshold {
		if l.logLevel >= gormlogger.Warn {
			fields["slow_query"] = true
			fields["threshold"] = l.slowThreshold.String()

			l.logger.Warn("Slow SQL query detected", fields)
		}
		return
	}

	// Log successful queries at info level
	if l.logLevel >= gormlogger.Info {
		l.logger.Info("SQL query executed", fields)
	}
}
