# Build Report — Iteration 231: Critic Bug Fix Deployed

## What This Iteration Does

Closes the loop on the Critic's REVISE from iter 230: fixes the progress handler state machine bug and applies lesson 57 to the Critic (assign fix tasks).

## Files Changed

### Hive repo
| File | What |
|------|------|
| `pkg/runner/critic.go` | +6 lines: Critic now assigns fix tasks to the agent after creation (same pattern as Scout). |

### Site repo (deployed)
| File | What |
|------|------|
| `graph/handlers.go` | +4 lines: Added `node.State != StateActive` guard to progress handler. Tasks must be active to submit for review. |

## The Bug (caught by Critic in iter 230)

The `progress` handler moved tasks to review state without checking the current state. Any task — done, closed, open — could be moved to review, violating the state machine:

```
open → review ← WRONG (should be: active → review only)
done → review ← WRONG
```

Fix: `if node.State != StateActive { return 400 }`

## Full Bug Lifecycle

1. **Iter 229:** Builder autonomously shipped review/progress ops
2. **Iter 230:** Critic autonomously reviewed the commit, found missing state guard, returned REVISE, created fix task
3. **Iter 231:** Human fixed the bug, deployed, closed the fix task

The system works: Builder ships → Critic catches → fix deployed. Eventually: Builder ships → Critic catches → Builder fixes (fully autonomous).

## Build

- `go build ./...` ✓ (both repos)
- `go test ./...` ✓
- `flyctl deploy --remote-only` ✓ — state guard live on lovyou.ai
