## Event Model

Audoctl accepts flexible, schema-less events.

Each event represents a meaningful action in your system (e.g. user interaction, AI task, system event).

---

## Event Types

Audoctl does not enforce strict schemas, but we recommend using a structured naming convention:

### Naming convention

<domain>.<action>

### Examples

- user.login
- user.logout
- ui.click
- ui.view
- ai.prompt.sent
- ai.task.start
- ai.task.complete

---

## Storage

Audoctl stores events in an append-only structure optimized for timeline queries.

### Example schema (SQLite)

```sql
CREATE TABLE events (
  id TEXT PRIMARY KEY,
  session_id TEXT NOT NULL,
  type TEXT NOT NULL,
  source TEXT NOT NULL,
  actor_id TEXT,
  data TEXT,
  created_at DATETIME NOT NULL
);

CREATE INDEX idx_events_session_time
ON events(session_id, created_at);

CREATE INDEX idx_events_type
ON events(type);