# Fiber Server Package

Production-grade Fiber v3 server implementation for audoctl with comprehensive configuration support.

## Features

✅ **Comprehensive Configuration** - 50+ configurable parameters via YAML or ENV  
✅ **Production-Ready** - Timeouts, rate limiting, compression, security headers  
✅ **Observability** - Request ID, trace ID, session tracking  
✅ **Security** - CORS, API key auth, TLS, security headers  
✅ **Performance** - Prefork, concurrency control, buffer tuning  
✅ **Middleware Suite** - Recovery, logging, compression, rate limiting  
✅ **Graceful Shutdown** - Configurable shutdown timeout  
✅ **Health Checks** - Built-in health check and ping endpoints  

## Quick Start

### 1. Basic Server Setup

```go
package main

import (
    "context"
    "github.com/audoctl/audoctl/configs"
    "github.com/audoctl/audoctl/pkg/errs"
    "github.com/audoctl/audoctl/pkg/fiberserver"
)

func main() {
    // Load config
    cfg, _ := configs.LoadConfig(context.Background())
    
    // Create server
    server := fiberserver.New(
        fiberserver.Config{
            Host:         cfg.HTTPServer.Host,
            Port:         cfg.HTTPServer.Port,
            ReadTimeout:  cfg.HTTPServer.ReadTimeout,
            WriteTimeout: cfg.HTTPServer.WriteTimeout,
            BodyLimit:    cfg.HTTPServer.BodyLimit,
        },
        errs.NewFiberHandler(true).Middleware(),
        fiberserver.WithRecoverMiddleware(true),
        fiberserver.WithRequestIDMiddleware(),
        fiberserver.WithPing("/ping"),
    )
    
    // Start server
    server.Listen()
}
```

### 2. Configuration

Create `config.yaml`:

```yaml
http_server:
  host: 0.0.0.0
  port: 3000
  read_timeout: 30s
  write_timeout: 30s
  body_limit: 10485760  # 10MB
  enable_compression: true
  enable_request_id: true
  enable_rate_limit: true
  rate_limit_max: 100
  rate_limit_duration: 1m

cors:
  enabled: true
  allow_origins: "*"
  allow_methods: "GET,POST,PUT,DELETE,PATCH,OPTIONS"
  allow_credentials: false

security:
  enable_security_headers: true
  xss_protection: "1; mode=block"
  content_type_nosniff: "nosniff"
```

Or use environment variables:

```bash
export HTTP_PORT=3000
export HTTP_BODY_LIMIT=10485760
export HTTP_ENABLE_COMPRESSION=true
export CORS_ALLOW_ORIGINS="https://yourdomain.com"
```

### 3. Full Server Example

