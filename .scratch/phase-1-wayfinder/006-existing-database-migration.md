---
status: closed
assignee: copilot
created: 2026-07-23
wayfinder:task
parent: 000-map
---

# Define the existing-database migration

## Question

How should an existing installation migrate its `receipts.receipt_number` column from `NOT NULL` to nullable so persisted drafts remain valid without losing data or requiring manual database recreation?

## Resolution

Use an explicit, numbered schema migration that detects the legacy `NOT NULL` definition via `PRAGMA table_info(receipts)`. When detected, migrate inside a startup-only transaction by creating the current receipts table with nullable `receipt_number`, copying every existing receipt row by column name, recreating the receipts indexes, and dropping the legacy table only after the copy succeeds. Preserve existing receipt numbers and statuses; do not rewrite finalized data or assign numbers to drafts.

Add a schema-migrations table so this repair is applied once and future schema changes are ordered and idempotent. The migration must fail loudly and leave the original database intact if the copy, constraint checks, or index recreation fails. Clean installs continue using the current `CREATE TABLE IF NOT EXISTS` schema, while existing installs are upgraded in place without manual database recreation.
