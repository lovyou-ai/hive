# Test Report: Zero causes links — Invariant 2 causality chain

**Build ref:** 8a13ac7f2f2fc395041446e39413b2cb3ea9ce09
**Tester:** hive:tester
**Timestamp:** 2026-03-28

## What Was Tested

The iteration added `causes []string` to `CreateTask`, `CreateDocument`, and `AssertClaim`, then wired causality through the critic→fix-task and build-document→task chains.

## Tests Already Present (Builder-written)

| Test | File | Verdict |
|------|------|---------|
| `TestCreateTaskSendsCauses` | `pkg/api/client_test.go` | PASS |
| `TestCreateTaskNilCausesOmitted` | `pkg/api/client_test.go` | PASS |
| `TestCreateDocumentSendsCauses` | `pkg/api/client_test.go` | PASS |
| `TestAssertClaimSendsCauses` | `pkg/api/client_test.go` | PASS |
| `TestPostOpStringFieldsPreserved` | `pkg/api/client_test.go` | PASS |
| `TestReviewCommitFixTaskHasCauses` | `pkg/runner/critic_test.go` | PASS |

## Coverage Gaps Found

Two paths were untested: `writeBuildArtifact` calling `CreateDocument` with causes, and `runArchitect` creating subtasks with causes when a milestone exists. Both are Invariant 2 violations if wrong.

## Tests Added

### `TestWriteBuildArtifactDocumentCauses` — `pkg/runner/runner_test.go`

Verifies that `writeBuildArtifact` calls `CreateDocument` with `causes: [task.ID]`. Previously the three existing build artifact tests used `Config` with no `APIClient`, so the `CreateDocument` call was silently skipped. This test wires a mock HTTP server and confirms the causality field is present and contains the triggering task ID.

**Result:** PASS

### `TestRunArchitectSubtasksHaveCauses` — `pkg/runner/architect_test.go`

Verifies that when the Architect decomposes a milestone into subtasks (via the `Reason()` fallback path), each created subtask carries `causes: [milestoneID]`. Previously only failure paths were tested (no parseable subtasks). This test provides a milestone with a >200-char body, a provider returning valid `SUBTASK_TITLE:` format, and verifies all `intend` ops include the milestone ID in `causes`.

**Result:** PASS

## Full Suite Results

```
ok  github.com/lovyou-ai/hive/pkg/api     (all 5 tests)
ok  github.com/lovyou-ai/hive/pkg/runner  (all tests, including 2 new)
```

## Coverage Notes

- `writeBuildArtifact` with `APIClient = nil` (causes skipped) is implicitly covered by existing tests — no gap there
- Architect `Operate()` path (canOperate=true) creates tasks via curl in the LLM — not unit-testable without a real operator; causality in that path relies on the Architect's prompt instruction
- Multiple-cause scenarios (len(causes) > 1) not tested; the mechanism is identical to single-cause and covered by the API-layer tests

## Verdict

**PASS** — causality chain is fully wired and verified for the three paths added in this iteration.

@Critic
