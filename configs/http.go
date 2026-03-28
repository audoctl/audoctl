package configs

import (
	"fmt"
	"strings"
	"time"

	"github.com/audoctl/audoctl/pkg/fiberserver"
)

// HTTPServer contains HTTP server configuration
type HTTPServer struct {
	// Server basics
	Host string `yaml:"host" env:"HOST,default=0.0.0.0"`
	Port int    `yaml:"port" env:"PORT,default=3000"`

	// Timeouts
	ReadTimeout  time.Duration `yaml:"read_timeout" env:"READ_TIMEOUT,default=10s"`
	WriteTimeout time.Duration `yaml:"write_timeout" env:"WRITE_TIMEOUT,default=10s"`
	IdleTimeout  time.Duration `yaml:"idle_timeout" env:"IDLE_TIMEOUT,default=120s"`

	// Body limits (for AI agent event payloads)
	BodyLimit      int `yaml:"body_limit" env:"BODY_LIMIT,default=10485760"`         // 10MB default
	MaxHeaderBytes int `yaml:"max_header_bytes" env:"MAX_HEADER_BYTES,default=8192"` // 8KB

	// Buffer sizes
	ReadBufferSize  int `yaml:"read_buffer_size" env:"READ_BUFFER_SIZE,default=8192"`   // 8KB
	WriteBufferSize int `yaml:"write_buffer_size" env:"WRITE_BUFFER_SIZE,default=8192"` // 8KB

	// Performance
	Prefork                      bool `yaml:"prefork" env:"PREFORK,default=false"`
	Concurrency                  int  `yaml:"concurrency" env:"CONCURRENCY,default=262144"` // 256K
	ReduceMemoryUsage            bool `yaml:"reduce_memory_usage" env:"REDUCE_MEMORY_USAGE,default=false"`
	DisableKeepalive             bool `yaml:"disable_keepalive" env:"DISABLE_KEEPALIVE,default=false"`
	StreamRequestBody            bool `yaml:"stream_request_body" env:"STREAM_REQUEST_BODY,default=false"`
	DisablePreParseMultipartForm bool `yaml:"disable_preparse_multipart" env:"DISABLE_PREPARSE_MULTIPART,default=true"`

	// Compression
	EnableCompression bool `yaml:"enable_compression" env:"ENABLE_COMPRESSION,default=true"`
	CompressLevel     int  `yaml:"compress_level" env:"COMPRESS_LEVEL,default=1"` // 1=fast, 9=best compression

	// Security
	EnableRequestID    bool   `yaml:"enable_request_id" env:"ENABLE_REQUEST_ID,default=true"`
	RequestIDHeader    string `yaml:"request_id_header" env:"REQUEST_ID_HEADER,default=X-Request-ID"`
	TrustedProxies     string `yaml:"trusted_proxies" env:"TRUSTED_PROXIES,default="`
	ProxyHeader        string `yaml:"proxy_header" env:"PROXY_HEADER,default=X-Forwarded-For"`
	EnableIPValidation bool   `yaml:"enable_ip_validation" env:"ENABLE_IP_VALIDATION,default=false"`

	// Rate limiting (per IP)
	EnableRateLimit   bool          `yaml:"enable_rate_limit" env:"ENABLE_RATE_LIMIT,default=false"`
	RateLimitMax      int           `yaml:"rate_limit_max" env:"RATE_LIMIT_MAX,default=100"`
	RateLimitDuration time.Duration `yaml:"rate_limit_duration" env:"RATE_LIMIT_DURATION,default=1m"`

	// Graceful shutdown
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout" env:"SHUTDOWN_TIMEOUT,default=30s"`

	// Development/Debug
	DisableStartupMessage bool `yaml:"disable_startup_message" env:"DISABLE_STARTUP_MESSAGE,default=false"`
	EnablePrintRoutes     bool `yaml:"enable_print_routes" env:"ENABLE_PRINT_ROUTES,default=false"`
	EnableStackTrace      bool `yaml:"enable_stack_trace" env:"ENABLE_STACK_TRACE,default=false"`

	// Observability
	EnableMetrics        bool   `yaml:"enable_metrics" env:"ENABLE_METRICS,default=false"`
	MetricsPath          string `yaml:"metrics_path" env:"METRICS_PATH,default=/metrics"`
	EnableHealthCheck    bool   `yaml:"enable_health_check" env:"ENABLE_HEALTH_CHECK,default=true"`
	HealthCheckPath      string `yaml:"health_check_path" env:"HEALTH_CHECK_PATH,default=/health"`
	EnablePprofEndpoints bool   `yaml:"enable_pprof" env:"ENABLE_PPROF,default=false"`
	EnableSwagger        bool   `yaml:"enable_swagger" env:"ENABLE_SWAGGER,default=false"`
}

// CORS contains CORS configuration
type CORS struct {
	Enabled          bool   `yaml:"enabled" env:"ENABLED,default=true"`
	AllowOrigins     string `yaml:"allow_origins" env:"ALLOW_ORIGINS,default=*"`
	AllowMethods     string `yaml:"allow_methods" env:"ALLOW_METHODS,default=GET,POST,PUT,DELETE,PATCH,OPTIONS"`
	AllowHeaders     string `yaml:"allow_headers" env:"ALLOW_HEADERS,default=Origin,Content-Type,Accept,Authorization,X-Request-ID,X-Trace-ID"`
	AllowCredentials bool   `yaml:"allow_credentials" env:"ALLOW_CREDENTIALS,default=false"`
	ExposeHeaders    string `yaml:"expose_headers" env:"EXPOSE_HEADERS,default=X-Request-ID,X-Trace-ID"`
	MaxAge           int    `yaml:"max_age" env:"MAX_AGE,default=3600"` // 1 hour
}

