# Build: Fix: populateFormFromJSON deploy + observer fallback cause

- **Iteration:** 404
- **Date:** 2026-03-29
- **Gap:** CAUSALITY invariant violated in production — array causes silently rejected; Observer Reason path emitting causeless tasks

## What Was Built

Three items completed from Scout scope:

### 1. Site Deploy — populateFormFromJSON fix

Deployed `site/graph/handlers.go` `populateFormFromJSON` fix to production via `flyctl deploy --remote-only`.

**Verification:** POST to `https://lovyou.ai/app/hive/op` with `"causes":["test-verify-404"]` returned node with `"causes":["test-verify-404"]` confirmed. Previously this returned "unknown op". Production now accepts JSON array causes.

### 2. Observer fallback cause — `pkg/runner/observer.go`

**Root cause:** When LLM emits `TASK_CAUSE: none`, `parseObserverTasks` correctly filters it to empty string. But `runObserverReason` then called `CreateTask` with nil causes — violating Invariant 2 (CAUSALITY).

**Fix:**
- `runObserver`: extract `fallbackCauseID = claims[0].ID` from pre-fetched claims; pass alongside `claimsSummary` to `runObserverReason`
- `runObserverReason` signature: `(ctx, claimsSummary, fallbackCauseID string)`
- Task loop: when `t.causeID == ""`, apply `fallbackCauseID` before building the causes slice

This closes task c2ab9f11. Observer Reason path no longer emits causeless nodes when claims are available.

### 3. Test — `pkg/runner/observer_test.go`

Added `TestRunObserverReason_FallbackCause`:
- Mock provider returns `TASK_CAUSE: none`
- Test server captures the CreateTask HTTP request body
- Asserts `"causes":["claim-fallback-123"]` is present — fallback is applied

## Changes

**`hive/pkg/runner/observer.go`**
- `runObserver`: declare `fallbackCauseID`; set from `claims[0].ID` when claims are available
- `runObserver`: pass `fallbackCauseID` to `runObserverReason` fallback call
- `runObserverReason`: add `fallbackCauseID string` parameter
- `runObserverReason` task loop: apply fallback when `t.causeID == ""`

**`hive/pkg/runner/observer_test.go`**
- Add imports: `context`, `encoding/json`, `io`, `net/http`, `net/http/httptest`
- Add `TestRunObserverReason_FallbackCause`

**`site/graph/handlers.go`** — deployed only, no code change needed

## Build

```
go.exe build -buildvcs=false ./pkg/runner/  → clean
go.exe test -buildvcs=false ./pkg/runner/   → PASS (all tests)
flyctl deploy --remote-only                 → deployed, 2 machines healthy
```

## Invariants

- **CAUSALITY (2):** Observer Reason path now always includes a cause (fallback to first claim when LLM emits none). Production accepts JSON array causes.
- **VERIFIED (12):** New test covers the fallback cause code path.
