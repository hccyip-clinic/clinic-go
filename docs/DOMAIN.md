# Clinic Management Domain Model

## Purpose

This document extracts the **pure domain knowledge** from the clinic management system, independent of any implementation technology (Vite, React, Hugo, Go, etc.).

Use this as the authoritative reference when implementing the application in any technology stack.

---

## Core Entities

### Receipt

**Definition:** A record of a patient visit containing clinical and financial information.

**Attributes:**
- `receiptNumber` (string, unique) — Format: `RCP-YYYYMMDD-XXXXX` (prefix + date + random)
- `patientId` (string) — Reference to Patient
- `visitDate` (date) — Date of consultation
- `diagnosis` (string, optional) — Clinical diagnosis
- `lineItems` (array) — List of treatments/services
- `subtotal` (integer, cents) — Sum of line item subtotals
- `discountType` (enum) — `none`, `percent`, `fixed`
- `discountValue` (integer) — Percentage points (0-100) or cents
- `grandTotal` (integer, cents) — Final amount after discount
- `status` (enum) — `draft`, `finalized`, `archived`
- `createdAt` (datetime)
- `updatedAt` (datetime)

**Invariants:**
1. Receipt number must be unique across all receipts
2. Grand total = subtotal - discount (calculated, not stored independently)
3. All monetary values stored as integers (cents), never floats
4. Receipt number generated only on finalization (not for drafts)

**Lifecycle:**
1. Created as `draft` (incomplete, editable)
2. Validated (all required fields present, format checks pass)
3. Finalized (receipt number assigned, immutable except for corrections)
4. Archived (after retention period, read-only)

---

### Patient

**Definition:** A person receiving treatment at the clinic.

**Attributes:**
- `id` (string, UUID) — Unique identifier
- `name` (string) — Full name (Unicode supported)
- `hkid` (string) — Hong Kong Identity Card number
- `hkidHash` (string, SHA-256) — Hash for deduplication without exposing raw HKID
- `gender` (enum) — `M`, `F`, `O` (Male, Female, Other)
- `createdAt` (datetime)
- `updatedAt` (datetime)

**Invariants:**
1. HKID must be unique across all patients
2. HKID format: 1-2 letters + 6 digits + check digit in parentheses (e.g., `A123456(7)`)
3. HKID hash derived from uppercase HKID without formatting (e.g., `A1234567`)

**HKID Validation Rules:**
- Format regex: `/^[A-Z]{1,2}[0-9]{6}\([0-9A]\)$/`
- Check digit algorithm: Weighted sum modulo 11
- Example: `A123456(7)` → valid, `A123456(8)` → invalid

---

### Clinic Settings

**Definition:** Global configuration applied to all receipts and documents.

**Attributes:**
- `clinicName` (string)
- `clinicAddress` (string)
- `clinicPhone` (string)
- `practitionerName` (string)
- `practitionerRegistration` (string) — Medical council registration number
- `receiptPrefix` (string, default: `RCP`)
- `retentionYears` (integer, default: 3)

**Invariants:**
1. Only one settings record exists (singleton)
2. All fields required except address

---

### Line Item

**Definition:** An individual treatment or service on a receipt.

**Attributes:**
- `description` (string) — Treatment name/description
- `quantity` (integer) — Number of units
- `unitPrice` (integer, cents) — Price per unit
- `subtotal` (integer, cents) — `quantity × unitPrice`

**Invariants:**
1. Subtotal calculated, not stored independently
2. Quantity must be positive (> 0)
3. Unit price must be non-negative (≥ 0)

---

## Value Objects

### Money

**Definition:** Monetary amount stored as integer cents.

**Operations:**
- `add(Money, Money) → Money`
- `subtract(Money, Money) → Money`
- `multiply(Money, integer) → Money`
- `applyPercent(Money, percent) → Money`

**Rules:**
1. Never use floating-point arithmetic
2. Always round down (floor) for discounts
3. Display format: `$XX.XX` (divide by 100, format with 2 decimals)

---

### HKID

**Definition:** Hong Kong Identity Card number.

**Format:** `[A-Z]{1,2}[0-9]{6}([0-9A])`

**Examples:**
- `A123456(7)` — Valid
- `AB123456(7)` — Valid (2 letters)
- `A123456(8)` — Invalid (wrong check digit)

