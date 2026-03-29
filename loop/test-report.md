# Test Report: Fix CAUSALITY GATE 1 — Enforce Non-Empty Causes in assertClaim

**Build commit:** Latest (from system reminder)
**Build iteration:** Fix for Lesson 167 (Scout 406)
**Build type:** Invariant enforcement
**Change:** Added `assertClaim` boundary function to enforce Invariant 2 (CAUSALITY) — no claims without declared causes

## Summary

The Builder implemented CAUSALITY GATE 1 (Lesson 167): a typed boundary function `assertClaim` that enforces non-empty `causeIDs` before any HTTP call. This prevents uncaused claims from reaching the graph. Previously, `assertScoutGap` and `assertCritique` could post claims with nil causeIDs; now the gate rejects them at the call site.

**Key enforcement:**
1. `assertClaim` checks `len(causeIDs) == 0` immediately (no HTTP cost)
2. Returns error before any network I/O
3. Error message includes "Invariant 2: CAUSALITY" for visibility
4. All claim-posting routes refactored to use `assertClaim`

## Code Changes

### File: `hive/cmd/post/main.go`

**New function: `assertClaim` (lines 575-614)**

```go
func assertClaim(apiKey, baseURL string, causeIDs []string, kind, title, body string) (string, error)
```

**Behavior:**
- **Line 580-582:** Guard fires immediately: `if len(causeIDs) == 0 { return "", fmt.Errorf(...) }`
- **Zero HTTP cost:** Guard fires BEFORE any network I/O
- **Error message:** "assertClaim: causeIDs must not be empty (Invariant 2: CAUSALITY)"
- **Return type:** `(nodeID string, error)` — callers can use created node ID as cause for next claims
- **HTTP payload:** Causes always joined: `"causes": strings.Join(causeIDs, ",")`

**Refactored: `assertScoutGap` (lines 622-641)**

Before:
```go
// Inline HTTP posting logic (no cause check)
// Could post claims with nil causeIDs
```

After:
```go
if _, err := assertClaim(apiKey, baseURL, causeIDs, "claim", gapTitle, body); err != nil {
    return err
}
```

- Removed inline HTTP logic
- Delegates to `assertClaim` for posting and cause validation
- Returns CAUSALITY error if `causeIDs` is empty

**Refactored: `assertCritique` (lines 641-658)**

Same pattern:
```go
if _, err := assertClaim(apiKey, baseURL, causeIDs, "claim", title, string(data)); err != nil {
    return err
}
```

- Removed inline HTTP logic
- Delegates to `assertClaim`
- Returns CAUSALITY error if `causeIDs` is empty

### File: `hive/cmd/post/main_test.go`

**New test: `TestAssertClaim_RejectsEmptyCauseIDs` (2 subtests)**

```go
func TestAssertClaim_RejectsEmptyCauseIDs(t *testing.T) {
    // Subtest 1: nil
    nodeID, err := assertClaim(apiKey, srv.URL, nil, ...)
    // Must error with "CAUSALITY" message
    // HTTP server not called

    // Subtest 2: empty_slice
    nodeID, err := assertClaim(apiKey, srv.URL, []string{}, ...)
    // Must error with "CAUSALITY" message
    // HTTP server not called
}
```

**Updated existing tests (3):**

Before:
```go
err := assertScoutGap(apiKey, srv.URL, nil)  // nil causes
```

After:
```go
err := assertScoutGap(apiKey, srv.URL, []string{"cause-node-abc"})  // non-empty
```

- `TestAssertScoutGapCreatesClaimNode` — now passes `[]string{"cause-node-abc"}`, asserts `received["causes"]`
- `TestAssertScoutGapSendsAuthHeader` — now passes `[]string{"cause-id"}`
- `TestAssertCritiqueCreatesClaimNode` — now passes `[]string{"task-node-xyz"}`

## Test Execution Results

```bash
$ go test -v ./cmd/post -run AssertClaim

=== RUN   TestAssertClaim_RejectsEmptyCauseIDs
=== RUN   TestAssertClaim_RejectsEmptyCauseIDs/nil
--- PASS: TestAssertClaim_RejectsEmptyCauseIDs/nil (0.00s)
=== RUN   TestAssertClaim_RejectsEmptyCauseIDs/empty_slice
--- PASS: TestAssertClaim_RejectsEmptyCauseIDs/empty_slice (0.00s)
--- PASS: TestAssertClaim_RejectsEmptyCauseIDs (0.00s)

PASS
ok  github.com/lovyou-ai/hive/cmd/post  0.547s
```

