# ADR-0005: Dashboard Layout and Navigation Structure

**Date**: 2026-07-17  
**Status**: Proposed  
**Drivers**: @development-team

## Context

Building the base dashboard layout for the clinic management system. The application is currently single-user (Practitioner + Admin roles held by same person) but must be auth-ready for future multi-user support.

**Architecture**: Go + HTMX + SQLite — server-rendered HTML using Go's `html/template` package, no React/TypeScript, no Hugo static site generator.

### Requirements

1. **Quick actions** for primary workflows + **system status** cards
2. **Four navigation domains**: Dashboard, Receipts, Patients, Reports, Settings
3. **Minimum viable auth** with direct permissions (skip role abstraction for now)
4. **Header** with clinic name + notification bell
5. **Home landing page** with quick actions, status cards, and weekly metrics chart
6. **Server-side rendering** — all pages rendered by Go templates, HTMX for dynamic interactions

### Considered Options

**Option 1: Full TailAdmin Replica**
- Complex sidebar with collapsible sections, breadcrumbs, full notification center
- Pros: Professional appearance, familiar pattern
- Cons: Overkill for single-user, implementation complexity

**Option 2: Minimal Shell**
- Sidebar + header + content only, no dashboard home widgets
- Pros: Fast to implement, focused
- Cons: Missing at-a-glance information, less actionable

**Option 3: Hybrid (Chosen)**
- TailAdmin visual style, simplified structure
- Dashboard home with curated quick actions + status cards + metrics
- Auth-ready infrastructure without over-engineering
- Server-rendered HTML with HTMX swaps for dynamic content

## Decision

**Adopt Hybrid approach**: Professional dashboard UI with single-user simplicity, auth-ready infrastructure, server-rendered HTML.

### Architecture

```
┌─────────────────────────────────────────────────────────┐
│  Header                                                  │
│  ┌─────────────────────────────┐ ┌──────────────────┐  │
│  │  Clinic Name (editable)     │ │ 🔔 Notifications │  │
│  └─────────────────────────────┘ └──────────────────┘  │
├─────────────────────────────────────────────────────────┤
│ Sidebar │  Content Area                                 │
│ ┌─────┐ │  ┌─────────────────────────────────────────┐ │
│ │📊   │ │  │  Quick Actions (Hero Section)           │ │
│ │📋   │ │  │  [Create Receipt] [Search Patient]      │ │
│ │👥   │ │  │  [Today's Receipts]                     │ │
│ │📈   │ │  ├─────────────────────────────────────────┤ │
│ │⚙️   │ │  │  Weekly Metrics (Combo Chart)           │ │
│ │     │ │  │  Mon-Sat: Patients + Receipts + Total   │ │
│ └─────┘ │  ├─────────────────────────────────────────┤ │
│         │  │  System Status Cards                    │ │
│         │  │  ┌───────────┐ ┌─────────────────────┐  │ │
│         │  │  │ Backup    │ │ Storage + Export    │  │ │
│         │  │  │ Status    │ │ Last Export Date    │  │ │
│         │  │  └───────────┘ └─────────────────────┘  │ │
│         │  └─────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────┘
```

### Navigation Structure

**Primary Sidebar (5 items)**:
1. 📊 **Dashboard** — Home landing page (quick actions + status + metrics)
2. 📋 **Receipts** — Create, finalize, print, archive receipts
3. 👥 **Patients** — Search, view history, manage patient records
4. 📈 **Reports** — Financial year statements, CSV exports, HK Inland Revenue
5. ⚙️ **Settings** — Clinic info, backup management, data retention, user config

**Icon Vocabulary** (Heroicons SVG inline):
- Dashboard: `home` (Heroicons outline)
- Receipts: `document-text` (Heroicons outline)
- Patients: `users` (Heroicons outline)
- Reports: `chart-bar` (Heroicons outline)
- Settings: `cog` (Heroicons outline)

**Implementation**: Inline SVG from Heroicons (no npm dependency), styled with TailwindCSS. Go templates render the sidebar with active state highlighting.

### Dashboard Home Content

**Quick Actions (3 cards)**:
1. **Create New Receipt** — Primary workflow entry point (`/receipts/new`)
2. **Search Patient** — Secondary workflow (`/patients/search`)
3. **View Today's Receipts** — Quick lookup (`/receipts?date=today`)

