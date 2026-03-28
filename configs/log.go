package configs

import "time"

type Log struct {
	App   App    `env:"APP" default:"[NO-APP-NAME]"`
	Level string `env:"LEVEL" default:"info"`
	Env   string `env:"ENV" default:"dev"`

	// Database logging
	Database DatabaseLog `yaml:"database"`

	// HTTP logging
	HTTP HTTPLog `yaml:"http"`
}

// DatabaseLog contains database logging configuration
type DatabaseLog struct {
	Enabled           bool          `yaml:"enabled" env:"DATABASE_ENABLED,default=true"`
	LogLevel          string        `yaml:"log_level" env:"DATABASE_LOG_LEVEL,default=warn"`        // silent, error, warn, info
	SlowThreshold     time.Duration `yaml:"slow_threshold" env:"DATABASE_SLOW_THRESHOLD,default=1s"` // queries slower than this will be logged
	IgnoreNotFound    bool          `yaml:"ignore_not_found" env:"DATABASE_IGNORE_NOT_FOUND,default=true"`
	ParameterizedQueries bool       `yaml:"parameterized_queries" env:"DATABASE_PARAMETERIZED_QUERIES,default=false"` // log with ? placeholders instead of values
}

// HTTPLog contains HTTP request/response logging configuration
type HTTPLog struct {
	Enabled           bool     `yaml:"enabled" env:"HTTP_ENABLED,default=true"`
	LogLevel          string   `yaml:"log_level" env:"HTTP_LOG_LEVEL,default=info"` // debug, info, warn, error
	LogRequestBody    bool     `yaml:"log_request_body" env:"HTTP_LOG_REQUEST_BODY,default=false"`
	LogResponseBody   bool     `yaml:"log_response_body" env:"HTTP_LOG_RESPONSE_BODY,default=false"`
	SkipPaths         []string `yaml:"skip_paths" env:"HTTP_SKIP_PATHS"`                             // paths to skip logging (e.g., /health, /metrics)
	SlowThreshold     time.Duration `yaml:"slow_threshold" env:"HTTP_SLOW_THRESHOLD,default=5s"`     // requests slower than this will be logged as warning
	MaxBodyLogSize    int      `yaml:"max_body_log_size" env:"HTTP_MAX_BODY_LOG_SIZE,default=1024"` // max bytes to log for request/response body
}

type App string

const (
	Audoctl App = "[AUDOCTL]"
)

const (
	AudoctlAppName = "Audoctl"
)

func (t App) IsValid() bool {
	types := map[App]struct{}{
		Audoctl: {},
	}

	_, ok := types[t]
	return ok
}
