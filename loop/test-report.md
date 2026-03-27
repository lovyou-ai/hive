# Test Report: Observer reads /board for claim audit

- **Date:** 2026-03-28
- **Build commit:** 165d4f0

## What Was Tested

The build added a second curl to `/knowledge?tab=claims&limit=50` in `buildPart2Instruction` so the Observer sees existing claims instead of always reporting zero.

### New tests added in `pkg/runner/observer_test.go`

#### `TestBuildPart2InstructionBoardAndClaims`
Verifies that when an API key is set, the instruction includes **both** the `/board` URL and the `knowledge?tab=claims` URL (with `limit=50`), and that exactly 2 `Authorization: Bearer` headers appear (one per curl command).

#### `TestParseObserverTasks` (10 sub-tests)
`parseObserverTasks` had zero tests. Full coverage added:

| Sub-test | What it checks |
|----------|---------------|
| empty input | Returns nil, not a zero-value task |
| no issue found text | Unrecognised text produces no tasks |
| single complete task | Title, priority, description parsed correctly |
| two tasks | Second task flushed when third TASK_TITLE: encountered |
| invalid priority → medium | "critical" normalised to "medium" |
| missing priority → medium | Zero-value priority normalised to "medium" |
| all valid priorities accepted | urgent/high/medium/low pass through unchanged |
| whitespace trimmed | Leading/trailing spaces on all fields stripped |
| title only (no other fields) | Task emitted with empty desc and defaulted priority |
| unrecognised lines ignored | Noise between directives doesn't corrupt parse |

## Results

```
ok  github.com/lovyou-ai/hive/pkg/runner  3.661s
```

All 11 new tests pass. No regressions across the full runner package (28 total tests).

## Coverage Notes

- `parseObserverTasks` was completely untested; now covered across all branches.
- The "exactly 2 curls" assertion catches future regressions where the claims fetch is dropped from the instruction.
- `buildObserverInstruction`, `buildOutputInstruction` were already covered; no changes needed.
