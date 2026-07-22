---
status: closed
assignee: copilot
created: 2026-07-23
wayfinder:task
parent: 000-map
---

# Define the Tailwind build artifact

## Question

Should Phase 1 include a checked-in/generated `static/css/styles.css` built from a real `src/index.css`, or should styling remain an explicit follow-up while the server-rendered workflow is stabilized?

## Resolution

Tailwind CSS is part of the Phase 1 runtime artifact. Add `clinic-hcc-app/src/index.css` as the Tailwind v4 entry point, keep `npm run build:css` as the release build, and check in the generated `clinic-hcc-app/static/css/styles.css` because the Go server serves that file directly and deployment must not require Node.js.

The source must scan the server-rendered templates (and any production HTML asset paths) so every class used by the application is emitted. The existing layout already references `/static/css/styles.css`; a missing source and output is therefore an implementation defect, not a reason to defer styling.
