---
status: done
created: 2026-07-23
wayfinder:grilling
parent: 000-map
blocked-by: 001-app-shell-and-navigation
assignee: Copilot
---

# Define the calm clinical visual system

## Question

Which palette, typography, spacing scale, surfaces, borders, elevation, status colors, icon treatment, focus states, and motion rules should make the clinic app feel modern, professional, clean, and trustworthy while remaining information-dense?

## Resolution

Adopt a restrained calm-clinical system:

- **Color**: deep teal as the primary accent for active navigation and primary actions; warm-neutral page surfaces; slate text; restrained green, amber, and red semantic states.
- **Typography**: one legible system sans-serif stack with a Traditional Chinese fallback such as Noto Sans CJK. Avoid decorative or network-dependent fonts.
- **Surfaces and depth**: warm off-white page canvas, white content surfaces, subtle 1px borders as the default separation, and shadows only for the sidebar, menus, and dialogs.
- **Spacing and density**: a 4/8-based scale with generous page and section spacing, balanced by compact tables and receipt controls.
- **Shape**: consistent moderate 8px rounding for cards and controls; reserve pill shapes for status badges.
- **Icons and motion**: one consistent 16–20px outline icon set, never emoji, with short purposeful transitions and reduced-motion support. No decorative animation.
