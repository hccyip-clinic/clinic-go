---
status: done
created: 2026-07-22
wayfinder:task
parent: 000-map
---

# Choose and migrate to a pure-Go SQLite driver

## Question

Which pure-Go SQLite driver and compatibility adjustments are required before Phase 1 CRUD work, while preserving the current schema, WAL mode, foreign keys, and repository API?

## Resolution

Use `modernc.org/sqlite` v1.54.0. The database package keeps its existing `*database/sql` wrapper and `New` API, changes the registered driver name from `sqlite3` to `sqlite`, and preserves WAL mode and foreign-key pragmas. The focused database tests and server build pass.
