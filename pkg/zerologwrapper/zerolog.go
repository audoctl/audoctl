package zerologwrapper

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

type Logger struct {
	z     zerolog.Logger
	level zerolog.Level
}

// Config struct
type Config struct {
	Level       string
	ConsoleMode bool // true -> pretty console, false -> JSON
}

// NewJSONLogger creates a new Logger
func NewJSONLogger(cfg Config) (*Logger, error) {
	level := zerolog.InfoLevel
	if cfg.Level != "" {
		l, err := zerolog.ParseLevel(cfg.Level)
		if err != nil {
			return nil, fmt.Errorf("invalid log level: %w", err)
		}
		level = l
	}

	zerolog.TimeFieldFormat = time.RFC3339Nano

	var output io.Writer = os.Stdout
	if cfg.ConsoleMode {
		output = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: "15:04:05",
		}
	}

	z := zerolog.New(output).
		Level(level).
		With().
		Timestamp().
		Logger()

	return &Logger{z: z, level: level}, nil
}

// ----------------- Basic Logging -----------------

func (l *Logger) Info(msg string, fields map[string]any) {
	if fields != nil {
		l.z.Info().Fields(fields).Msg(msg)
	} else {
		l.z.Info().Msg(msg)
	}
}

func (l *Logger) Debug(msg string, fields map[string]any) {
	if fields != nil {
		l.z.Debug().Fields(fields).Msg(msg)
	} else {
		l.z.Debug().Msg(msg)
	}
}

func (l *Logger) Warn(msg string, fields map[string]any) {
	if fields != nil {
		l.z.Warn().Fields(fields).Msg(msg)
	} else {
		l.z.Warn().Msg(msg)
	}
}

func (l *Logger) Error(msg string, fields map[string]any) {
	if fields != nil {
		l.z.Error().Fields(fields).Msg(msg)
	} else {
		l.z.Error().Msg(msg)
	}
}

func (l *Logger) Fatal(msg string, fields map[string]any) {
	if fields != nil {
		l.z.Fatal().Fields(fields).Msg(msg)
	} else {
		l.z.Fatal().Msg(msg)
	}
	os.Exit(1)
}

// ----------------- Context Support -----------------

func (l *Logger) WithContext(ctx context.Context) zerolog.Logger {
	return l.z.With().Ctx(ctx).Logger()
}

// ----------------- Helpers -----------------

func (l *Logger) GetLevel() string {
	return l.level.String()
}

// GetZerolog returns the underlying zerolog.Logger
func (l *Logger) GetZerolog() zerolog.Logger {
	return l.z
}