**Rationale**: 3 actions fits cognitive load best. Receipts/Patients are primary workflows. "Today's Receipts" provides quick access without full navigation.

**Excluded from quick actions** (use sidebar instead):
- Export Backup → Settings > Backup
- Open Settings → Sidebar Settings item

---

**Weekly Metrics Chart (Combo Chart)**:
- **X-axis**: Monday–Saturday (current week so far, Mon→Today)
- **Primary Y-axis (bars)**: Patient count, Receipt count
- **Secondary Y-axis (line)**: Daily total revenue (HKD)
- **Data definitions**:
  - Patients = unique patients with receipts that day
  - Receipts = total receipt count (including repeat visits)
  - Day Total = sum of `total_amount_hkd` from receipts
- **Implementation**: Chart.js via CDN, server-side data preparation, `<canvas>` element with HTMX refresh on navigation

**Rationale**: Combo chart shows volume (bars) + revenue trend (line) without visual clutter. Chart.js is simpler than D3, good enough for this use case.

---

**System Status Cards (2 cards)**:

**Backup Status Card**:
- States:
  - ✅ Success: "Last backup: Today, 12:00 AM"
  - ⚠️ Warning: "Last backup: Yesterday" (1 day late)
  - ❌ Critical: "Last backup: 3 days ago" (action required)
  - 🔧 Not configured: "Backup not set up" (initial state)
- Click behavior: Navigate to Settings > Backup for details + retry
- Shows: Last success/failure timestamp, error details on failure

**Storage Card**:
- Shows:
  - Database size (e.g., "2.4 MB")
  - Record counts: "2,347 receipts / 891 patients"
  - Growth rate: "+12 receipts this week"
  - Last export date: "Last export: 2026-07-15"
  - Percentage of practical limit (threshold: 1000 MB = 100%)
- Click behavior: Navigate to Settings > Data Management for archive/export options
- Thresholds:
  - < 50% (< 500MB): ✅ No action
  - 50-80% (500-800MB): ⚠️ Info only
  - 80-100% (800-1000MB): ⚠️ "Approaching limit"
  - > 100% (> 1000MB): ❌ "Consider archiving"

**Rationale**: SQLite has no hard limit, but 1GB is practical threshold for performance. Three-year retention should keep database well under this limit.

---

### Header Design

**Clinic Name**:
- Source: Editable field in Settings > Clinic Information (stored in `clinic_settings` table)
- Default: "Hong Ching Clinic" (from `config.json` `clinic_name_en` field)
- Display: Left side of header, prominent
- Click behavior: None (informational only)

**Notification Bell**:
- Badge count: Unread critical + failed notifications (today only)
- Click behavior: Open notification center dropdown (HTMX-loaded partial)
- Notification categories:
  - **System** (persistent, stored in DB, 7-day retention):
    - Critical: Backup failed 3+ days, storage > 1000MB
    - Warning: Backup 1 day late, storage > 80%
    - Info: Backup completed, new version available
  - **User Operation** (ephemeral, flash message):
    - Success: Receipt finalized, backup exported
    - Fail: Export failed, validation error

**Notification Center** (dropdown panel):
- Grouped by category (System: Critical/Warning/Info)
- Shows timestamp, message, action button (if applicable)
- Mark as read/unread functionality (HTMX POST to `/notifications/{id}/read`)
- "View all" link to full notifications page (future enhancement)

**Rationale**: Separating system notifications (persistent, actionable) from operation feedback (ephemeral flash messages) reduces noise. Bell badge shows only urgent items.

---

### Authentication Model

**Direct Permissions (Option C)**:
```go
type Permission string

const (
    PermReceiptsCreate   Permission = "receipts:create"
    PermReceiptsRead     Permission = "receipts:read"
    PermReceiptsUpdate   Permission = "receipts:update"
    PermReceiptsFinalize Permission = "receipts:finalize"
    PermReceiptsArchive  Permission = "receipts:archive"
    PermPatientsRead     Permission = "patients:read"
    PermPatientsCreate   Permission = "patients:create"
    PermPatientsUpdate   Permission = "patients:update"
    PermReportsGenerate  Permission = "reports:generate"
    PermReportsExport    Permission = "reports:export"
    PermSettingsRead     Permission = "settings:read"
    PermSettingsUpdate   Permission = "settings:update"
    PermBackupManage     Permission = "backup:manage"
    PermNotificationsRead Permission = "notifications:read"
)

type User struct {
    ID          string
    Username    string
    Permissions []Permission
}
```

