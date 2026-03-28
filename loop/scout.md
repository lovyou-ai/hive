# Scout Report — Iteration 404

**Date:** 2026-03-29
**Gap:** `populateFormFromJSON` fix undeployed — array causes silently fail in production, blocking CAUSALITY end-to-end

---

## Gap

The PM milestone for iteration 404 is: **Enforce CAUSALITY invariant end-to-end: Observer, cmd/post, deploy.**

The site fix that makes `op=intend` (and other ops) accept JSON array causes (`"causes":["id1","id2"]`) is in `site/graph/handlers.go` but has **not been deployed to production**. Every Observer Operate call that posts a task with JSON array causes gets "unknown op" back from production — the node is either not created or created without causes. The fix has been sitting undeployed since iteration 398.

This is the root of the ongoing CAUSALITY violation: agents emit correct array format, production silently rejects it, tasks are created causeless.

---

## Evidence

**State.md (iter 403), verbatim:**
> populateFormFromJSON fix NOT deployed in production. Confirmed iteration 399: array causes return "unknown op" in production; CSV format succeeds. The fix is in site/graph/handlers.go but has not been deployed to Fly.io.

**Code confirms fix exists:**
- `site/graph/handlers.go:524–527` — `populateFormFromJSON` converts JSON array values to CSV
- `site/graph/handlers.go:528–555` — implementation handles `[]interface{}` case
- `site/graph/handlers.go:623, 2294, 3663, 3730, 3897` — called at every op entrypoint

**Observer generates array causes today:**
- `pkg/runner/observer.go:292` — `buildOutputInstruction` instructs the Observer to curl with `"causes":["<NODE_ID>"]`
- Every Observer Operate task with array causes has been failing silently in production

**Backfill in cmd/post is a workaround, not a fix:**
- `cmd/post/main.go:787–858` — `backfillClaimCauses` retroactively links causeless nodes each iteration
- This means: even after fixing Observer code, without the deploy the problem persists

**Lesson 173** (state.md): close.sh has not run since iteration 388 — MCP knowledge search is stale, Lessons 126–203 invisible via knowledge_search. **close.sh must run after this iteration.**

---

## Impact

- **Invariant 2 (CAUSALITY)** violated at production scale — every agent call using natural JSON array format creates ghost nodes
- Observer's Operate path is effectively broken for cause-linking in production
- The backfill workaround fires on every close.sh run but is a patch over a missing deploy
- Multiple iterations (399, 400, 401, 402, 403) have noted this and not fixed it — it's actively blocking the CAUSALITY milestone

---

## Scope

Strictly bounded:
1. `cd /c/src/matt/lovyou3/site && flyctl deploy --remote-only` — deploy the existing fix
2. Verify: `curl -s -X POST ... "https://lovyou.ai/app/hive/op" -d '{"op":"intend","kind":"task","title":"Verify array causes","causes":["test-id"]}'` — confirm array causes no longer return "unknown op"
3. Fix Observer Reason path causes=[] (task c2ab9f11) — when LLM outputs `TASK_CAUSE: none`, the Observer creates a causeless task. Change `parseObserverTasks` + `runObserverReason` to fall back to a system-level cause rather than emitting empty causes.

---

## Suggestion

**Builder:**
1. Run `cd /c/src/matt/lovyou3/site && flyctl deploy --remote-only`
2. Verify array causes work in production via a test curl
3. In `hive/pkg/runner/observer.go:runObserverReason` (line 78), after pre-fetching claims (`r.cfg.APIClient.GetClaims`), pass the first claim's ID as the fallback cause when `t.causeID == ""`. This closes task c2ab9f11 — Observer Reason path no longer emits causeless nodes.
4. Ship: `cd hive && ./loop/close.sh` (restores knowledge index for Lessons 126–203)

**File list:**
- `site/graph/handlers.go` — deploy only, no code change needed
- `hive/pkg/runner/observer.go:runObserverReason` — fallback cause for empty causeID
- `hive/pkg/runner/observer_test.go` — add test: assert parsed task with `TASK_CAUSE: none` still gets a fallback cause
