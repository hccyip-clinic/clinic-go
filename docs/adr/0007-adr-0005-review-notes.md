# ADR-0005 Review: Unclear Items & Open Questions

**Date**: 2026-07-18  
**Status**: Needs Review  
**Reviewer**: @development-team

---

## Summary

ADR-0005 (Dashboard Layout and Navigation) is **mostly complete** but has several unclear items that need resolution before implementation.

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

## ⚠️ Unclear Items (Needs Decision)

### 1. Chart.js Implementation Details

**Current State:**
```html
<canvas id="weeklyChart" 
        data-patients="{{.PatientCounts}}" 
        data-receipts="{{.ReceiptCounts}}" 
        data-totals="{{.DailyTotals}}">
</canvas>
```

**Unclear:**
- How is the Chart.js instance initialized? (inline `<script>` tag in template?)
- Is Chart.js vendored or loaded from CDN?
- What happens if CDN is unavailable (offline scenario)?
- How does HTMX refresh the chart data without full page reload?

**Decision Needed:**
- [ ] **CDN vs Vendored**: CDN for simplicity, vendored for reliability?
- [ ] **Initialization**: Inline script in base template, or separate `chart.js` file?
- [ ] **HTMX Refresh**: Does the `<canvas>` get re-rendered on navigation, or does HTMX update data attributes and trigger a JS event?

**Recommendation:**
- Use CDN for prototype (`<script src="https://cdn.jsdelivr.net/npm/chart.js"></script>`)
- Inline initialization script in `base.html` template
- HTMX swaps the entire chart partial on dashboard navigation (simpler than data attribute updates)

---

### 2. Notification Center Dropdown

**Current State:**
```html
<!-- Header partial -->
<div hx-get="/notifications/unread-count" 
     hx-trigger="every 30s" 
     hx-swap="outerHTML">
  <span class="badge">{{.UnreadCount}}</span>
</div>
```

**Unclear:**
- What does the notification dropdown HTML structure look like?
- How does the user open/close the dropdown? (click toggle, hover, HTMX trigger?)
- Is the dropdown content loaded via HTMX on click, or pre-rendered and hidden?
- How does "mark as read" work? (HTMX POST, then update badge count?)

**Decision Needed:**
- [ ] **Dropdown Trigger**: Click vs hover? (Click is more accessible)
- [ ] **Content Loading**: HTMX on-click load vs pre-rendered hidden div?
- [ ] **Mark as Read**: Does it dismiss the notification, or just mark read (kept in history)?

**Recommendation:**
- Click-to-open dropdown (Alpine.js `x-data="{ open: false }"` or vanilla JS)
- Pre-render dropdown content in header partial (simpler than HTMX load on click)
- "Mark as read" via HTMX POST to `/notifications/{id}/read`, returns updated badge count

---

### 3. Weekly Metrics Cache Strategy

**Current State:**
```sql
CREATE TABLE weekly_metrics_cache (
  date DATE PRIMARY KEY,
  patient_count INTEGER,
  receipt_count INTEGER,
  total_amount INTEGER -- in cents
);
```

**Unclear:**
- When is the cache refreshed? (On every receipt finalization? Cron job?)
- What happens if receipt is deleted/modified after cache update?
- Is the cache worth the complexity, or should we just run the aggregation query on dashboard load?

**Decision Needed:**
- [ ] **Cache vs Live Query**: Is caching necessary for this scale? (Single-user, <100 receipts/week)
- [ ] **If Cached**: Refresh strategy? (Trigger on receipt create/update/delete)
- [ ] **If Live Query**: Query optimization? (Index on `visit_date`, filter by current week)

**Recommendation:**
- **Skip the cache for Phase 1** — run live aggregation query on dashboard load
- Query should be fast enough for <1000 receipts: 
  ```sql
  SELECT 
    DATE(visit_date) as date,
    COUNT(DISTINCT patient_id) as patient_count,
    COUNT(*) as receipt_count,
    SUM(grand_total) as total_amount
  FROM receipts
  WHERE visit_date >= DATE('now', 'weekday 0', '-6 days')
    AND visit_date <= DATE('now')
  GROUP BY DATE(visit_date)
  ```
- Add index: `CREATE INDEX idx_receipts_visit_date ON receipts(visit_date)`
- Revisit caching if dashboard load time > 500ms

---

### 4. Flash Message Implementation

