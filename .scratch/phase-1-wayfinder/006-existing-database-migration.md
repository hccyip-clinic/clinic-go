---
status: needs-triage
created: 2026-07-23
wayfinder:task
parent: 000-map
---

# Define the existing-database migration

## Question

How should an existing installation migrate its `receipts.receipt_number` column from `NOT NULL` to nullable so persisted drafts remain valid without losing data or requiring manual database recreation?
