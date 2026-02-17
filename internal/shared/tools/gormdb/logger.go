package gormdb

import (
	"context"
	"time"

	"github.com/audoctl/audoctl/configs"
	"github.com/audoctl/audoctl/internal/shared/helper"
	"github.com/audoctl/audoctl/internal/shared/tools/logger"
	gormlogger "gorm.io/gorm/logger"
)

type GormLogger struct {
	logger logger.Logger
}

func NewGormLogger(l logger.Logger) *GormLogger {
	return &GormLogger{logger: l}
}

func (l *GormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	// This method would ideally set the log level of the logger,
	// but since Zerolog doesn't support changing the log level on the fly easily,
	// we might need to recreate the logger if needed or ignore if not applicable.
	// This is a placeholder for how it could be implemented.
	return l
}

func buildBaseFields(ctx context.Context) map[string]any {
	return map[string]any{
		helper.CorrelationIdKey: helper.GetCorrelationId(ctx),
		"env":                   configs.GetEnv(),
		"appName":               configs.GetAppName(),
		"hostName":              configs.GetHostname(),
		"type":                  "[SQL]",
	}
}

func (l *GormLogger) Info(ctx context.Context, msg string, args ...interface{}) {
	fields := buildBaseFields(ctx)
	argsToMap(fields, args...)

	l.logger.Info(msg, fields)
}

func (l *GormLogger) Warn(ctx context.Context, msg string, args ...interface{}) {
	fields := buildBaseFields(ctx)
	argsToMap(fields, args...)
	l.logger.Warn(msg, fields)
}

func (l *GormLogger) Error(ctx context.Context, msg string, args ...interface{}) {
	fields := buildBaseFields(ctx)
	argsToMap(fields, args...)
	l.logger.Error(msg, fields)
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	sql, rows := fc()
	latency := time.Since(begin).String()
	if err != nil {
		fields := buildBaseFields(ctx)
		argsToMap(fields, "error", err, "sql", sql, "rowsAffected", rows, "latency", latency)
		l.logger.Error("SQL error", fields)
	} else {
		fields := buildBaseFields(ctx)
		argsToMap(fields, "sql", sql, "rowsAffected", rows, "latency", latency)
		l.logger.Info("SQL trace", fields)
	}
}

func argsToMap(fields map[string]any, args ...interface{}) {
	for i := 0; i < len(args)-1; i += 2 {
		key, ok := args[i].(string)
		if !ok {
			continue
		}
		fields[key] = args[i+1]
	}
}
