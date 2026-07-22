---
status: done
created: 2026-07-22
wayfinder:grilling
parent: 000-map
blocked-by: 001-sqlite-driver-migration
---

# Define patient identity and validation

## Question

What normalized HKID, name, gender, uniqueness, search, and validation behavior must the Phase 1 patient workflow enforce, including how raw HKIDs and hashes are maintained?

## Resolution

Normalize HKID input by trimming whitespace, removing optional spaces and hyphens, uppercasing letters, and storing the canonical parenthesized form. Validate the official check digit for both one- and two-letter prefixes; reject invalid values. Store the canonical raw HKID for local lookup and a hash for deduplication, enforce uniqueness, and reject duplicate patient creation in favor of the existing patient. Names are required Unicode text after trimming, gender is one of `M`, `F`, or `O`, search supports name substring or exact canonical HKID matching, and HKIDs are immutable after patient creation in Phase 1.
