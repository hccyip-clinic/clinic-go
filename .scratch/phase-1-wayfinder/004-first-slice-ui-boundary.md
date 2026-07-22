---
status: done
created: 2026-07-22
wayfinder:prototype
parent: 000-map
blocked-by: 002-receipt-lifecycle, 003-patient-identity-and-validation
---

# Shape the first server-rendered UI slice

## Question

What is the smallest coherent set of Go templates, routes, forms, and local HTMX interactions that demonstrates patient creation, receipt draft/finalization, and receipt listing without prematurely implementing the full dashboard?

## Resolution

Use Variant A, the guided flow, as the first production UI boundary: patient details, visit details, charges, and review are explicit steps, with draft save available before finalization. The full three-variant prototype is retained at `clinic-hcc-app/static/prototypes/phase1-ui.html` until the production templates are implemented; it is not production UI.
