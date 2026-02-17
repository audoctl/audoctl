# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial project setup with Go 1.26
- Production-ready HTTP server with Fiber v3
- Comprehensive configuration system (YAML + environment variables)
- Database support (PostgreSQL, MySQL, SQLite) with GORM
- Structured logging with zerolog
- Health check and ping endpoints
- Request ID and Trace ID middleware
- Security headers middleware
- CORS middleware with full configuration
- Rate limiting support
- Compression middleware (gzip, brotli)
- Graceful shutdown handling
- Profiling endpoints (pprof) for development
- Metrics endpoint (Prometheus-ready)
- CLI interface with Cobra
- Version command with build information
- Makefile for easy development
- Docker support (Dockerfile + docker-compose)
- Comprehensive README with documentation
- Example configuration files
- .gitignore and .dockerignore

### Features
- 🚀 High-performance HTTP server
- 🔒 Security-first design with multiple security layers
- 📊 Built-in observability (logging, health checks, metrics)
- ⚙️ Highly configurable via YAML or environment variables
- 🛡️ Production-ready with proper error handling
- 🗄️ Multi-database support
- 📝 Structured JSON logging
- 🔄 Graceful shutdown with context timeout

### Security
- Security headers (XSS, Frame Options, CSP, HSTS, etc.)
- CORS with configurable origins and methods
- Request rate limiting
- API key authentication support
- TLS/SSL support
- Trusted proxy configuration
- IP validation

### Developer Experience
- Hot reload support (via air)
- Comprehensive Makefile commands
- Docker support for easy deployment
- Clear project structure
- Extensive documentation
- Example configurations

## [1.0.0] - 2026-02-17

### Initial Release
- First production-ready release of Audoctl
- Full HTTP server implementation
- Complete configuration system
- Database integration
- Security features
- Observability tools
- Docker support
- Comprehensive documentation

[Unreleased]: https://github.com/audoctl/audoctl/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/audoctl/audoctl/releases/tag/v1.0.0
