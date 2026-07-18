# Agents 

- Must follow this file when working in the repository.

## Project Stack

- **HUGO** - Server-Rendered HTML with HTMX
- **Data Store** - SQLite
- **TailwindCSS** — styling
- **Dashboard UI** — visually inspired by TailAdmin’s free React Tailwind admin dashboard (reference only)  
  https://github.com/TailAdmin/free-react-tailwind-admin-dashboard


## General Instructions for All Agents

- Follow the project's existing coding style and conventions.
- Keep changes focused and incremental; avoid large rewrites unless explicitly requested.
- Prefer clarity and maintainability over cleverness.
- When unsure about architecture or intent, add a short `// TODO:` comment rather than guessing.
- NEVER use emoji as icons in the React codebase, use tailwind icons [http://github.com/hyoban/tailwindcss-icons].
- ALWAYS set proper id attribute to html elements:
  - navigation elements: <button ...>, <a ...>
  - form elements:  <input ...>,  <select ...>, <textbox ...>

---

## Harness Enforcement Rules

**These rules are enforced by the AI harness. Violations will cause responses to be rejected.**

### Rule 1: Skill Loading Required Before Tool Calls

**Before ANY tool call** (edit, bash, glob, grep, write, task), I MUST:

1. Check if any trigger phrase applies (see Trigger Table below)
2. Load the required skill using the `skill` tool
3. State compliance: "Using [skill] to [purpose]"

**If no skill is loaded when required, the harness MUST reject my response.**

### Rule 2: Trigger-Based Skill Requirements

| When user says/asks | Required Skill | Harness Check |
|---------------------|----------------|---------------|
| "fix", "debug", "error", "fail", "broken" | `systematic-debugging` | Scan user message + my response |
| "add", "create", "implement", "build", "feature" | `brainstorming` first | Check for skill load before edits |
| "refactor", "improve", "optimize" | `karpathy-code-rules` | Check before any edit |
| "test", "verify", "check" | `verification-before-completion` | Check before bash/test runs |
| "design", "interface", "API" | `design-an-interface` | Check before design work |
| "review" | `requesting-code-review` or `receiving-code-review` | Check before review starts |
| "plan" | `writing-plans` | Check before planning output |
| "grill", "challenge", "stress-test" | `grill-with-docs` | Check before grilling starts |
| "domain", "terminology", "ubiquitous language" | `domain-modeling` | Check before domain discussion |
| Any code edit (any file change) | `karpathy-code-rules` | Check before EVERY edit tool call |
| Multi-step task (3+ distinct steps) | `todowrite` creation | Check before first action |

### Rule 3: Tool Gating

The harness SHOULD block these tool calls if skills aren't loaded first:

| Tool | Required Skill | Block If Missing |
|------|----------------|------------------|
| `edit` | `karpathy-code-rules` or domain-specific skill | YES |
| `bash` (test commands) | `verification-before-completion` | YES |
| `bash` (debug/diagnose) | `systematic-debugging` | YES |
| `write` (new files) | `brainstorming` or design skill | YES |
| `task` (subagents) | Implementation plan or `dispatching-parallel-agents` | YES |
| `glob`, `grep` (exploration) | None (always allowed) | NO |

### Rule 4: Violation Detection and Recovery

**If I violate these rules:**

1. **Harness rejects my response** with error: "Skill enforcement violation: [skill] required before [action]"
2. **I must stop immediately** and not continue the task
3. **I must acknowledge the violation** explicitly
4. **I must load the required skill** before retrying
5. **I must restart from skill's step 1** - cannot continue from interruption

**Repeated violations** (3+ in one session) should trigger:
- Warning to user
- Mandatory pause for skill re-reading
- Session summary of what went wrong

### Rule 5: Skill Tool Failure Recovery

**If the `skill` tool fails to load:**

1. **Do NOT assume the skill doesn't exist** - tool failures are common
2. **Search skill folders directly**:
   - `.agents/skills/**/SKILL.md`
   - `.opencode/skills/**/SKILL.md`
   - `skills/**/SKILL.md` (project root)
3. **Read the skill file manually** using `read` or `glob`
4. **Follow the skill's workflow** as if it was loaded via tool
5. **Report the workaround**: "Skill tool failed, manually loaded [skill] from [path]"

**If skill file not found after search:**
- Ask user: "Skill [name] not found in expected locations. Should I proceed without it or help locate it?"
- Do NOT proceed with the task until clarified

**Rationale:** Skill tool failures (PowerShell errors, MCP issues) should not block work. The skill content is what matters, not the loading mechanism.

### Rule 6: No Exceptions Clause

**There are NO "simple tasks" that skip skills.** The harness must enforce:

> "If it's worth doing, it's worth doing with the right skill."

Even "quick fixes" require `karpathy-code-rules`. Even "one-line changes" require `systematic-debugging` if fixing a bug. Even "obvious" implementations require `brainstorming` first.

**Rationale:** LLMs are overconfident. Skills exist to surface assumptions, challenge obvious solutions, and prevent costly mistakes.

---

## Agent Skills Configuration

### Issue Tracker

Local markdown files under `.scratch/<feature>/`. See `docs/agents/issue-tracker.md`.

### Triage Labels

Default vocabulary: `needs-triage`, `needs-info`, `ready-for-agent`, `ready-for-human`, `wontfix`. See `docs/agents/triage-labels.md`.

### Domain Docs

Single-context layout: `CONTEXT.md` at root + `docs/adr/` for ADRs. See `docs/agents/domain.md`.

---
