# Audoctl

> Audit and execution control for AI systems

A production-ready HTTP server built with Go and Fiber v3, designed for audit and execution control in AI systems.

## Features

- 🚀 **High Performance** - Built on Fiber v3 (powered by fasthttp)
- 🔒 **Security First** - Security headers, CORS, rate limiting, and more
- 📊 **Observability** - Structured logging, health checks, metrics, and profiling
- ⚙️ **Configurable** - YAML config files with environment variable overrides
- 🛡️ **Production Ready** - Graceful shutdown, compression, timeouts
- 🗄️ **Database Support** - PostgreSQL, MySQL, and SQLite
- 📝 **Structured Logging** - JSON logging with zerolog
- 🔄 **Graceful Shutdown** - Proper cleanup and shutdown handling

## Quick Start

### Prerequisites

- Go 1.25 or later
- Make (optional, for using Makefile commands)

### Installation

```bash
# Clone the repository
git clone https://github.com/audoctl/audoctl.git
cd audoctl

# Install dependencies
go mod download

# Build the application
make build
# or
go build -o bin/audoctl ./cmd
```

### Configuration

Create a `config.yaml` file in the project root (see `config.yaml` for example):

```yaml
database:
  driver: sqlite
  dsn: "audoctl.db"

application:
  name: Audoctl
  version: 1.0.0

log:
  level: info
  env: production

http_server:
  host: 0.0.0.0
  port: 8080
  enable_compression: true
  enable_request_id: true
```

Or use environment variables (see `.env.example`):

```bash
cp .env.example .env
# Edit .env with your settings
```

### Running

```bash
# Using Makefile
make run

# Or directly
./bin/audoctl audoctl

# With custom config
AUDOCTL_CONFIG=config.yaml ./bin/audoctl audoctl
```

## Usage

### Commands

```bash
# Start the server
audoctl audoctl

# Display version information
audoctl version

# Show help
audoctl --help
```

### API Endpoints

- `GET /ping` - Ping/liveness check
- `GET /health` - Health check endpoint
- `GET /api/*` - Your API routes (base path)
- `GET /metrics` - Prometheus metrics (if enabled)
- `GET /debug/pprof/*` - Profiling endpoints (if enabled)

## Configuration

The application supports configuration via:

1. **YAML file** (default: `config.yaml`)
2. **Environment variables** (override YAML values)

### Configuration Priority

Environment variables > YAML config > Default values

### Key Configuration Options

#### HTTP Server

- `http_server.port` - Server port (default: 8080)
- `http_server.enable_compression` - Enable gzip compression
- `http_server.enable_rate_limit` - Enable rate limiting
- `http_server.enable_request_id` - Add request ID to each request
- `http_server.enable_stack_trace` - Include stack traces in errors (dev only)

#### Security

- `security.enable_security_headers` - Add security headers
- `cors.enabled` - Enable CORS
- `cors.allow_origins` - Allowed origins

#### Database

- `database.driver` - Database driver (postgres, mysql, sqlite)
- `database.host` - Database host
- `database.port` - Database port
- `database.dsn` - Connection string (for SQLite)

#### Logging

- `log.level` - Log level (debug, info, warn, error)
- `log.env` - Environment (development, production)

See `config.yaml` for all available options.

## Development

### Prerequisites

- Go 1.25+
- Make
- golangci-lint (for linting)
- air (for hot reload, optional)

### Commands

```bash
# Format code
make fmt

# Run linter
make lint

# Run tests
make test

# Run tests with coverage
make test-coverage

# Run all checks
make check

# Hot reload development mode
make dev

# Clean build artifacts
make clean
```

### Project Structure

```
audoctl/
├── cmd/                    # Application entrypoints
│   ├── audoctl/           # Main application
│   │   ├── ctl.go         # CLI command
│   │   ├── server.go      # Server setup
│   │   └── starter.go     # Application starter
│   ├── version/           # Version command
│   ├── main.go            # Main entrypoint
│   └── root.go            # Root command
├── configs/               # Configuration
│   ├── config.go          # Config loader
│   ├── http.go            # HTTP config
│   ├── database.go        # Database config
│   └── ...
├── internal/              # Private application code
│   └── shared/
│       ├── helper/        # Helper utilities
│       └── tools/         # Shared tools
│           ├── gormdb/    # Database utilities
│           └── logger/    # Logger interface
├── pkg/                   # Public libraries
│   ├── fiberserver/       # Fiber server wrapper
│   │   ├── server.go      # Server implementation
│   │   ├── config.go      # Server config
│   │   ├── middleware.go  # Middleware implementations
│   │   └── server_option.go # Server options
│   ├── graceful/          # Graceful shutdown
│   └── zerologwrapper/    # Zerolog wrapper
├── config.yaml            # Configuration file
├── .env.example           # Example environment variables
├── Makefile               # Build automation
└── README.md              # This file
```

## Production Deployment

### Build for Production

```bash
# Build with version information
make VERSION=1.0.0 build

# The binary will be in bin/audoctl
```

### Docker Deployment

```dockerfile
# Example Dockerfile
FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o audoctl ./cmd

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/audoctl .
COPY config.yaml .
CMD ["./audoctl", "audoctl"]
```

### Systemd Service

```ini
[Unit]
Description=Audoctl Service
After=network.target

[Service]
Type=simple
User=audoctl
WorkingDirectory=/opt/audoctl
ExecStart=/opt/audoctl/bin/audoctl audoctl
Restart=on-failure
RestartSec=5s

[Install]
WantedBy=multi-user.target
```

## Monitoring

### Health Checks

```bash
# Liveness check
curl http://localhost:8080/ping

# Readiness check
curl http://localhost:8080/health
```

### Metrics

Enable metrics in config:

```yaml
http_server:
  enable_metrics: true
  metrics_path: /metrics
```

### Profiling

Enable pprof endpoints (development only):

```yaml
http_server:
  enable_pprof: true
```

Access profiling at `http://localhost:8080/debug/pprof/`

## Security

- Enable security headers in production
- Use TLS/HTTPS for external access
- Configure CORS properly for your domain
- Use environment variables for sensitive data
- Enable rate limiting to prevent abuse
- Review and configure security headers

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- [Fiber](https://github.com/gofiber/fiber) - Web framework
- [Zerolog](https://github.com/rs/zerolog) - Logging
- [GORM](https://gorm.io/) - ORM
- [Cobra](https://github.com/spf13/cobra) - CLI framework

## Support

For support, please open an issue in the GitHub repository.