**Default User** (single-user phase):
```go
defaultUser := User{
    ID:       "user-default",
    Username: "practitioner",
    Permissions: []Permission{
        PermReceiptsCreate, PermReceiptsRead, PermReceiptsUpdate,
        PermReceiptsFinalize, PermReceiptsArchive,
        PermPatientsRead, PermPatientsCreate, PermPatientsUpdate,
        PermReportsGenerate, PermReportsExport,
        PermSettingsRead, PermSettingsUpdate,
        PermBackupManage, PermNotificationsRead,
    },
}
```

**Auth Infrastructure**:
- Session cookie (httpOnly, secure, sameSite=strict)
- In-memory session store with cookie expiry (single-user, restart = logout)
- Middleware checks `HasPermission("receipts:create")` on routes
- User object includes permissions array (future: role-based)
- No UI differences (single user sees everything)

**Rationale**: Direct permissions avoid role abstraction complexity while remaining auth-ready. In-memory sessions are sufficient for single-user desktop app. Adding roles later requires only adding `roles` table and permission inheritance logic.

---

## Consequences

### Positive

**User Experience**:
- At-a-glance system status (backup, storage) drives action
- Quick actions reduce navigation clicks for primary workflows
- Weekly metrics provide business insight without full report navigation

**Technical**:
- Auth-ready without over-engineering (direct permissions, session cookie)
- Notification system supports future multi-user scenarios
- Modular layout (sidebar/header/content) supports easy page addition
- Server-rendered HTML = simpler debugging, no client-side state management

**Maintainability**:
- Clear navigation hierarchy (5 primary domains)
- Notification categorization prevents alert fatigue
- Storage thresholds are explicit and adjustable
- Go templates = type-safe rendering (compile-time checks)

### Negative

**Implementation Complexity**:
- Chart.js integration requires CDN + careful script loading
- Notification system needs database table + cleanup job (7-day retention)
- Auth middleware on every route (even if always passes in single-user phase)
- HTMX coordination for dynamic content (event triggers, swap targets)

**Performance**:
- Weekly metrics query on dashboard load (should use SQLite materialized view or cache)
- Notification count query on every page load (needs indexing on `is_read`, `expires_at`)

**Design Debt**:
- Hardcoded permission strings in Go (will need refactoring for RBAC)
- Notification categories may need expansion (e.g., "Error" separate from "Fail")
- Chart.js CDN dependency (could vendor if needed)

## Implementation Notes

### File Structure (Go + TailwindCSS v4)

```
├── templates/
│   ├── layouts/
│   │   ├── base.html            # Base template with <head>, sidebar, header
│   │   └── dashboard.html       # Dashboard home page
│   ├── partials/
│   │   ├── sidebar.html         # Navigation with Heroicons SVG
│   │   ├── header.html          # Clinic name + notification bell
│   │   ├── notification-center.html  # Dropdown panel (HTMX partial)
│   │   └── weekly-chart.html    # Chart.js canvas + data attributes
│   └── dashboard/
│       ├── quick-actions.html   # 3 action cards
│       ├── backup-status.html   # Backup status widget
│       └── storage-card.html    # Storage metrics widget
├── pages/
│   ├── receipts/
│   │   ├── list.html            # Receipt list with filters
│   │   ├── new.html             # Create receipt form
│   │   └── view.html            # Single receipt view
│   ├── patients/
│   │   ├── search.html          # Patient search + results
│   │   └── view.html            # Patient detail + history
│   ├── reports/
│   │   └── financial.html       # Financial year reports
│   └── settings/
│       ├── clinic.html          # Clinic info form
│       ├── backup.html          # Backup management
│       └── data.html            # Data retention + export
├── static/
│   ├── js/
│   │   └── chart.min.js         # Chart.js (vendored or CDN)
│   └── css/
│       └── tailwind.css         # Built TailwindCSS
├── handlers/
│   ├── auth.go                  # Session management, permission checks
│   ├── dashboard.go             # Dashboard page handler
│   ├── notifications.go         # Notification CRUD
│   └── permissions.go           # Permission constants + middleware
└── config.json                  # Clinic name, defaults
```

