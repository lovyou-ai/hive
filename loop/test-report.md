# Test Report: Iteration 367 — Knowledge API causes field

**Tester:** Tester
**Date:** 2026-03-28
**Verdict:** PASS

## What Was Tested

The Builder fixed `syncClaims` in `cmd/post/main.go` to include the `Causes` field when decoding the `/knowledge` API response and writing `loop/claims.md`.

## Tests Run

```
go.exe test ./cmd/post/
ok  github.com/lovyou-ai/hive/cmd/post  0.605s
```

32 tests, all pass.

## New Test Added

**`TestSyncClaimsMultipleCauses`** — verifies that when a claim has multiple cause IDs, all of them are written as a comma-joined `**Causes:**` line in `claims.md`. This exercises the `strings.Join(c.Causes, ", ")` path that the Builder introduced. The Builder's `TestSyncClaimsWritesCauses` only covers the single-cause case; this covers the multi-cause case.

## Coverage Notes

The key fix (adding `Causes []string` to the decode struct) is pinned by `TestSyncClaimsWritesCauses`. Removal of the field would cause that test to fail.

All relevant paths covered:
- Single cause: `TestSyncClaimsWritesCauses`
- Multiple causes: `TestSyncClaimsMultipleCauses` (added this iteration)
- No causes: `TestSyncClaimsClaimWithNoMetadata` / `TestSyncClaimsWritesFile` (no `**Causes:**` line emitted)
- Empty response: `TestSyncClaimsEmptyDoesNotWrite`
- API error: `TestSyncClaimsAPIError`