```go
func CreateServer(cfg *configs.Config) *fiberserver.Server {
    isDev := cfg.Log.Env == "development"
    
    server := fiberserver.New(
        fiberserver.Config{
            Host:                  cfg.HTTPServer.Host,
            Port:                  cfg.HTTPServer.Port,
            ReadTimeout:           cfg.HTTPServer.ReadTimeout,
            WriteTimeout:          cfg.HTTPServer.WriteTimeout,
            IdleTimeout:           cfg.HTTPServer.IdleTimeout,
            BodyLimit:             cfg.HTTPServer.BodyLimit,
            ReadBufferSize:        cfg.HTTPServer.ReadBufferSize,
            WriteBufferSize:       cfg.HTTPServer.WriteBufferSize,
            Prefork:               cfg.HTTPServer.Prefork,
            EnableCompression:     cfg.HTTPServer.EnableCompression,
            CompressLevel:         cfg.HTTPServer.CompressLevel,
            EnableRequestID:       cfg.HTTPServer.EnableRequestID,
            RequestIDHeader:       cfg.HTTPServer.RequestIDHeader,
            EnableRateLimit:       cfg.HTTPServer.EnableRateLimit,
            DisableStartupMessage: cfg.HTTPServer.DisableStartupMessage,
            EnableStackTrace:      isDev,
        },
        errs.NewFiberHandler(isDev).Middleware(),
        
        // Core middleware
        fiberserver.WithRecoverMiddleware(isDev),
        fiberserver.WithRequestIDMiddleware(),
        fiberserver.WithSessionIDMiddleware(),
        
        // CORS
        fiberserver.WithCORS(
            cfg.CORS.AllowOrigins,
            cfg.CORS.AllowMethods,
            cfg.CORS.AllowHeaders,
            cfg.CORS.AllowCredentials,
        ),
        
        // Security
        fiberserver.WithSecurityHeaders(fiberserver.SecurityHeadersConfig{
            XSSProtection:         cfg.Security.XSSProtection,
            ContentTypeNosniff:    cfg.Security.ContentTypeNosniff,
            XFrameOptions:         cfg.Security.XFrameOptions,
            HSTSMaxAge:            cfg.Security.HSTSMaxAge,
            HSTSIncludeSubdomains: cfg.Security.HSTSIncludeSubdomains,
            ContentSecurityPolicy: cfg.Security.ContentSecurityPolicy,
            ReferrerPolicy:        cfg.Security.ReferrerPolicy,
        }),
        
        // Compression
        fiberserver.WithCompression(),
        
        // Health checks
        fiberserver.WithPing("/ping"),
        
        // Base router
        fiberserver.WithBaseRouter("/api/v1"),
    )
    
    // Register routes
    server.RegisterHandlers(
        sessionHandler,
        eventHandler,
    )
    
    return server
}
```

## Configuration Parameters

### Server Basics

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `host` | string | `0.0.0.0` | Server bind address |
| `port` | int | `3000` | Server port |

### Timeouts

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `read_timeout` | duration | `10s` | Max time to read request |
| `write_timeout` | duration | `10s` | Max time to write response |
| `idle_timeout` | duration | `120s` | Keep-alive timeout |
| `shutdown_timeout` | duration | `30s` | Graceful shutdown timeout |

### Body Limits

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `body_limit` | int | `10485760` | Max request body size (bytes) |
| `max_header_bytes` | int | `8192` | Max header size (bytes) |

### Performance

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `prefork` | bool | `false` | Enable prefork mode (one process per CPU) |
| `concurrency` | int | `262144` | Max concurrent connections |
| `reduce_memory_usage` | bool | `false` | Reduce memory at cost of performance |
| `disable_keepalive` | bool | `false` | Disable HTTP keep-alive |

### Compression

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `enable_compression` | bool | `true` | Enable response compression |
| `compress_level` | int | `1` | Compression level (1=fast, 9=best) |

### Security

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `enable_request_id` | bool | `true` | Generate request IDs |
| `request_id_header` | string | `X-Request-ID` | Request ID header name |
| `enable_rate_limit` | bool | `false` | Enable rate limiting |
| `rate_limit_max` | int | `100` | Max requests per duration |
| `rate_limit_duration` | duration | `1m` | Rate limit window |

### Observability

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `enable_metrics` | bool | `false` | Enable Prometheus metrics |
| `metrics_path` | string | `/metrics` | Metrics endpoint path |
| `enable_health_check` | bool | `true` | Enable health checks |
| `health_check_path` | string | `/health` | Health check path |
| `enable_pprof` | bool | `false` | Enable pprof endpoints |

## Middleware

### Built-in Middleware

```go
// Recovery - catches panics
WithRecoverMiddleware(includeDebug bool)

// Request ID - generates unique IDs
WithRequestIDMiddleware(header string)

// Trace ID - distributed tracing
WithTraceIDMiddleware(header string)

// Session ID - extracts session ID
WithSessionIDMiddleware()

// Security Headers
WithSecurityHeaders(cfg SecurityHeadersConfig)

// Compression
WithCompression(level compress.Level)

// CORS
WithCORS(origins, methods, headers string, credentials bool)

// Rate Limiting
WithRateLimiter(max int, duration time.Duration)

// API Key Auth
WithAPIKeyAuth(apiKeys []string, header string)

// Timeout
WithTimeout(timeout time.Duration)

// Logging
WithLogging()
```

