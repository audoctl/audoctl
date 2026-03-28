# Sessions Module – Audoctl

## Overview

The **Sessions module** provides context for events.
Each session represents a single user, AI agent, or process instance.
All events are tied to a session for proper timeline grouping and querying.

---

## Session Model

```go
type Session struct {
    TimestampedEntity
    ActorID   string    // user, bot, or system identifier
    Type      string    // type of session (user, ai, system)
}

### Example schema (Postgres)

```sql
CREATE TABLE sessions (
    id TEXT PRIMARY KEY,
    actor_id TEXT,
    type TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

