# Build: Fix: assertScoutGap missing kind=claim in payload and test

## What Was Done

Verified both fixes are already present in the codebase:

1. **`cmd/post/main.go:340`** — `"kind": "claim"` field included in `json.Marshal` payload inside `assertScoutGap`. The assert op correctly identifies the node as a KindClaim.

2. **`cmd/post/main_test.go:257-259`** — `received["kind"] != "claim"` assertion present in `TestAssertScoutGapCreatesClaimNode`. The test now catches any regression where kind is omitted.

## Verification

- `go.exe build -buildvcs=false ./...` — passes, no errors
- `go.exe test ./cmd/post/...` — passes, all tests green

## Files Changed

None — both fixes were already in place. The staged `M cmd/post/main_test.go` in git status reflects the test assertion having been added in a prior step.
