# ADR-0005 Implementation Decisions

**Date**: 2026-07-18  
**Status**: **Approved** — Decisions incorporated into ADR-0005  
**Superseded By**: ADR-0005 (updated 2026-07-18)

---

## Purpose

This ADR documents the implementation decisions for ADR-0005 (Dashboard Layout and Navigation). These decisions resolve unclear items and provide concrete implementation guidance.

**Note**: All decisions from this ADR have been incorporated into ADR-0005. This document is retained for historical reference.

---

## ✅ Clear Items (No Action Needed)

These sections are well-defined and ready for implementation:

1. **Navigation Structure** — 5 primary sidebar items with clear icons
2. **Dashboard Home Content** — Quick actions (3 cards), weekly metrics chart, system status cards
3. **Header Design** — Clinic name (editable), notification bell with badge count
4. **Authentication Model** — Direct permissions (no role abstraction), in-memory sessions
5. **Icon Vocabulary** — Heroicons SVG inline, styled with TailwindCSS
6. **Notification Categories** — System (persistent) vs User Operation (ephemeral)

---

## ✅ Decisions Made

All unclear items have been resolved and incorporated into ADR-0005:

1. **Chart.js Implementation** ✅
   - **Decision**: Load from jsDelivr CDN
   - **Initialization**: Inline `<script>` in base template
   - **HTMX Refresh**: Full chart partial swap on navigation

2. **Notification Center Dropdown** ✅
   - **Decision**: Pre-rendered content, click-to-open
   - **Toggle**: Vanilla JS or Alpine.js for open/close
   - **Mark as Read**: HTMX POST to `/notifications/{id}/read`

3. **Weekly Metrics Cache** ✅
   - **Decision**: Skip cache, use live aggregation query
   - **Index**: `CREATE INDEX idx_receipts_visit_date ON receipts(visit_date)`
   - **Revisit**: If dashboard load > 500ms

4. **Flash Messages** ✅
   - **Decision**: Gorilla Sessions with cookie store
   - **Lifecycle**: Auto-clear after first display
   - **Pattern**: Middleware injects into template context

5. **Cleanup Job** ✅
   - **Decision**: Daily background goroutine (24-hour ticker)
   - **Trigger**: On app startup + every 24 hours
   - **Query**: `DELETE FROM notifications WHERE expires_at < CURRENT_TIMESTAMP`

6. **Permission Middleware** ✅
   - **Decision**: Single middleware with route-level checks
   - **Pattern**: `RequirePermission(perm)` wrapper
   - **Error**: HTTP 403 "Forbidden"

---

## 📋 Historical Action Items (Completed)

These items were resolved and implemented in ADR-0005:

- [x] Decide on Chart.js loading strategy (CDN from jsDelivr)
- [x] Define notification dropdown HTML structure (pre-rendered, click-to-open)
- [x] Confirm weekly metrics query approach (live query, no cache)
- [x] Choose flash message library (Gorilla Sessions)
- [x] Define cleanup job pattern (daily goroutine)
- [x] Document permission middleware pattern (single middleware)

---

## 📝 Related ADRs

- **ADR-0004**: Architecture decision (Go + HTMX + SQLite)
- **ADR-0005**: Dashboard layout (updated with these decisions)
- **ADR-0006**: Export format decision (CSV/PDF/Excel in Go)

---

**Last Updated**: 2026-07-18  
**Next Review**: When implementing notification system or permissions