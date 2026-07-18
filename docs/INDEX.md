# Clinic Management — Documentation Index

## 📚 Documentation Set for Go + HTMX + SQLite Prototype

Complete specification for building a self-hosted clinic management application with Go + HTMX + SQLite.

---

## 🎯 Purpose

**Goal:** Build a simple, self-hosted desktop application for clinic receipt management.

**Approach:**
- Single binary deployment (Go + embedded SQLite)
- Server-rendered HTML with TailwindCSS
- "Just enough" interactivity via HTMX
- Team-friendly stack (Go + HTML)

---

## 📖 Document Overview

### 1. **Domain Model** (`DOMAIN.md`)

**Purpose:** Pure domain knowledge, implementation-agnostic.

**Contents:**
- Core entities (Receipt, Patient, Settings, Line Item)
- Value objects (Money, HKID)
- Business rules (receipt numbering, discount calculation, retention policy)
- Validation rules
- Use cases (create receipt, search patient, generate reports)
- Domain events
- Glossary

**Use when:** Implementing business logic in any technology stack.

---

### 2. **Prototype Specification** (`prototype-spec.md`)

**Purpose:** Technical specification for the Go + HTMX + SQLite prototype.

**Contents:**
- Technology stack details
- UI requirements (dashboard layout, pages, widgets)
- SQLite database schema
- Key invariants
- Prototyping phases (5 phases from setup to polish)

**Use when:** Starting the prototype implementation.

---

### 3. **Quick Start Guide** (`QUICKSTART.md`)

**Purpose:** Step-by-step setup instructions for the new workspace.

**Contents:**
- Prerequisites
- Project initialization
- Directory structure
- Database setup
- HTTP server example (Go + Chi router)
- Handler examples
- HTMX form examples
- TailwindCSS configuration
- Development workflow
- Deployment instructions

**Use when:** Setting up the new Go + HTMX workspace.

---

### 4. **Architecture Decision Record** (`docs/adr/0004-hugo-htmx-sqlite-architecture.md`)

**Purpose:** Documents the decision to adopt Go + HTMX + SQLite, including trade-offs.

**Contents:**
- Context (requirements for simple CRUD app)
- Considered options (React, Go, Python, Blazor)
- Decision rationale
- Consequences (positive/negative)
- Testing strategy
- When to revisit

**Use when:** Understanding why this architecture was chosen.

---

### 5. **Decision Matrix** (DELETED - was `hugo-htmx-sqlite.md`)

**Status:** Deleted — contained misleading "Hugo" naming.

**Replacement:** See `docs/adr/0004-hugo-htmx-sqlite-architecture.md` for the corrected architecture decision.

---

## 🚀 Getting Started

### Step 1: Read the Domain Model

Start with `DOMAIN.md` to understand the business logic independent of technology.

### Step 2: Review the Architecture

Read `docs/adr/0004-hugo-htmx-sqlite-architecture.md` to understand the trade-offs.

### Step 3: Follow the Quick Start

Use `QUICKSTART.md` to set up the new workspace.

### Step 4: Implement in Phases

Follow the 5 phases in `prototype-spec.md`:
1. **Setup** — Go module, SQLite, HTTP server, Tailwind
2. **Core CRUD** — Patient management, receipt form, receipts list
3. **HTMX Interactivity** — Validation, dynamic line items, calculations
4. **Dashboard & Reports** — Widgets, financial reports, CSV export
5. **Polish** — Auto-backup, retention, theme toggle, mobile

---

## 📋 File Locations

```
docs/
├── INDEX.md                          ← This file
├── DOMAIN.md                         ← Pure domain knowledge
├── prototype-spec.md                 ← Technical specification
├── QUICKSTART.md                     ← Setup guide
└── adr\
    └── 0004-hugo-htmx-sqlite-architecture.md  ← Architecture decision
    └── 0005-dashboard-layout-navigation.md    ← Dashboard UI
    └── 0006-go-standard-library-for-exports.md ← CSV/PDF/Excel exports
```

---

## 🔑 Key Concepts

### Self-Hosted Desktop App

**Definition:** Single binary that runs on desktop and opens browser to `http://localhost:PORT`.

**Characteristics:**
- No external network dependency
- Embedded SQLite database
- Browser as UI (not a native window)
- Server-rendered HTML

**Not "Local-First" in PWA sense:** Requires server process running; no offline support.

---

### Server-Rendered HTML with HTMX

**Pattern:**
1. Browser requests page
2. Server renders full HTML
3. Browser displays page
4. User interacts (form submit, button click)
5. HTMX sends AJAX request
6. Server returns HTML fragment
7. HTMX swaps fragment into page

**Benefits:**
- No client-side state management
- Server validates everything
- "Just enough" interactivity

**Trade-offs:**
- Network round-trips for interactions
- No offline support
- HTMX coordination complexity

---

### TailAdmin-Inspired Dashboard

**Layout Pattern:**
- Fixed sidebar (navigation)
- Top header (search, profile, theme toggle)
- Main content area (responsive grid of cards)

**Components:**
- StatCard (metrics)
- TableCard (data tables)
- ChartCard (visualizations)
- FormCard (input forms)

**Styling:**
- TailwindCSS v4
- Tailwind Icons
- Dark mode support

---

## 🎯 Success Criteria

The prototype is successful when:
- ✅ All CRUD operations work (create, read, update, delete receipts/patients)
- ✅ HTMX provides smooth interactivity (validation, line items, calculations)
- ✅ Dashboard displays key metrics
- ✅ Reports generate correctly (monthly, financial year)
- ✅ CSV export works (masked by default)
- ✅ Auto-backup runs daily
- ✅ Retention policy enforced (3 years)
- ✅ Single binary deployment works

---

## 📞 Next Steps

1. **Create new workspace** — Separate Git repository
2. **Initialize Go module** — `go mod init clinic-app`
3. **Set up SQLite** — Create database with schema from `prototype-spec.md`
4. **Build HTTP server** — Follow `QUICKSTART.md` examples
5. **Implement Phase 1** — Setup and core CRUD
6. **Iterate through phases** — Add HTMX, dashboard, polish

---

## 📚 Additional Resources

- **HTMX Reference:** https://htmx.org/reference/
- **TailwindCSS v4:** https://tailwindcss.com/docs
- **Go SQLite:** https://github.com/mattn/go-sqlite3
- **Chi Router:** https://github.com/go-chi/chi
- **TailAdmin (Demo):** https://tailadmin.com/

---

## 🏷️ Version

**Documentation Version:** 1.0  
**Last Updated:** 2026-07-15  
**Maintained By:** Development Team

---

## 📝 License

This documentation is provided as-is for prototyping purposes.