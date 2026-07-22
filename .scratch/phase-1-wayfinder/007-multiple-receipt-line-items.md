---
status: closed
assignee: copilot
created: 2026-07-23
wayfinder:grilling
parent: 000-map
---

# Decide multiple receipt line items

## Question

Should the Phase 1 receipt flow support adding, editing, and removing multiple treatment line items, or intentionally limit the first vertical slice to one line item while preserving the multi-item domain model?

## Resolution

Phase 1 production receipts support multiple line items. A new draft starts with one blank row; saving requires at least one valid row. Users may add, edit, and remove rows while the receipt is a draft, but finalized receipts and their line items are immutable.

Quantities are positive whole numbers, unit prices remain integer cents, and row subtotals are calculated as quantity multiplied by unit price. Line-item order is persisted explicitly so drafts, printed receipts, and later reads remain deterministic.