### Custom Middleware

```go
server := fiberserver.New(
    cfg,
    errorHandler,
    fiberserver.WithMiddleware(
        func(c fiber.Ctx) error {
            // Custom middleware logic
            return c.Next()
        },
    ),
)
```

## Health Checks

### Basic Health Check

```go
type DBHealthCheck struct {
    db *sql.DB
}

func (h *DBHealthCheck) Ping() error {
    return h.db.Ping()
}

server.RegisterOptions(
    fiberserver.WithHealthCheck(&DBHealthCheck{db}, "database"),
)
```

Access: `GET /health`

Response:
```json
{
  "success": true,
  "data": {
    "status": "healthy",
    "service": "database"
  }
}
```

## Error Handling

The server integrates with the `errs` package for structured error responses:

```go
func GetSession(c fiber.Ctx) error {
    sessionID := c.Params("id")
    
    session, err := repo.FindByID(sessionID)
    if err != nil {
        return errs.SessionNotFound(sessionID).
            WithOptions(errs.ExtractRequestContext(c)...)
    }
    
    return response.Success(c, session)
}
```

Error response:
```json
{
  "success": false,
  "error": {
    "id": "err_1708123456789",
    "key": "session.not_found",
    "category": "not_found",
    "message": "session not found",
    "status": 404,
    "request_id": "req_abc123",
    "session_id": "sess_123"
  }
}
```

## Production Deployment

### Docker Example

```dockerfile
FROM golang:1.26-alpine AS builder

WORKDIR /app
COPY . .
RUN go build -o audoctl ./cmd/audoctl

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/audoctl .
COPY --from=builder /app/config.yaml .

EXPOSE 3000
CMD ["./audoctl", "start"]
```

### Environment Variables

```bash
# Production
HTTP_PORT=8080
HTTP_PREFORK=true
HTTP_ENABLE_COMPRESSION=true
HTTP_ENABLE_RATE_LIMIT=true
HTTP_RATE_LIMIT_MAX=1000
HTTP_BODY_LIMIT=52428800  # 50MB
HTTP_READ_TIMEOUT=60s
HTTP_WRITE_TIMEOUT=60s
TLS_ENABLED=true
TLS_CERT_FILE=/etc/ssl/certs/audoctl.crt
TLS_KEY_FILE=/etc/ssl/private/audoctl.key
CORS_ALLOW_ORIGINS=https://audoctl.yourdomain.com
```

### Kubernetes Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: audoctl
spec:
  replicas: 3
  template:
    spec:
      containers:
      - name: audoctl
        image: audoctl:latest
        ports:
        - containerPort: 3000
        env:
        - name: HTTP_PORT
          value: "3000"
        - name: HTTP_PREFORK
          value: "true"
        - name: HTTP_ENABLE_RATE_LIMIT
          value: "true"
        livenessProbe:
          httpGet:
            path: /ping
            port: 3000
          initialDelaySeconds: 5
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health
            port: 3000
          initialDelaySeconds: 5
          periodSeconds: 5
```

## Performance Tuning

### For High Throughput

```yaml
http_server:
  prefork: true                 # One process per CPU
  concurrency: 524288           # 512K connections
  body_limit: 52428800          # 50MB for large payloads
  read_timeout: 60s
  write_timeout: 60s
  enable_compression: true
  compress_level: 1             # Fast compression
  reduce_memory_usage: false
```

### For Low Memory

```yaml
http_server:
  reduce_memory_usage: true
  concurrency: 65536            # 64K connections
  body_limit: 5242880           # 5MB
  disable_keepalive: true
  stream_request_body: true
```

### For Long-Running AI Operations

```yaml
http_server:
  read_timeout: 300s            # 5 minutes
  write_timeout: 300s           # 5 minutes
  body_limit: 104857600         # 100MB for large agent payloads
  stream_request_body: true
  disable_preparse_multipart: true
```

## License

Part of audoctl project - MIT License
