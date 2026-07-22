# Clinic Management System — Context

## Business Overview

A self-hosted desktop application for **clinic receipt management**, designed for single-practitioner traditional Chinese medicine clinics in Hong Kong.

**Primary purpose:** Generate billing receipts, print receipts, and issue sick leave certificates for patient visits.

**Deployment model:** Single binary (Go + embedded SQLite) running on localhost, accessed via browser.

**Architecture Clarification:** This is NOT a Hugo static site. The stack is:
- **Backend:** Go HTTP server with `html/template` for server-rendered HTML
- **Dynamic UI:** HTMX for partial page updates without full React complexity
- **Data:** Embedded SQLite database
- **Hugo's role:** None — Hugo is a static site generator and is NOT used in this project

---

## Key User Roles

### Practitioner
- Medical professional providing treatments (Tui Na, acupuncture, internal medicine)
- Creates and finalizes receipts
- Views patient history
- Generates reports for HK Inland Revenue Department

### Admin (same person as Practitioner in single-practitioner setup)
- Configures clinic settings
- Manages data retention and backups
- Exports financial reports

---

## Core Workflows

### 1. Create Receipt
Patient visit → Enter treatments → Calculate totals → Finalize → Print receipt

### 2. Patient Search
Search by name/HKID → Select patient → Auto-fill details → Link to receipt

### 3. Generate Reports
Select period (daily/monthly/financial year) → Calculate totals → Export CSV (masked by default)

### 4. Auto-Backup
Daily at midnight → Mask patient data → Write delta backup → Enforce 30-day retention

---

## Domain Terminology

**Receipt:** A patient visit record with clinical and financial information. Lifecycle: draft → finalized → archived.

**HKID:** Hong Kong Identity Card number. Input may contain spaces or hyphens, but is normalized to uppercase canonical form `[A-Z]{1,2}[0-9]{6}([0-9A])` before validation and storage. Both one- and two-letter prefixes use the official check digit algorithm.

**Financial Year:** April 1 — March 31 (Hong Kong standard).

**Grand Total:** Final amount after discount. All values stored as integer cents.

**Receipt Number:** Unique identifier. Format: `{prefix}-{YYYYMMDD}-{random}` (e.g., `RCP-20260715-A3F2B`).

**Masked Data:** Patient data with HKID/name partially obscured (e.g., `A*****6(7)`) for privacy in exports/backups.

**Retention Period:** 3 years from financial year end. Receipts auto-purged after this period.

**Tui Na:** Chinese therapeutic massage.

**Acupuncture:** Traditional Chinese medicine needle therapy.

**Internal Medicine:** Traditional Chinese herbal medicine prescription.

**Dashboard:** Home landing page showing quick actions, system status cards, and weekly metrics chart.

**Quick Actions:** Curated set of 3 primary workflow entry points on dashboard (Create Receipt, Search Patient, View Today's Receipts).

**System Status:** At-a-glance widgets showing backup health and storage metrics with actionable thresholds.

**Notifications:** System-generated alerts categorized as Critical/Warning/Info (persistent, stored in SQLite, 7-day retention) or Success/Fail (ephemeral flash messages). Bell badge shows unread critical + fail count for today.

**Permissions:** Direct permission strings assigned to users (e.g., `receipts:create`, `patients:read`). No role abstraction in Phase 1. Defined as Go constants in `handlers/permissions.go`.

**Financial Week:** Monday–Saturday (Hong Kong business week). Dashboard weekly metrics chart shows Mon→Today.

**Session:** In-memory session store with httpOnly cookie expiry. Single-user desktop app — restart clears sessions.

**Clinic Settings:** Single-row SQLite table (`clinic_settings`) storing clinic name, address, telephone. Default from `config.json` on first run.

**Patient Identity:** A patient has one unique canonical HKID. Names and gender may be updated, but the HKID is immutable after creation in Phase 1.

---

## System Boundaries

### In Scope (Phase 1)
- Receipt creation and management
- Receipt printing
- Patient database with HKID validation
- Basic financial reports
- Privacy-masked exports
- Auto-backup with retention

### Out of Scope (Phase 1)
- Payment channel integration (system records payments, doesn't process them)
- Insurance claims processing
- Appointment scheduling

### Future Phases
- Patient health history tracking
- Monthly income analytics
- Financial year statements for HK Inland Revenue Department

---

## Key Business Rules

1. **No sales tax in Hong Kong** — receipts show totals without tax
2. **HKID must be unique** — validated with format regex and check digit algorithm
3. **Receipt numbers generated only on finalization** — not for drafts
4. **All monetary values in integer cents** — no floating-point arithmetic
5. **Privacy by default** — exports masked unless user explicitly confirms full data access
6. **Financial year runs April 1 — March 31** — not calendar year
7. **Full Unicode/UTF-8 support** — patient names, diagnosis, and receipts must display Chinese characters correctly

---

## External Dependencies

### Hong Kong Inland Revenue Department
- Financial year statements must follow April 1 — March 31 format
- Receipts retained for 3 years per regulatory requirements

### Privacy Ordinance (Hong Kong)
- Patient data must be protected in backups/exports
- Masking required for data leaving the clinic system

---

## Success Metrics

- Receipt creation time: < 2 minutes per patient
- Zero floating-point rounding errors in financial calculations
- 100% HKID format validation accuracy
- Automated daily backups with 30-day retention
- CSV exports compliant with privacy requirements