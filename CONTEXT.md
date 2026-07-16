# Clinic Management System — Context

## Business Overview

A self-hosted desktop application for **clinic receipt management**, designed for single-practitioner traditional Chinese medicine clinics in Hong Kong.

**Primary purpose:** Generate billing receipts, print receipts, and issue sick leave certificates for patient visits.

**Deployment model:** Single binary (Go + embedded SQLite) running on localhost, accessed via browser.

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

**HKID:** Hong Kong Identity Card number. Format: `[A-Z]{1,2}[0-9]{6}([0-9A])`. Validated with check digit algorithm.

**Financial Year:** April 1 — March 31 (Hong Kong standard).

**Grand Total:** Final amount after discount. All values stored as integer cents.

**Receipt Number:** Unique identifier. Format: `{prefix}-{YYYYMMDD}-{random}` (e.g., `RCP-20260715-A3F2B`).

**Masked Data:** Patient data with HKID/name partially obscured (e.g., `A*****6(7)`) for privacy in exports/backups.

**Retention Period:** 3 years from financial year end. Receipts auto-purged after this period.

**Tui Na:** Chinese therapeutic massage.

**Acupuncture:** Traditional Chinese medicine needle therapy.

**Internal Medicine:** Traditional Chinese herbal medicine prescription.

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