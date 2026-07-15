# Clinic Management System — Prototype Specification

## Overview

A self-hosted desktop application for clinic receipt management, designed for single-practitioner use. The application runs as a local server process and uses the browser as its UI.

**Key characteristics:**
- Single binary deployment (Go + embedded SQLite)
- Server-rendered HTML with HTMX for interactivity
- TailwindCSS for styling (TailAdmin-inspired dashboard)
- No external network dependency — runs entirely on localhost

---

## Domain Model

### Core Concepts

#### Receipt

A patient visit record containing:
- **Patient identity**: name, HKID (Hong Kong Identity Card), gender
- **Visit details**: diagnosis, visit date, treatments provided
- **Financial data**: line items, subtotal, discounts, grand total
- **Administrative metadata**: receipt number, issue date, practitioner details

**Lifecycle:** Created as draft → Validated → Saved to SQLite → Optionally printed → Archived after retention period

#### Patient

A person receiving treatment:
- **Name**: Full name (Unicode supported)
- **HKID**: Hong Kong Identity Card number (format validated)
- **Gender**: Male / Female / Other
- **Visit history**: List of associated receipts

#### Clinic Settings

Configuration applied to all receipts:
- Clinic name, address, phone
- Practitioner name, registration number
- Default discount policies
- Receipt numbering format

#### Financial Year

April 1 — March 31 (Hong Kong standard). Used for reporting and analytics.

---

## Glossary

**Receipt Number**: Unique identifier for each receipt (timestamp + random suffix). Generated on final save.

**Draft Receipt**: In-progress receipt not yet saved to database. Stored in temporary table or session.

**HKID**: Hong Kong Identity Card number. Format: 1-2 letters + 6 digits + check digit in parentheses (e.g., `A123456(7)`).

**Grand Total**: Final amount after all discounts applied.

**Line Item**: Individual treatment or service on a receipt. Includes description, quantity, unit price, and subtotal.

**Delta Backup**: Backup containing only changes since last backup (privacy-masked).

**Masked Data**: Patient data with HKID and name partially obscured for privacy (e.g., `A*****6(7)`).

**Retention Period**: 3 years. Receipts older than this are automatically purged per regulatory requirements.

**Practitioner**: The medical professional (e.g., Tui Na therapist, acupuncturist) providing services.

---

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│              Single Binary (Go + Hugo)                  │
│  ┌───────────────────────────────────────────────────┐  │
│  │  HTTP Server (localhost:PORT)                     │  │
│  │  - Route handlers                                 │  │
│  │  - Session management                             │  │
│  │  - Template rendering                             │  │
│  └───────────────────────────────────────────────────┘  │
│  ┌───────────────────────────────────────────────────┐  │
│  │  SQLite Database (embedded)                       │  │
│  │  - receipts table                                 │  │
│  │  - patients table                                 │  │
│  │  - settings table                                 │  │
│  │  - backups table                                  │  │
│  └───────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────┘
           ↓ serves to
