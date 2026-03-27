# Build Report — Observer claim context injection

## Gap
Observer was calling `/board` (which returns only `kind=task` nodes) and concluding "zero claims" every iteration. There are 65 actual claims at `/knowledge?tab=claims`. The LLM, given only the board response, had no evidence claims existed and kept generating false fix tasks.

## Root Cause
`buildObserverInstruction` told the Observer to *fetch* claims via curl, but didn't inject the count as ground truth. The LLM could ignore or misread the curl result. For `runObserverReason` (fallback path), there was no claims context at all.

## Fix

### `pkg/runner/observer.go`
1. **Added `buildClaimsSummary(claims []api.Node) string`** — formats pre-fetched claims as a grounding string: `"N claims exist (and X more). Titles: "T1", "T2", ..."`. Returns empty if no claims.

2. **`buildPart2Instruction` now accepts `claimsSummary string`** — when non-empty, injects a "Ground truth (pre-fetched by runner — do not contradict)" block before the curl commands. The LLM sees the count before it evaluates any curl output.

3. **`buildObserverInstruction` now accepts `claimsSummary string`** — threads through to `buildPart2Instruction`.

4. **`runObserver` pre-fetches claims before building the instruction** — calls `r.cfg.APIClient.GetClaims(slug, 50)` when APIClient is set. Logs the count. Passes `claimsSummary` to both `buildObserverInstruction` and `runObserverReason`.

5. **`runObserverReason` accepts `claimsSummary string`** — includes a "Claims (ground truth, pre-fetched)" section in the Reason prompt when non-empty.

### `pkg/runner/observer_test.go`
- Updated all callers of `buildPart2Instruction` and `buildObserverInstruction` to pass the new `claimsSummary` parameter.
- Added test cases in `TestBuildPart2Instruction` for ground-truth injection (with and without API key).
- Added `TestBuildClaimsSummary` covering: empty input, single claim, 5 claims (all shown), 6 claims (remainder count), 10 claims.

## Verification
- `go.exe build -buildvcs=false ./...` — clean
- `go.exe test ./...` — all pass