**Current State:**
```html
{{if .FlashSuccess}}
  <div class="flash flash-success">{{.FlashSuccess}}</div>
{{end}}
```

**Unclear:**
- How are flash messages stored? (Session cookie, server-side session store?)
- How are they cleared after display? (Auto-dismiss on next request?)
- What's the Go implementation pattern? (Chi middleware + session context?)

**Decision Needed:**
- [ ] **Storage**: Session-based (Gorilla Sessions) or cookie-based flash messages?
- [ ] **Lifecycle**: Auto-clear after first display, or manual dismissal?
- [ ] **Template Context**: How to pass flash data to all templates? (Middleware injects into context)

**Recommendation:**
- Use **Gorilla Sessions** with cookie store for flash messages
- Flash lifecycle:
  1. Handler sets `session.AddFlash("Receipt finalized!", "success")`
  2. Handler saves session, redirects
  3. Next request reads flash, displays, clears
- Middleware injects `.FlashSuccess` and `.FlashError` into all template contexts

---

### 5. Cleanup Job Scheduling

**Current State:**
```
Notification auto-cleanup:
- Trigger: First app startup on Monday
- Query: `DELETE FROM notifications WHERE expires_at < CURRENT_TIMESTAMP`
```

**Unclear:**
- How is the "first startup on Monday" detected? (Check day of week in `init()`?)
- What if the app is never restarted on Monday? (Cleanup never runs?)
- Should cleanup be a scheduled background job instead?

**Decision Needed:**
- [ ] **Trigger Mechanism**: Startup check vs scheduled goroutine?
- [ ] **Frequency**: Weekly (Monday) vs daily (check if 7 days since last cleanup)?

**Recommendation:**
- **Daily background goroutine** — simpler and more reliable:
  ```go
  func startCleanupJob(db *sql.DB) {
      go func() {
          ticker := time.NewTicker(24 * time.Hour)
          for range ticker.C {
              _, err := db.Exec("DELETE FROM notifications WHERE expires_at < CURRENT_TIMESTAMP")
              if err != nil {
                  log.Printf("Cleanup failed: %v", err)
              }
          }
      }()
  }
  ```
- Run on app startup with `time.Since(lastCleanup) > 24 hours` check

---

### 6. Permission Strings Hardcoding

**Current State:**
```go
const (
    PermReceiptsCreate   Permission = "receipts:create"
    PermReceiptsRead     Permission = "receipts:read"
    // ... 12 more permissions
)
```

**Unclear:**
- Where are these permission checks placed? (Route middleware, handler logic?)
- How to test permission checks? (Unit test examples needed)
- What happens if a permission check fails? (403 error page, redirect to dashboard?)

**Decision Needed:**
- [ ] **Middleware Pattern**: One middleware per permission, or single middleware with permission list?
- [ ] **Error Handling**: Custom 403 page, or generic error?
- [ ] **Testing**: Unit test for each permission, or integration test?

**Recommendation:**
- **Single middleware** with route-level permission requirements:
  ```go
  func RequirePermission(perm Permission) func(http.Handler) http.Handler {
      return func(next http.Handler) http.Handler {
          return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
              user := getUserFromSession(r)
              if !user.HasPermission(perm) {
                  http.Error(w, "Forbidden", http.StatusForbidden)
                  return
              }
              next.ServeHTTP(w, r)
          })
      }
  }
  
  // Usage:
  r.Use(RequirePermission(PermReceiptsCreate))
  ```
- Test with unit tests: `TestRequirePermission_Allowed`, `TestRequirePermission_Forbidden`

---

## 📋 Action Items

Before starting ADR-0005 implementation:

1. **[ ] Decide on Chart.js loading strategy** (CDN vs vendored)
2. **[ ] Define notification dropdown HTML structure** (mockup needed)
3. **[ ] Confirm weekly metrics query approach** (live query, no cache)
4. **[ ] Choose flash message library** (Gorilla Sessions recommended)
5. **[ ] Define cleanup job pattern** (daily goroutine)
6. **[ ] Document permission middleware pattern** (code example)

---

## 📝 Related ADRs

- **ADR-0004**: Architecture decision (Go + HTMX + SQLite)
- **ADR-0006**: Export format decision (CSV/PDF/Excel in Go)

---

**Last Updated**: 2026-07-18  
**Next Review**: After implementation decisions are made