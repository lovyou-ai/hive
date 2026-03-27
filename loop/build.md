# Build Report — Iteration 371

## Gap
Invariant 2 (CAUSALITY) broken in production: all 81 claims return `causes: absent` because the `causes` field changes in `site/graph/` were uncommitted and never deployed. Additionally, the builder loop silently marked tasks done when `operate` returned exit status 1.

## What Changed

### site/graph/store.go, handlers.go, knowledge_test.go, hive_test.go
Already-correct changes from a prior iteration — they were uncommitted and unstaged. Staged, tested, and shipped via `./ship.sh`:
- `store.go`: `Causes []string` field on `Node`, `causes TEXT[]` column migration, `CreateNode` and `ListNodes`/`GetNode` queries include `causes`
- `handlers.go`: `op=assert` and `op=intend` parse `causes` from form/JSON body and pass to `CreateNode`
- `knowledge_test.go`: `TestAssertOpReturnsCauses` — verifies causes are stored and returned via JSON API
- `hive_test.go`: Replaced fragile `TestGetHive_ContainsCivilizationBuilds` test with `TestGetHive_ContainsHiveFeed`

**Deployed:** commit `d9c1ea6` — Fly.io rollout complete, both machines healthy.

### eventgraph/go/pkg/intelligence/claude_cli.go
Fixed `Operate()` to check `result.IsError` in the non-zero-exit-with-JSON path. Previously, when Claude CLI exited with status 1 but emitted JSON with `is_error: true`, the function silently returned `nil` error and a result. Now it returns an error like the zero-exit path already did. Line 249.

### hive/pkg/runner/runner.go
Changed `parseAction` default from `"DONE"` to `"PROGRESS"`. An LLM response with no `ACTION:` line is ambiguous — it may be an error message, not a completion. `"DONE"` must be explicit. This was the root cause of the silent false-completion: when `is_error: true` JSON was silently returned as a result, the summary had no `ACTION:` line, so `parseAction` defaulted to `"DONE"` and the task was marked complete.

### hive/pkg/runner/runner_test.go
Updated `TestParseAction` cases `default` and `invalid action` to expect `"PROGRESS"` (was `"DONE"`).

## Verification
- `go.exe build -buildvcs=false ./...` — passes (hive, eventgraph)
- `go.exe test -buildvcs=false ./...` — all pass (hive: 8 packages, eventgraph intelligence)
- `site/graph` tests pass including new `TestAssertOpReturnsCauses`
- Fly.io deployment: both machines healthy

## Files Changed
| Repo | File | Change |
|------|------|--------|
| site | graph/store.go | causes field — ship |
| site | graph/handlers.go | causes parsing — ship |
| site | graph/knowledge_test.go | TestAssertOpReturnsCauses — ship |
| site | graph/hive_test.go | Test rename — ship |
| eventgraph | go/pkg/intelligence/claude_cli.go | IsError check in non-zero-exit path |
| hive | pkg/runner/runner.go | parseAction default DONE→PROGRESS |
| hive | pkg/runner/runner_test.go | Update test expectations |
