# Audoctl

> GitHub Actions for AI Agents – Audit & execution control for AI systems

[![Go Reference](https://pkg.go.dev/badge/github.com/audoctl/audoctl.svg)](https://pkg.go.dev/github.com/audoctl/audoctl)
[![Build Status](https://github.com/audoctl/audoctl/actions/workflows/ci.yml/badge.svg)](https://github.com/audoctl/audoctl/actions)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Docker Pulls](https://img.shields.io/docker/pulls/audoctl/audoctl.svg)](https://hub.docker.com/r/audoctl/audoctl)
[![Stars](https://img.shields.io/github/stars/audoctl/audoctl)](https://github.com/audoctl/audoctl/stargazers)

A production-ready HTTP server built with **Go 1.26 + Fiber v3**, designed to **track, manage, and query events** from AI workflows in real-time.

---

## 🚀 Features

- High Performance – Fiber v3 powered by fasthttp
- Security First – Security headers, CORS, rate limiting
- Observability – JSON logging, health checks, metrics, profiling
- Configurable – YAML config files with environment variable overrides
- Database Support – SQLite (default), MySQL, PostgreSQL
- Production Ready – Graceful shutdown, compression, timeouts

---

## 🏁 Quick Start

```bash
# Clone repo
git clone https://github.com/audoctl/audoctl.git
cd audoctl

# Build
make build
# or
go build -o bin/audoctl ./cmd

# Run with default config
./bin/audoctl audoctl

# Or with custom config
AUDOCTL_CONFIG=config.development.yaml ./bin/audoctl audoctl