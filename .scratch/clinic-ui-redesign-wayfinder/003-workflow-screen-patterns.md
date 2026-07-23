---
status: done
created: 2026-07-23
wayfinder:grilling
parent: 000-map
blocked-by: 001-app-shell-and-navigation, 002-visual-design-system
assignee: Copilot
---

# Shape receipt and patient workflow screens

## Question

How should the dashboard, receipt list/form/view, patient search/form, reports, settings, flash messages, validation errors, and destructive actions be structured and sequenced so the core clinic workflows are fast, clear, and recoverable?

## Resolution

Use these screen patterns:

- **Dashboard**: show a short context line, then the three quick actions with Create receipt first; place operational status and recent receipts after the actions, with analytics secondary.
- **Receipt creation**: keep one page with numbered Patient, Visit, Charges, and Review/Actions sections. Use a searchable patient picker with a nearby create-new path, repeatable treatment rows, inline numeric validation, and a visible live totals summary.
- **Receipt lifecycle**: use a stable action bar with Save draft as the safe editing action, Finalize as a distinct confirmation action on the review/view screen, and Delete draft separated as a confirmed destructive action.
- **Patient management**: use a search-first page with New patient at the top, a scannable results table, and a dedicated detail/edit view that can later show receipt history.
- **Feedback and errors**: show field-level errors, add a concise form summary when multiple fields fail, preserve entered values after errors, and show clear dismissible success feedback after save or finalization. Never rely on color alone.
