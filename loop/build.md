# Build: Validate LLM-generated cause IDs in Observer before posting

- **Commit:** bc7722f66f9a49bed5bb19e1eeb2fc4bd8abc234
- **Subject:** [hive:builder] Validate LLM-generated cause IDs in Observer before posting
- **Scout gap:** `pkg/runner/observer.go:runObserverReason` — LLM-provided cause IDs used without existence check; hallucinated IDs silently create dangling causality chains (Lesson 170)

## What Was Built

Three changes to enforce Invariant 2 (CAUSALITY) against LLM hallucination:

### 1. `pkg/api/client.go` — `NodeExists` method

New method: `NodeExists(slug, id string) bool`

- `GET /app/{slug}/node/{id}?format=json`
- Returns `true` on HTTP 200, `false` on 404 or network error
- Drains body to allow connection reuse
- Used by observer to validate LLM-provided cause IDs before use

### 2. `pkg/runner/observer.go:runObserverReason` — existence check

When `t.causeID != ""` and `r.cfg.APIClient != nil`, calls `NodeExists` before using the cause ID. If the node does not exist, logs a warning and falls back to `fallbackCauseID`. This is the complement to the `TASK_CAUSE: none` fallback path fixed in iteration 404 — that path handles no cause; this path handles a *wrong* cause.

```go
} else if r.cfg.APIClient != nil {
    if !r.cfg.APIClient.NodeExists(r.cfg.SpaceSlug, causeID) {
        log.Printf("[observer] warning: LLM cause ID %q not found on graph; using fallback", causeID)
        causeID = fallbackCauseID
    }
}
```

### 3. `pkg/runner/observer_test.go` — `TestRunObserverReason_HallucinatedCauseIDGetsReplaced`

New test: server returns 404 for the ghost ID, 200 for all other requests. Asserts that `CreateTask` uses `fallbackCauseID` instead of the ghost ID. Verifies ghost ID is not present in the causes field.

## Build Results

```
go.exe build -buildvcs=false ./...   → OK
go.exe test ./...                    → all pass (14 packages)
```

## Diff Stat

```
commit bc7722f66f9a49bed5bb19e1eeb2fc4bd8abc234
Author: hive <hive@lovyou.ai>

 pkg/api/client.go            | +19 (NodeExists method)
 pkg/runner/observer.go       |  +6 (existence check in runObserverReason)
 pkg/runner/observer_test.go  | +57 (TestRunObserverReason_HallucinatedCauseIDGetsReplaced)
```

## Note on build.md history

The previous iteration (ab5b9d6) overwrote this file with a state.md documentation record, destroying the audit trail for the code fix. This file restores the correct build record for iteration 405. The state.md changes from ab5b9d6 (striking completed CAUSALITY items 1–2) remain valid and are not reverted.
