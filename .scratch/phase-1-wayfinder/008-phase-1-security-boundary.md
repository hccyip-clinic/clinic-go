---
status: closed
assignee: copilot
created: 2026-07-23
wayfinder:grilling
parent: 000-map
---

# Set the Phase 1 security boundary

## Question

What authentication, session-cookie, authorization, and CSRF guarantees must be present before the localhost Phase 1 slice is considered shippable, given the current session secret and direct-permission model?

## Resolution

Phase 1 requires a single-user login even on localhost. Credentials are initialized through a first-run setup flow and stored only as a salted Argon2id password hash in the SQLite `settings` row; plaintext passwords are never persisted. A local filesystem-authorized reset command is required, and changing the password invalidates all in-memory sessions.

Authenticated sessions use an in-memory store with an `httpOnly`, `SameSite` cookie. Every state-changing request requires a server-validated per-session CSRF token, including HTMX requests. The existing direct-permission model remains the authorization seam: the single account receives the fixed full permission set, and middleware checks permissions on protected routes without adding role-management UI.

The server binds to `127.0.0.1` by default. Listening on another interface requires explicit configuration and is not the Phase 1 default.
