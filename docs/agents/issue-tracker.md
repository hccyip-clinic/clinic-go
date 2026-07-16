# Issue Tracker

This repository uses **local markdown files** as its issue tracker.

## Location

Issues are stored under `.scratch/<feature>/` as markdown files.

## Workflow

Skills that create issues (`to-tickets`, `to-spec`, `qa`, `triage`) will:
1. Create a `.scratch/<feature-name>/` directory for the feature or epic
2. Write individual issue files as `<issue-number>-<slug>.md` within that directory
3. Track state via frontmatter or status badges in the markdown files

## Issue File Template

```markdown
---
status: needs-triage
created: YYYY-MM-DD
---

# Issue Title

## Description

## Acceptance Criteria

## Notes
```

## PRs as a Request Surface

External PRs from contributors are **not** managed through this issue tracker. This is a solo/internal project workflow.