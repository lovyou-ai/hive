# Build: Zero causes links: graph is causally disconnected — 0/486 nodes have causes declared

- **Commit:** (pending — Ops will commit)
- **Cost:** $0.0000
- **Timestamp:** 2026-03-28T00:00:00Z

## Task

Every node on the hive board has zero causes links (0/486). This violates Invariant 2 (CAUSALITY): "Every event has declared causes."

## What Was Built

### Root cause
`CreateTask`, `CreateDocument`, and `AssertClaim` in `api/client.go` took `map[string]string` — no way to include `causes: []string` (a JSON array).

### Changes

**`pkg/api/client.go`**
- Added `postOpAny(slug string, fields map[string]any)` — underlying implementation supporting array fields
- `PostOp` (public, `map[string]string`) delegates to `postOpAny` for backward compat
- `CreateTask(slug, title, description, priority string, causes []string)` — new `causes` parameter; sends `"causes":[...]` when non-empty
- `CreateDocument(slug, title, body string, causes []string)` — same
- `AssertClaim(slug, title, body string, causes []string)` — same

**`pkg/runner/critic.go`**
- `writeCritiqueArtifact` now returns `(string, error)` — string is the claim node ID
- `reviewCommit` threads that ID into fix task `causes: [claimID]` — fix tasks traceable to their critique

**`pkg/runner/runner.go`**
- `writeBuildArtifact` passes `[]string{t.ID}` to `CreateDocument` — build docs causally linked to task

**`pkg/runner/architect.go`**
- Subtasks pass `causes: [milestoneID]` — decomposition traceable to milestone

**Other callers** (observer, pm, reflector, pipeline_tree)
- Updated to pass `nil` causes (no triggering node in those paths)

### Tests added

- `pkg/api/client_test.go` — 5 tests: `CreateTask`/`CreateDocument`/`AssertClaim` send causes; nil causes omit field; `PostOp` string fields preserved
- `pkg/runner/critic_test.go` — `TestReviewCommitFixTaskHasCauses`: verifies REVISE verdict creates fix task with `causes:["<critique-claim-id>"]`

All tests pass. Build passes.
