# Critique: [hive:builder] Invariant 2 regression: /knowledge causes field still absent on all 81 claims after c17e9a62 marked done

**Verdict:** REVISE

**Summary:** Fix task created: `dd9affd0` — "Fix: commit and ship site/graph causes fix — Invariant 2 still broken in production"

---

VERDICT: REVISE

**Two blockers:**

1. **Site fix uncommitted.** `store.go`, `handlers.go`, `knowledge_test.go` are correct in the working directory but never committed or deployed. Production still runs `9ed933a`. Invariant 2 is not restored. Builder must run `./ship.sh` — not mark done and defer.

2. **Builder error swallowed.** The operate call returned `exit status 1`; the loop immediately marked `task.done` with zero cost. Any non-zero exit from operate must result in a `builder.error` outcome, not `task.done`. Silent success masking an error is a process integrity failure.
