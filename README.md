![CI](https://github.com/audoctl/audoctl/actions/workflows/ci.yml/badge.svg)
![Release](https://img.shields.io/github/v/release/audoctl/audoctl)
![License](https://img.shields.io/github/license/audoctl/audoctl)
![Go](https://img.shields.io/badge/go-1.26-blue)

# audoctl

> Control plane for AI agents.
**audoctl** is a self-hosted execution and observability engine for AI agents and LLM-powered systems. It captures every step an agent takes—prompts, LLM calls, tool executions, and errors—into an ordered, append-only timeline, enabling deterministic debugging, replay, and auditability.

---

## 🚀 Key Features

* **Execution Timeline** – Track every agent step in order.
* **Deterministic Replay** – Reproduce agent runs exactly as they happened.
* **Cost & Token Tracking** – Monitor model usage and estimate operational costs.
* **Audit & Compliance Ready** – Immutable event history for governance.
* **Self-Hosted** – Run locally, in Docker, or embedded in Go services.

---

## 📦 Project Structure

```txt
audoctl/
├── cmd/
│   └── audoctl/
│       └── main.go             # CLI + HTTP entrypoint
├── internal/
│   ├── module/                 # Core modules
│   │   ├── session/
│   │   │   ├── handler.go      # session-specific API endpoints (internal)
│   │   │   ├── service.go
│   │   │   ├── repository.go
│   │   │   └── dto.go          # request/response objects
│   │   ├── event/
│   │   │   ├── handler.go
│   │   │   ├── service.go
│   │   │   ├── repository.go
│   │   │   └── dto.go
│   │   └── storage/
│   │       ├── handler.go      # ops if needed
│   │       ├── service.go
│   │       ├── repository.go
│   │       └── dto.go
│   └── model/
│       ├── entity/
│       │   ├── session.go
│       │   └── event.go
│       └── enum/
│           └── event_type.go
├── pkg/
│   └── api/
│       └── http.go             # public API entrypoint, routes
├── go.mod
└── README.md
```

---

## 💡 Quick Start

### Run with Docker (Postgres)

```bash
docker-compose up -d
```

### Start API server

```bash
go run ./cmd/audoctl
```

---

## 🧪 Example: Instrumenting an AI Agent

```go
sess := audoctl.StartSession(ctx, audoctl.SessionConfig{
    Agent: "refund_agent",
})

defer sess.End("success")

sess.Event("prompt", map[string]any{
    "template": "refund_v3",
    "input":    userMessage,
})
```

---

## 📡 MVP HTTP API

* **Create Session**: `POST /v1/sessions`
* **Append Event**: `POST /v1/sessions/{id}/events`
* **Get Timeline**: `GET /v1/sessions/{id}/timeline`
* **Finish Session**: `POST /v1/sessions/{id}/finish`

---

## 🧭 Roadmap

* [ ] SQLite adapter
* [ ] CLI (`audoctl timeline <session>`)
* [ ] Cost & token tracking
* [ ] Deterministic replay engine
* [ ] Optional Web UI

---

## 🤝 Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for details. Fork, create a branch, and submit a PR.

---

## 📄 License

MIT License. See [LICENSE](LICENSE).

---

## 🧠 Why audoctl?

AI agents are autonomous systems, but debugging and auditing them is still primitive. audoctl gives developers the power to **trace, replay, and audit** agent execution, making AI workflows **deterministic, observable, and product