┌─────────────────────────────────────────────────────────┐
│              Browser (http://localhost:PORT)            │
│  ┌───────────────────────────────────────────────────┐  │
│  │  TailwindCSS (styled HTML)                        │  │
│  │  - Dashboard layout                               │  │
│  │  - Receipt forms                                  │  │
│  │  - Tables and charts                              │  │
│  └───────────────────────────────────────────────────┘  │
│  ┌───────────────────────────────────────────────────┐  │
│  │  HTMX                                             │  │
│  │  - Partial page updates                           │  │
│  │  - Form submissions                               │  │
│  │  - Validation feedback                            │  │
│  └───────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────┘
```

**Characteristics:**
- Single binary deployment
- Server-rendered HTML
- HTMX for dynamic interactions
- SQLite for persistent storage
- No build step (Hugo static generation + Go HTTP server)

---

## Technology Stack

### Backend

| Component | Technology | Purpose |
|-----------|------------|---------|
| **Language** | Go 1.21+ | HTTP server, business logic |
| **Web Framework** | Hugo (Go templates) | Server-side rendering |
| **Database** | SQLite 3 | Embedded relational database |
| **HTTP Server** | Go `net/http` (or `chi` router) | Localhost server |

### Frontend

| Component | Technology | Purpose |
|-----------|------------|---------|
| **Styling** | TailwindCSS v4 | Utility-first CSS |
| **Interactivity** | HTMX 2.x | AJAX, partial updates |
| **Layout** | TailAdmin-inspired | Dashboard UI patterns |
| **Icons** | Tailwind Icons | UI icons |

### Development Tools

| Tool | Purpose |
|------|---------|
| `go mod` | Dependency management |
| `hugo` | Template generation (optional, can use Go templates directly) |
| `tailwindcss` CLI | CSS compilation |
| `air` | Hot reload for development |

---

## UI Requirements

### Dashboard Layout

**Structure:**
- **Sidebar** (persistent, collapsible on mobile)
  - Navigation: Dashboard, Receipts, Patients, Reports, Settings
  - Icons + labels
- **Header** (top bar)
  - Search bar
  - Theme toggle (light/dark)
  - User profile / settings dropdown
- **Main Content Area**
  - Responsive grid layout
  - Card-based components

### Pages

#### 1. Dashboard (Home)

**Widgets:**
- Today's receipts count
- This month's revenue
- Recent patients list
- Quick actions (New Receipt, Search)

#### 2. Receipts List

**Features:**
- Paginated table
- Filter by date range
- Search by patient name / HKID / receipt number
- Actions: View, Edit, Print, Delete

#### 3. Receipt Form

**Sections:**
- Patient info (name, HKID, gender) — with autocomplete
- Visit details (date, diagnosis)
- Line items (description, quantity, price) — add/remove rows
- Discount (percentage or fixed amount)
- Totals (subtotal, discount, grand total) — auto-calculated
- Actions: Save Draft, Finalize, Print

**Validation:**
- Required fields highlighted
- HKID format validation
- Real-time total calculations
- Error messages displayed inline

#### 4. Patient List

**Features:**
- Paginated table
- Search by name / HKID
- View visit history
- Actions: Edit, View Receipts

#### 5. Reports

**Types:**
- Daily receipts summary
- Monthly revenue report
- Financial year report (April 1 — March 31)
- Patient visit frequency

**Export:**
- CSV download
- Print-friendly layout

#### 6. Settings

**Sections:**
- Clinic info (name, address, phone)
- Practitioner details
- Receipt numbering format
- Backup settings
- Data retention policy

---

## Data Model (SQLite Schema)

```sql
-- Patients table
CREATE TABLE patients (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    hkid TEXT UNIQUE NOT NULL,
    hkid_hash TEXT NOT NULL,  -- SHA-256 for deduplication
    gender TEXT CHECK(gender IN ('M', 'F', 'O')),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Receipts table
CREATE TABLE receipts (
    id TEXT PRIMARY KEY,
    receipt_number TEXT UNIQUE NOT NULL,
    patient_id TEXT NOT NULL,
    visit_date DATE NOT NULL,
    diagnosis TEXT,
    subtotal INTEGER NOT NULL,  -- in cents
    discount_type TEXT CHECK(discount_type IN ('percent', 'fixed', 'none')),
    discount_value INTEGER DEFAULT 0,  -- in cents or percentage points
    grand_total INTEGER NOT NULL,  -- in cents
    status TEXT CHECK(status IN ('draft', 'finalized', 'archived')),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (patient_id) REFERENCES patients(id)
);

-- Receipt line items
CREATE TABLE receipt_items (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    receipt_id TEXT NOT NULL,
    description TEXT NOT NULL,
    quantity INTEGER NOT NULL,
    unit_price INTEGER NOT NULL,  -- in cents
    subtotal INTEGER NOT NULL,  -- in cents
    FOREIGN KEY (receipt_id) REFERENCES receipts(id) ON DELETE CASCADE
);

-- Clinic settings (single row)
CREATE TABLE settings (
    id INTEGER PRIMARY KEY CHECK(id = 1),
    clinic_name TEXT NOT NULL,
    clinic_address TEXT,
    clinic_phone TEXT,
    practitioner_name TEXT NOT NULL,
    practitioner_registration TEXT,
    receipt_prefix TEXT DEFAULT 'RCP',
    retention_years INTEGER DEFAULT 3,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Auto-backups log
CREATE TABLE backups (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    backup_type TEXT CHECK(backup_type IN ('delta', 'full')),
    backup_date DATE NOT NULL,
    file_path TEXT,
    masked INTEGER DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX idx_receipts_patient ON receipts(patient_id);
CREATE INDEX idx_receipts_visit_date ON receipts(visit_date);
CREATE INDEX idx_receipts_status ON receipts(status);
CREATE INDEX idx_patients_hkid ON patients(hkid);
```

---

## Key Invariants

1. **Receipt number uniqueness**: Each receipt has a unique identifier (timestamp + random suffix)
2. **Data durability**: Critical data written to SQLite with WAL mode for reliability
3. **Privacy by default**: All exports default to masked data; full exports require explicit confirmation
4. **Retention enforcement**: Receipts older than 3 years are automatically purged
5. **Currency precision**: All monetary values stored as integers (cents), never floats

---

## Prototyping Steps

### Phase 1: Setup

1. Initialize Go module
2. Set up SQLite database with schema
3. Create basic HTTP server (localhost only)
4. Configure TailwindCSS v4

### Phase 2: Core CRUD

1. Patient management (list, create, edit)
2. Receipt form (single page, full post)
3. Receipt list (paginated table)
4. Settings page

### Phase 3: HTMX Interactivity

1. Form validation with fragment updates
2. Dynamic line items (add/remove rows)
3. Real-time total calculations
4. Patient autocomplete

### Phase 4: Dashboard & Reports

1. Dashboard widgets
2. Financial reports
3. CSV export
4. Print layouts

### Phase 5: Polish

1. Auto-backup scheduling
2. Retention policy enforcement
3. Theme toggle (light/dark)
4. Mobile-responsive adjustments

---

## References

- **Domain Model**: `DOMAIN.md`
- **Quick Start Guide**: `QUICKSTART.md`
- **Architecture Decision**: `docs/adr/0004-hugo-htmx-sqlite-architecture.md`
- **Decision Matrix**: `hugo-htmx-sqlite.md`