**Masking Rule:**
- Display format: `A*****6(7)` (first char + asterisks + last 2 chars)
- Full HKID only shown to authorized users

---

## Business Rules

### Receipt Numbering

**Format:** `{prefix}-{YYYYMMDD}-{random}`

**Example:** `RCP-20260715-A3F2B`

**Rules:**
1. Prefix from settings (default: `RCP`)
2. Date is receipt creation date (not visit date)
3. Random suffix: 5 alphanumeric characters (case-sensitive)
4. Generated only on finalization (not for drafts)
5. Must be unique across all receipts

---

### Discount Calculation

**Formula:**
```
if discountType == 'none':
    grandTotal = subtotal
elif discountType == 'percent':
    grandTotal = subtotal - floor(subtotal * discountValue / 100)
elif discountType == 'fixed':
    grandTotal = max(0, subtotal - discountValue)
```

**Rules:**
1. Percent discounts rounded down (floor)
2. Fixed discounts cannot exceed subtotal (grand total ≥ 0)
3. Discount applied to subtotal, not individual line items

---

### Financial Year

**Definition:** April 1 — March 31 (Hong Kong standard)

**Usage:**
- Reports grouped by financial year
- Retention policy calculated from financial year end
- Tax documents organized by financial year

**Calculation:**
```
func financialYear(date Date) FinancialYear {
    if date.Month >= 4:
        return FinancialYear{Start: date.Year, End: date.Year + 1}
    else:
        return FinancialYear{Start: date.Year - 1, End: date.Year}
}
```

---

### Retention Policy

**Rule:** Receipts older than 3 years are automatically purged.

**Calculation:**
- Retention period starts from **financial year end**, not receipt date
- Example: Receipt from June 2023 → Financial year 2023-2024 → Retention ends March 31, 2027

**Purge Schedule:**
- Run daily at midnight
- Soft delete (mark as `archived`) for 30 days before hard delete
- Log all purged receipts for audit trail

---

### Privacy Masking

**When to Mask:**
- Backups (delta and full)
- CSV exports
- Reports shared outside clinic
- Analytics data

**Masking Rules:**
- **HKID:** `A123456(7)` → `A*****6(7)` (preserve first char, last 2 chars)
- **Patient Name:** `John Doe` → `J*** D**` (preserve first letter of each word)
- **Patient ID Hash:** Always stored (not masked) for deduplication

**Unmasked Access:**
- Requires explicit user confirmation
- Logged for audit purposes
- Only for authorized users (practitioner, admin)

---

## Validation Rules

### Receipt Validation

**Required Fields:**
- Patient (name, HKID, gender)
- Visit date
- At least one line item
- Practitioner details (from settings)

**Format Validation:**
- HKID format (regex + check digit)
- Visit date (not in future)
- Receipt date (not in future)

**Business Validation:**
- Grand total > 0
- Line item subtotals match calculation
- Discount within reasonable range (0-100% for percent, 0-subtotal for fixed)

**Severity Levels:**
- **Error:** Blocks save (missing required field, invalid HKID)
- **Warning:** Allows save but notifies user (unusually high discount, future visit date)

---

### Patient Validation

**Required Fields:**
- Name (non-empty)
- HKID (valid format)
- Gender (M/F/O)

**Uniqueness:**
- HKID must be unique (check against existing patients)
- Use HKID hash for comparison (case-insensitive)

---

## Use Cases

### Create Receipt

**Preconditions:**
- Patient exists (or is created during receipt creation)
- Clinic settings configured

**Flow:**
1. User opens "New Receipt" form
2. System pre-fills practitioner details from settings
3. User enters/selects patient
4. User adds line items (description, quantity, price)
5. System calculates subtotal, discount, grand total in real-time
6. User clicks "Finalize"
7. System validates receipt
8. System generates receipt number
9. System saves receipt to database
10. System offers to print receipt

**Postconditions:**
- Receipt saved with status `finalized`
- Patient visit history updated
- Receipt number assigned and immutable

---

### Search Patient

**Preconditions:**
- Patient database exists

**Flow:**
1. User types in patient search field
2. System searches by name (partial match, case-insensitive)
3. System searches by HKID (exact match)
4. System displays matching patients (max 10 results)
5. User selects patient
6. System auto-fills patient details in form

**Postconditions:**
- Patient linked to current receipt
- Patient last-visited timestamp updated

