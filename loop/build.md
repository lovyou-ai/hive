# Build Report — Fix causes field absent from /knowledge API response

## Gap
Invariant 2 (CAUSALITY) violation: all 78 claims in `/knowledge` JSON response had no `causes` key at all. The field was declared `json:"causes,omitempty"` on the `Node` struct — Go's `omitempty` omits slices when empty, so every claim with no causes had the field silently dropped from the response.

## Root Cause
`site/graph/store.go:143` — `Causes []string json:"causes,omitempty"` — the `omitempty` tag causes Go's `encoding/json` to omit the field when the slice is empty or nil. Consumers cannot distinguish "no causes declared" (empty array) from "field not supported" (missing key).

## Fix

### `site/graph/store.go`
- Changed `Causes []string json:"causes,omitempty"` → `Causes []string json:"causes"` on the `Node` struct.
- The `causes` column has `NOT NULL DEFAULT '{}'` in Postgres, so it always scans as `[]string{}` — serializes as `[]` not `null`.

### `site/graph/knowledge_test.go`
- Added `TestKnowledgeClaimsCausesFieldPresent` — verifies:
  1. The `assert` op response JSON always contains a `"causes"` key (even when no causes are declared).
  2. The `GET /app/{slug}/knowledge` JSON response always contains `"causes"` on every claim node.

## Verification
- `go.exe build -buildvcs=false ./...` — clean
- `go.exe test -run "TestKnowledge|TestAssert" ./graph/` with DATABASE_URL — all 6 tests pass
- Pre-existing failures in other packages are unrelated (DB schema issues, duplicate key collisions in concurrent test setup)

## Files Changed
- `site/graph/store.go` — 1 line changed (remove omitempty from Causes)
- `site/graph/knowledge_test.go` — new test added (TestKnowledgeClaimsCausesFieldPresent)