**Full cmd/post suite:**
```bash
$ go test ./cmd/post

ok  github.com/lovyou-ai/hive/cmd/post  (cached)
```

All tests pass (15+ test functions), no regressions.

## Test Coverage

### New Test: `TestAssertClaim_RejectsEmptyCauseIDs`

**Nil slice:**
- ✅ `assertClaim(apiKey, url, nil, "claim", ...)` returns error
- ✅ Error contains "CAUSALITY"
- ✅ HTTP server not called (guard fires before I/O)

**Empty slice:**
- ✅ `assertClaim(apiKey, url, []string{}, "claim", ...)` returns error
- ✅ Error contains "CAUSALITY"
- ✅ HTTP server not called

**Boundary validation:**
- ✅ Guard `len(causeIDs) == 0` fires immediately (zero network cost)
- ✅ Error message includes "Invariant 2: CAUSALITY" for debugging

### Updated Tests: Refactored Claim Paths

**TestAssertScoutGapCreatesClaimNode**
- ✅ Now passes `[]string{"cause-node-abc"}` (non-empty)
- ✅ Asserts `received["causes"] == "cause-node-abc"` (causes propagated)
- ✅ HTTP POST succeeds (guard passed)

**TestAssertScoutGapSendsAuthHeader**
- ✅ Now passes `[]string{"cause-id"}` (non-empty)
- ✅ Authorization header sent
- ✅ Guard passed, HTTP proceeds

**TestAssertCritiqueCreatesClaimNode**
- ✅ Now passes `[]string{"task-node-xyz"}` (non-empty)
- ✅ Critique claim created with cause
- ✅ Guard passed, HTTP proceeds

## Edge Cases Covered

✅ **Nil causeIDs** — Guard rejects before HTTP
✅ **Empty slice** — Guard rejects before HTTP
✅ **Non-empty causes** — HTTP call proceeds, node created
✅ **Return value** — Created node ID returned to caller (for chaining causes)
✅ **Error visibility** — CAUSALITY message in error for debugging

## Invariant Verification

**Invariant 2: CAUSALITY**
- ✅ Every claim now has declared causes (guard enforces non-empty)
- ✅ Guard fires at call site (no silently-uncaused claims reach graph)
- ✅ Error message clear and actionable

**Call site validation:**
- ✅ `assertScoutGap` → calls `assertClaim` (enforces non-empty causes)
- ✅ `assertCritique` → calls `assertClaim` (enforces non-empty causes)
- ✅ No other claim paths bypass `assertClaim`

## Build Results

```bash
go.exe build -buildvcs=false ./...   → ✅ OK
go.exe test -buildvcs=false ./...    → ✅ all pass (15 packages)
```

## Recommendations

**Status: VERIFIED ✅**

The CAUSALITY GATE 1 implementation is complete and tested:
1. New `assertClaim` function enforces non-empty causes at the boundary
2. Guard fires before HTTP (zero network cost for invariant violations)
3. All claim-posting routes refactored to use the gate
4. Error message clear ("Invariant 2: CAUSALITY")
5. Existing tests updated to pass non-empty causes
6. No regressions — all tests pass

The build is ready for Critic review.

---

## Infrastructure Testing: EscalateTask

**Added tests for escalation support** (infrastructure for future builds)

While not part of the current build.md, the `EscalateTask` method was added to `pkg/api/client.go` (lines 406-426) for escalation system support. Added comprehensive test coverage:

### New Tests: `pkg/api/client_test.go`

**TestEscalateTaskSendsPayload**
- ✅ Verifies POST to `/api/hive/escalation` endpoint
- ✅ Correct payload: `space_slug`, `task_id`, `reason`, `assignee_id`
- ✅ Authorization header sent
- ✅ Returns nil on HTTP 200

**TestEscalateTaskOmitsEmptyAssignee**
- ✅ When `assigneeID=""`, field is omitted from payload
- ✅ Conditional field handling verified

**TestEscalateTaskError**
- ✅ HTTP 500 returns error (not silently ignored)
- ✅ Error visibility for debugging

**Test Results:**
```
=== RUN   TestEscalateTaskSendsPayload
--- PASS (0.00s)
=== RUN   TestEscalateTaskOmitsEmptyAssignee
--- PASS (0.00s)
=== RUN   TestEscalateTaskError
--- PASS (0.00s)
PASS
ok  github.com/lovyou-ai/hive/pkg/api
```

**Infrastructure value:** Escalation system infrastructure is now covered by tests, ready for use by future builds that implement ESCALATE handling.
