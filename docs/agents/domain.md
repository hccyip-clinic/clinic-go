# Domain Documentation

This repository uses a **single-context** layout for domain documentation.

## Structure

```
clinic-hugo/
├── CONTEXT.md          # Main domain context (to be created)
└── docs/
    └── adr/            # Architecture Decision Records
        └── NNNN-*.md
```

## File Responsibilities

### CONTEXT.md

The single source of truth for:
- Business domain overview
- Key user roles and personas
- Core workflows and use cases
- Domain terminology (ubiquitous language)
- System boundaries and integrations

### docs/adr/

Architecture Decision Records documenting:
- Significant architectural choices
- Context and consequences of decisions
- Status (proposed, accepted, deprecated, superseded)

## Consumer Rules for Agents

When working in this codebase, agents MUST:

1. **Read CONTEXT.md first** before implementing features or making architectural decisions
2. **Check docs/adr/** for existing decisions before proposing new architecture
3. **Update CONTEXT.md** when discovering new domain concepts during implementation
4. **Create new ADRs** in `docs/adr/` for significant architectural changes
5. **Use ubiquitous language** from CONTEXT.md in all communications and code comments

## ADR Template

Use the existing ADR format in `docs/adr/0004-hugo-htmx-sqlite-architecture.md` as a template for new ADRs.