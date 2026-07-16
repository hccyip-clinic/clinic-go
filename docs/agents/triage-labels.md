# Triage Labels

This repository uses the following canonical triage labels:

| Label | Description | When Applied |
|-------|-------------|--------------|
| `needs-triage` | New issue awaiting initial review | Default state for new issues |
| `needs-info` | Blocked, awaiting user/stakeholder input | Agent asked a question, waiting for answer |
| `ready-for-agent` | Specified and ready for implementation | All questions answered, clear acceptance criteria |
| `ready-for-human` | Work complete, awaiting human review | Agent finished, needs PR review or validation |
| `wontfix` | Issue will not be addressed | Deprioritized, obsolete, or rejected |

## State Machine

```
needs-triage → needs-info ↔ ready-for-agent → ready-for-human → (merged/closed)
     ↓
  wontfix
```

## Usage Rules

1. Every issue starts as `needs-triage`
2. Agent moves to `needs-info` if clarification needed
3. User moves to `ready-for-agent` when ready
4. Agent moves to `ready-for-human` when work is done
5. Human reviews and closes or requests changes