// TLS contains TLS/SSL configuration
type TLS struct {
	Enabled    bool   `yaml:"enabled" env:"ENABLED,default=false"`
	CertFile   string `yaml:"cert_file" env:"CERT_FILE,default="`
	KeyFile    string `yaml:"key_file" env:"KEY_FILE,default="`
	MinVersion string `yaml:"min_version" env:"MIN_VERSION,default=1.2"` // TLS 1.2
}

// Security contains security headers configuration
type Security struct {
	// Security headers
	EnableSecurityHeaders bool   `yaml:"enable_security_headers" env:"ENABLE_SECURITY_HEADERS,default=true"`
	XSSProtection         string `yaml:"xss_protection" env:"XSS_PROTECTION,default=1; mode=block"`
	ContentTypeNosniff    string `yaml:"content_type_nosniff" env:"CONTENT_TYPE_NOSNIFF,default=nosniff"`
	XFrameOptions         string `yaml:"x_frame_options" env:"X_FRAME_OPTIONS,default=SAMEORIGIN"`
	HSTSMaxAge            int    `yaml:"hsts_max_age" env:"HSTS_MAX_AGE,default=31536000"` // 1 year
	HSTSIncludeSubdomains bool   `yaml:"hsts_include_subdomains" env:"HSTS_INCLUDE_SUBDOMAINS,default=true"`
	ContentSecurityPolicy string `yaml:"content_security_policy" env:"CONTENT_SECURITY_POLICY,default=default-src 'self'"`
	ReferrerPolicy        string `yaml:"referrer_policy" env:"REFERRER_POLICY,default=strict-origin-when-cross-origin"`

	// API Key authentication (optional)
	EnableAPIKeyAuth bool   `yaml:"enable_api_key_auth" env:"ENABLE_API_KEY_AUTH,default=false"`
	APIKeyHeader     string `yaml:"api_key_header" env:"API_KEY_HEADER,default=X-API-Key"`
}

// GetAddress returns the full server address
func (h *HTTPServer) GetAddress() string {
	return fmt.Sprintf("%s:%d", h.Host, h.Port)
}

// IsProduction returns true if running in production mode
func (h *HTTPServer) IsProduction() bool {
	return !h.EnableStackTrace && !h.EnablePrintRoutes
}

// GetBodyLimitMB returns body limit in megabytes
func (h *HTTPServer) GetBodyLimitMB() int {
	return h.BodyLimit / 1024 / 1024
}

// ToFiberServerConfig converts HTTPServer config to fiberserver.Config
func (h *HTTPServer) ToFiberServerConfig() fiberserver.Config {
	// Parse trusted proxies
	var trustedProxies []string
	if h.TrustedProxies != "" {
		trustedProxies = strings.Split(h.TrustedProxies, ",")
		// Trim whitespace
		for i := range trustedProxies {
			trustedProxies[i] = strings.TrimSpace(trustedProxies[i])
		}
	}

	return fiberserver.Config{
		// Server basics
		Host: h.Host,
		Port: h.Port,

		// Timeouts
		ReadTimeout:  h.ReadTimeout,
		WriteTimeout: h.WriteTimeout,
		IdleTimeout:  h.IdleTimeout,

		// Body limits
		BodyLimit:      h.BodyLimit,
		MaxHeaderBytes: h.MaxHeaderBytes,

		// Buffer sizes
		ReadBufferSize:  h.ReadBufferSize,
		WriteBufferSize: h.WriteBufferSize,

		// Performance
		Prefork:                      h.Prefork,
		Concurrency:                  h.Concurrency,
		ReduceMemoryUsage:            h.ReduceMemoryUsage,
		DisableKeepalive:             h.DisableKeepalive,
		StreamRequestBody:            h.StreamRequestBody,
		DisablePreParseMultipartForm: h.DisablePreParseMultipartForm,

		// Compression
		EnableCompression: h.EnableCompression,
		CompressLevel:     h.CompressLevel,

		// Security
		EnableRequestID:    h.EnableRequestID,
		RequestIDHeader:    h.RequestIDHeader,
		TrustedProxies:     trustedProxies,
		ProxyHeader:        h.ProxyHeader,
		EnableIPValidation: h.EnableIPValidation,

		// Rate limiting
		EnableRateLimit:   h.EnableRateLimit,
		RateLimitMax:      h.RateLimitMax,
		RateLimitDuration: h.RateLimitDuration,

		// Graceful shutdown
		ShutdownTimeout: h.ShutdownTimeout,

		// Development/Debug
		DisableStartupMessage: h.DisableStartupMessage,
		EnablePrintRoutes:     h.EnablePrintRoutes,
		EnableStackTrace:      h.EnableStackTrace,

		// Observability
		EnableMetrics:     h.EnableMetrics,
		MetricsPath:       h.MetricsPath,
		EnableHealthCheck: h.EnableHealthCheck,
		HealthCheckPath:   h.HealthCheckPath,
		EnablePprof:       h.EnablePprofEndpoints,
	}
}
