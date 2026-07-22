---
status: done
created: 2026-07-22
wayfinder:grilling
parent: 000-map
blocked-by: 001-sqlite-driver-migration
---

# Define the Phase 1 receipt lifecycle

## Question

What are the precise patient, draft, finalization, editing, deletion, and receipt-numbering rules for the first end-to-end receipt workflow, including the transaction boundary?

## Resolution

Phase 1 persists drafts in SQLite. Drafts may be edited or deleted; finalization validates the patient, line items, totals, and discount invariants in one all-or-nothing transaction, assigns a unique `RCP-YYYYMMDD-XXXXXX`-style number, and makes the receipt immutable. Finalized receipts may be viewed, printed, or archived, but not edited or deleted. All monetary values remain integer cents, quantities and prices cannot be negative, discounts cannot exceed the subtotal, and percentage discounts are limited to 0–100.
