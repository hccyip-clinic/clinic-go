---
status: closed
assignee: copilot
created: 2026-07-23
wayfinder:grilling
parent: 000-map
---

# Define receipt archival policy

## Question

When and how does a finalized receipt become archived, and is archival automatic at the retention boundary, manually triggered, or both while remaining read-only?

## Resolution

Archival is deferred beyond Phase 1. Finalized receipts remain in the `finalized` state; Phase 1 exposes no manual archive transition and performs no purge or retention job. The existing `archived` status remains in the schema for forward compatibility, but application code must not transition receipts into it yet.

The future policy is read-only archival after three years from the end of the applicable Hong Kong financial year, executed by a later scheduled retention workflow. That workflow is separate from the Phase 1 vertical slice.
