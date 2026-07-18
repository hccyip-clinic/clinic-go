# ADR-0004: Go + HTMX + SQLite for Self-Hosted Desktop App

**Date**: 2026-07-15  
**Status**: Proposed  
**Drivers**: @development-team

## Context

Building a simple CRUD application for clinic receipt management requires:
- Single-binary deployment (runs on desktop, opens browser to localhost)
- Simple CRUD workflows (receipts, patients, reports)
- Minimal client-side interactivity (form validation, dynamic line items)
- Server-side data storage (SQLite)
- Dashboard UI (TailAdmin-inspired)
- Export formats: CSV, Excel, PDF for receipts and reports

The team is more comfortable with server-rendered HTML and Go than client-side JavaScript frameworks.

**Clarification:** Despite the historical name "Hugo" in this ADR title, Hugo (the static site generator) is NOT used. The architecture is:
- **Go HTTP server** with `html/template` for server-rendered HTML
- **HTMX** for partial page updates
- **Embedded SQLite** for data storage

### Considered Options

**Option 1: Keep Vite + React**
- Pros: Rich interactivity, existing codebase, offline-capable
- Cons: Complex debugging, build step, React learning curve

**Option 2: Go + HTMX + SQLite**
- Pros: Simple deployment, server-rendered HTML, minimal JavaScript, familiar mental model, single binary
- Cons: Network round-trips for interactions, no offline support, HTMX coordination complexity

**Option 3: Python + FastAPI/Flask + HTMX**
- Pros: Python ecosystem (pandas, reportlab), rapid prototyping, excellent HTMX support, team familiarity
- Cons: Python runtime dependency, larger deployment footprint, virtual environment management

**Option 4: .NET Blazor Server**
- Pros: Rich interactivity, C# ecosystem, single deployment model
- Cons: .NET runtime required, heavier resource usage, learning curve

## Decision

**Adopt Go + HTMX + SQLite** for the clinic management prototype.

### Rationale

1. **Team skills alignment** — Go + server-rendered HTML matches team expertise better than React
2. **Deployment simplicity** — Single binary with embedded SQLite
3. **Appropriate complexity** — HTMX provides "just enough" interactivity for CRUD workflows
4. **Fast prototyping** — Minimal setup, direct mapping from domain model to database schema

### Architecture

```
┌─────────────────────────────────────────────────────────┐
│  Single Binary (Go + html/template)                    │
│  - HTTP server (localhost:PORT)                         │
│  - SQLite database (embedded)                           │
│  - Business logic (Go handlers)                         │
└─────────────────────────────────────────────────────────┘
           ↓ serves HTML + TailwindCSS
┌─────────────────────────────────────────────────────────┐
│  Browser (http://localhost:PORT)                        │
│  - Server-rendered HTML                                 │
│  - TailwindCSS styling                                  │
│  - HTMX for dynamic interactions                        │
└─────────────────────────────────────────────────────────┘
```

## Consequences

### Positive

**Development Experience**
- Simpler debugging (server-side logic, no React DevTools)
- Faster iteration (no build step for Go code)
- Clearer mental model (request → handler → response)

**Deployment**
- Single binary distribution
- No npm/node runtime required for end users
- Embedded SQLite (no separate database server)

**Maintainability**
- Server-side validation (single source of truth)
- Easier to onboard developers familiar with traditional web development
- Less JavaScript to maintain

### Negative

**User Experience**
- Network round-trips for every interaction (form validation, line item updates)
- No offline support (requires server running)
- Slower perceived performance on high-latency systems

**Development**
- HTMX coordination complexity (event triggers, swap targets)
- Template testing less straightforward than pure functions
- Less reusable logic (Go handlers vs pure functions)

**Implementation Effort**
- Full rewrite from scratch
- Domain logic must be implemented in Go
- Testing strategy must be built from ground up

## Compliance

This decision aligns with:
- **Simplicity over flexibility** — CRUD workflows don't require complex state
- **Team skills optimization** — Go + HTML > React for this team

## Implementation Notes

### Core Patterns

**Server-Rendered HTML:**
- All pages rendered on server
- HTMX for partial page updates
- No client-side state management

**SQLite Repository:**
- Direct SQL queries (no ORM needed)
- WAL mode for reliability
- Integer cents for all monetary values

**HTMX Interactions:**
- Form validation with fragment swaps
- Dynamic line items (add/remove)
- Real-time calculations via server requests

## Testing Strategy

**Unit Tests** (Go):
- Handler logic (validation, calculations)
- Repository queries (CRUD operations)
- Business logic (discount math, receipt numbering)

**Integration Tests**:
- HTMX form submissions
- Template rendering
- Database transactions

**Manual Testing**:
- Full CRUD workflows
- HTMX interactions (add/remove line items, validation)
- Print layouts

## When to Revisit

This decision should be revisited if:
- Application requires rich client-side interactivity (drag-and-drop, real-time collaboration)
- Offline support becomes a requirement
- Performance issues with HTMX round-trips (>200ms perceived latency)

## References

- **Prototype Specification**: `docs/prototype-spec.md`
- **Quick Start Guide**: `docs/QUICKSTART.md`
- **Decision Matrix**: `docs/hugo-htmx-sqlite.md` (DELETED - misleading title)
- **Domain Model**: `DOMAIN.md`

---

**Last updated**: 2026-07-15  
**Maintained by**: Development Team