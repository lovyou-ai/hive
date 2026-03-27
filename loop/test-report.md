# Test Report: IsError fix in claude_cli.go

**Date:** 2026-03-28

## What Was Tested

Builder fixed `Operate()` in `eventgraph/go/pkg/intelligence/claude_cli.go` to return an error when Claude emits `is_error:true` JSON alongside a non-zero exit code. Without the fix, `Operate()` would silently succeed in this case.

## Tests Run

### Pre-existing test (Builder wrote this)
**`TestOperateIsErrorReturnsError`** — `pkg/intelligence/claude_cli_test.go`
- Simulates: claude exits with code 1, stdout contains `{"is_error":true,"result":"task failed: permission denied"}`
- Verifies: `Operate()` returns an error containing the result message
- Result: **PASS** ✓

### Added by Tester (gap coverage)
**`TestOperateIsErrorZeroExitReturnsError`** — `pkg/intelligence/claude_cli_test.go`
- Simulates: claude exits with code 0 but stdout contains `{"is_error":true,"result":"tool call rejected"}`
- Verifies: `Operate()` returns an error even when exit code is 0 (line 263 in claude_cli.go)
- Result: **PASS** ✓

### Full suite
All 38 packages pass. No regressions.

```
ok  github.com/lovyou-ai/eventgraph/go/pkg/intelligence  5.1s
ok  [all 37 other packages]
```

## Coverage Notes

Four `IsError` check sites in `claude_cli.go`:

| Site | Path | Tested? |
|------|------|---------|
| `Operate()` non-zero exit (line 249) | the fix | ✓ `TestOperateIsErrorReturnsError` |
| `Operate()` zero exit (line 263) | pre-existing | ✓ `TestOperateIsErrorZeroExitReturnsError` (added) |
| `Reason()` zero exit (line 172) | pre-existing | ✗ not tested |
| `Reason()` non-zero exit | **MISSING** | — no `IsError` check exists here |

## Gap Flagged: Reason() non-zero exit has no IsError check

In `Reason()` lines 155-165, if claude exits non-zero but still writes `is_error:true` JSON to stdout, the code falls through to `resultToResponse()` and returns success. This is inconsistent with `Operate()`. Not a regression from this iteration, but worth fixing.

## Verdict

**PASS.** The Builder's fix works. Added one test to cover the zero-exit `is_error:true` path in `Operate()`. One gap flagged in `Reason()` non-zero exit path for a future iteration.