### Database Schema Additions

```sql
-- Notifications table (7-day retention)
CREATE TABLE notifications (
  id TEXT PRIMARY KEY,
  category TEXT NOT NULL, -- 'critical' | 'warning' | 'info' | 'success' | 'fail'
  scope TEXT NOT NULL,    -- 'system' | 'user_operation'
  title TEXT NOT NULL,
  message TEXT NOT NULL,
  action_url TEXT,        -- Optional: navigate on click
  is_read INTEGER DEFAULT 0,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  expires_at DATETIME     -- Auto-cleanup after 7 days
);

-- Clinic settings (single-row table)
CREATE TABLE clinic_settings (
  id INTEGER PRIMARY KEY CHECK (id = 1), -- Enforce single row
  clinic_name_en TEXT NOT NULL,
  clinic_name_zh TEXT,
  clinic_reg_no TEXT,
  clinic_address TEXT,
  clinic_telephone TEXT,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Weekly metrics materialized view (refresh on receipt finalization)
CREATE TABLE weekly_metrics_cache (
  date DATE PRIMARY KEY,
  patient_count INTEGER,
  receipt_count INTEGER,
  total_amount INTEGER -- in cents
);
```

### HTMX Patterns

**Notification bell badge refresh**:
```html
<!-- Header partial -->
<div hx-get="/notifications/unread-count" 
     hx-trigger="every 30s" 
     hx-swap="outerHTML">
  <span class="badge">{{.UnreadCount}}</span>
</div>
```

**Chart data refresh on navigation**:
```html
<!-- Weekly chart partial -->
<canvas id="weeklyChart" 
        data-patients="{{.PatientCounts}}" 
        data-receipts="{{.ReceiptCounts}}" 
        data-totals="{{.DailyTotals}}">
</canvas>
```

**Flash messages for user operations**:
```html
<!-- Base template -->
{{if .FlashSuccess}}
  <div class="flash flash-success">{{.FlashSuccess}}</div>
{{end}}
{{if .FlashError}}
  <div class="flash flash-error">{{.FlashError}}</div>
{{end}}
```

### Cleanup Job

**Notification auto-cleanup**:
- Trigger: First app startup on Monday
- Query: `DELETE FROM notifications WHERE expires_at < CURRENT_TIMESTAMP`
- Implementation: Go `init()` function checks day of week, runs cleanup if Monday

### Testing Strategy

**Unit Tests** (Go):
- Permission checks (HasPermission middleware)
- Notification categorization logic
- Storage threshold calculations
- Backup status state machine
- Weekly metrics aggregation query

**Integration Tests**:
- Dashboard page renders with all widgets
- Weekly metrics query returns correct aggregation
- Notification bell badge updates on new critical notification
- Sidebar navigation highlights active item
- HTMX partials load correctly (notification center, chart refresh)

**Manual Testing** (Playwright):
- Create receipt flow from dashboard quick action
- Click backup status card → navigate to Settings > Backup
- Notification bell shows count, dropdown opens, mark as read
- Chart renders correctly with combo chart (bars + line)

## When to Revisit

This decision should be revisited if:
- Multi-user support becomes a requirement (add roles table, permission inheritance)
- Notification volume increases (need filtering, categories expansion)
- Dashboard performance degrades (need query optimization, caching strategy)
- Chart requirements expand (need dedicated analytics page)
- HTMX coordination becomes complex (consider Alpine.js for client state)

## References

- **Domain Model**: `CONTEXT.md`, `docs/hcc-data-model.md`
- **UI Reference**: TailAdmin (https://github.com/TailAdmin/free-react-tailwind-admin-dashboard)
- **Chart Library**: Chart.js (https://www.chartjs.org)
- **Icons**: Heroicons (https://heroicons.com)
- **HTMX**: https://htmx.org
- **Go Templates**: https://pkg.go.dev/html/template

---

**Last updated**: 2026-07-17  
**Maintained by**: Development Team