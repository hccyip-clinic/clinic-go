---
status: done
created: 2026-07-23
wayfinder:grilling
parent: 000-map
blocked-by: 001-app-shell-and-navigation, 002-visual-design-system, 003-workflow-screen-patterns
assignee: Copilot
---

# Define responsive, accessible, and bilingual content rules

## Question

What responsive breakpoints, keyboard and focus behavior, contrast requirements, table/form adaptations, Traditional Chinese text expansion rules, and narrow-screen fallbacks must every redesigned screen satisfy?

## Resolution

Apply these cross-screen rules:

- **Responsive behavior**: target desktop at 1024px and above, tablet at 768–1023px with a collapsible sidebar and reorganized grids, and a narrow-screen fallback below 768px with stacked content and scrollable data tables.
- **Accessibility**: follow WCAG 2.2 AA practices with full keyboard operation, visible focus indicators, 4.5:1 body-text contrast, semantic landmarks and headings, explicit labels, announced validation/status changes, and no color-only meaning.
- **Forms**: stack receipt sections on smaller screens and keep the primary Save draft action reachable in a bottom action area without covering fields or keyboard input.
- **Tables**: preserve tabular meaning with horizontal scrolling and a readable minimum width; reduce only secondary columns when an equivalent remains available.
- **Bilingual content**: keep English UI labels, allow flexible containers for Traditional Chinese names, diagnoses, and addresses, never truncate patient-identifying data, use locale-safe date and HKD formatting, and test mixed-language content at 200% zoom.
- **Interaction sizing**: require at least 44×44px touch targets and ensure all actions work without hover.
