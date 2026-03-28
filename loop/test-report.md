# Test Report: CAUSALITY Integration Tests (Invariant 2)

- **Iteration:** 404
- **Timestamp:** 2026-03-29
- **Build:** d1ec670 — Add integration test: every node creation code path must have non-empty causes

## What Was Tested

Builder added `pkg/loop/causality_test.go` with 3 tests pinning Invariant 2 (CAUSALITY) to CI for three node-creation code paths. Tester added 1 edge-case test.

### Loop package (`pkg/loop/causality_test.go`)

| Test | Who | Covers |
|---|---|---|
| `TestCausality_LoopTaskCommandPath` | Builder | `/task create` via `processTaskCommands` → task event has `len(Causes()) > 0` |
| `TestCausality_DirectAPICallPath` | Builder | `api.Client.CreateTask` → HTTP body contains `causes` array |
| `TestCausality_CmdPostPath` | Builder | `api.Client.AssertClaim` → HTTP body contains `causes` array |
| `TestCausality_LoopTaskCommandPath_MultipleTasks` | Tester | 3× `/task create` in one response → all 3 task events have causes (regression guard on shared-causes-slice refactor risk) |

## Results

```
TestCausality_LoopTaskCommandPath                  PASS
TestCausality_DirectAPICallPath                    PASS
TestCausality_LoopTaskCommandPath_MultipleTasks    PASS
TestCausality_CmdPostPath                          PASS

pkg/loop    PASS  4/4 causality tests
Full suite  PASS  all packages
```

**Total: 4/4 PASS. Full suite green.**

## Coverage Notes

- Loop `/task create` path: `processTaskCommands` → `executeTaskCommands(causes=[LastEvent])` confirmed non-empty for single and multiple commands
- Direct API path: `CreateTask(causes)` correctly serializes causes into HTTP JSON body
- cmd/post path: `AssertClaim(causes)` correctly serializes causes into HTTP JSON body
- `CreateDocument` with causes is covered separately in `pkg/api/client_test.go` (`TestCreateDocumentSendsCauses`)
- Known gap NOT addressed by this build: **Observer `Reason` path** — Observer creates nodes with empty causes (task c2ab9f11, state.md item 2). No test covers this because the production code is still broken; a test would need the Observer fix first.

## @Critic

Tests done. Ready for review.
