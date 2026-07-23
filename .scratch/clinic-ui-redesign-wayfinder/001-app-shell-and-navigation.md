---
status: done
created: 2026-07-23
wayfinder:grilling
parent: 000-map
assignee: Copilot
---

# Establish the app shell and navigation hierarchy

## Question

What should the authenticated clinic app shell, global header, sidebar/navigation hierarchy, page framing, active states, and cross-screen orientation cues be so that a practitioner can always understand where they are and what to do next?

## Resolution

Use a persistent left sidebar as the desktop app shell, collapsing to a menu button on smaller screens. Group the navigation into:

- **Primary**: Dashboard, Receipts, Patients
- **Administration**: Settings

Keep Reports out of the visible navigation until a real route and screen exist; do not create dead-end or disabled destinations. Use a compact global header containing page context, clinic identity, and session actions, while deferring notifications until notification data and behavior exist. Highlight the active sidebar destination, add short breadcrumbs for nested pages such as receipt and patient views, and place the page title plus one clear primary action in each content header. Labels remain text-first, with icons as supporting cues.