---

### Generate Monthly Report

**Preconditions:**
- Receipts exist for selected month
- User authorized to view reports

**Flow:**
1. User selects month and year
2. System filters receipts by visit date
3. System groups by practitioner, treatment type
4. System calculates totals (count, revenue)
5. System displays summary table
6. User clicks "Export CSV"
7. System generates masked CSV (patient names/HKIDs masked)
8. System prompts for unmasked export (requires confirmation)

**Postconditions:**
- Report generated
- Export logged for audit

---

### Auto-Backup

**Preconditions:**
- Backup schedule configured
- Storage location available

**Flow:**
1. System runs daily at midnight
2. System identifies receipts created/updated since last backup
3. System masks patient data (name, HKID)
4. System writes delta backup to storage
5. System logs backup metadata (date, count, type)
6. System enforces retention (delete backups older than 30 days)

**Postconditions:**
- Delta backup created
- Old backups purged per retention policy

---

## Domain Events

### ReceiptFinalized

**When:** Receipt transitions from `draft` to `finalized`

**Data:**
- `receiptId` (string)
- `receiptNumber` (string)
- `patientId` (string)
- `grandTotal` (integer)
- `finalizedAt` (datetime)

**Consumers:**
- Auto-backup system (triggers immediate backup)
- Analytics (updates revenue dashboard)
- Notification system (optional: send receipt to patient)

---

### ReceiptPurged

**When:** Receipt deleted per retention policy

**Data:**
- `receiptId` (string)
- `receiptNumber` (string)
- `purgedAt` (datetime)
- `reason` (string) — "retention_policy"

**Consumers:**
- Audit log (records deletion for compliance)
- Backup system (marks backup for deletion)

---

### PatientIdentified

**When:** Existing patient identified by HKID during receipt creation

**Data:**
- `patientId` (string)
- `hkidHash` (string)
- `matchedAt` (datetime)

**Consumers:**
- Patient history (updates last-visited timestamp)
- Analytics (updates patient visit frequency)

---

## Glossary

**Archived Receipt:** Receipt older than retention period, marked read-only pending deletion.

**Delta Backup:** Backup containing only records changed since last backup.

**Draft Receipt:** In-progress receipt not yet finalized. Editable, no receipt number assigned.

**Financial Year:** April 1 — March 31 (Hong Kong standard).

**Finalized Receipt:** Completed receipt with assigned receipt number. Immutable except for corrections.

**Grand Total:** Final amount after discount applied.

**HKID:** Hong Kong Identity Card number. Format: `[A-Z]{1,2}[0-9]{6}([0-9A])`.

**HKID Hash:** SHA-256 hash of uppercase HKID without formatting. Used for deduplication.

**Line Item:** Individual treatment/service on a receipt.

**Masked Data:** Patient data with HKID/name partially obscured for privacy.

**Practitioner:** Medical professional providing services (e.g., Tui Na therapist, acupuncturist).

**Receipt Number:** Unique identifier for finalized receipt. Format: `{prefix}-{YYYYMMDD}-{random}`.

**Retention Period:** 3 years from financial year end. After this, receipts are purged.

**Subtotal:** Sum of all line item subtotals, before discount.

**Tui Na:** Chinese therapeutic massage.

**Acupuncture:** Traditional Chinese medicine needle therapy.

**Internal Medicine:** Traditional Chinese herbal medicine prescription.

---

## Implementation-Agnostic Principles

1. **Single Source of Truth** — Receipt state owned by one module (whether React App.tsx or Go handler)
2. **Immutable Updates** — State changes create new objects, don't mutate existing (whether in-memory or database rows)
3. **Validation with Severity** — Distinguish blocking errors from warnings
4. **Auto-Save** — Valid drafts persisted on every change
5. **Privacy by Default** — Exports default to masked data; full data requires explicit confirmation
6. **Currency Precision** — Integer math (cents), never floating-point
7. **Audit Trail** — All sensitive operations logged (exports, unmasked access, deletions)

---

## References

- **Full Context**: `CONTEXT.md` (implementation-specific details)
- **Prototype Spec**: `docs/prototype-spec.md` (Go + HTMX implementation)
- **Quick Start**: `docs/QUICKSTART.md` (getting started guide)
- **Architecture Decision**: `docs/adr/0004-hugo-htmx-sqlite-architecture.md`