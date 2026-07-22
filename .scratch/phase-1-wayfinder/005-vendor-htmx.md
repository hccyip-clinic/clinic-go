---
status: done
created: 2026-07-22
wayfinder:task
parent: 000-map
blocked-by: 004-first-slice-ui-boundary
---

# Vendor HTMX for offline deployment

## Question

What local HTMX asset and serving/build arrangement should Phase 1 use so the guided receipt flow has no runtime dependency on `unpkg.com` or any other external network?

## Resolution

Use the `htmx.org` npm package at v2.0.7 and vendor `dist/htmx.min.js` as `static/js/htmx.min.js`. The prototype loads it from `/static/js/htmx.min.js`, so the server's existing static-file route serves HTMX without a CDN dependency.
