# Build: Causality fix is narrow: Observer-created nodes still have causes=[] after commits 274999c and 8a13ac7

- **Timestamp:** 2026-03-28T00:00:00Z

## Task

Causality fix tasks (1a90716f, 2f8fd5a0, ee651efd) were marked done, but 3 Observer-created board nodes still had causes=[]. The prior fix only covered the Architect Operate path (buildArchitectOperateInstruction). All other creation paths remained broken.

## What Was Built

Audited ALL `intend` op creation call sites and fixed each to declare causes:

### 1. `cmd/post/main.go` — `createTask()`
- Added `causeIDs []string` parameter
- Passes `causes` field in the JSON payload when non-empty
- Caller now passes `buildDocID` so every board task is linked to the build document that triggered it
- Added `TestCreateTaskSendsCauses` test to pin the invariant

### 2. `pkg/runner/observer.go` — Observer Operate path
- Updated `buildOutputInstruction()` curl template to include `"causes":["<NODE_ID>"]`
- Added instruction: "Replace <NODE_ID> with the ID of the specific board node, claim, or document that triggered this finding (Invariant 2: CAUSALITY)"

### 3. `pkg/runner/observer.go` — Observer Reason path
- Added `causeID string` field to `observerTask` struct
- Updated `parseObserverTasks()` to parse `TASK_CAUSE:` lines from LLM output
- Updated prompt to request `TASK_CAUSE: <node_id>` for each finding
- Updated `runObserverReason()` to pass the parsed cause ID to `CreateTask`

### 4. `pkg/runner/pm.go` — PM Operate path
- Updated `runPMOperate()` curl template to include `"causes":["<PINNED_GOAL_NODE_ID>"]`
- Instructs PM to use the pinned goal node it reads in step 1 as the cause

### 5. `pkg/runner/pm.go` — PM Reason path
- `runPMReason()` now looks up the most recent `Build:` document and uses its ID as cause when creating the milestone

### 6. `pkg/runner/critic.go` — Critic Operate path
- Before calling Operate, looks up the `Build: {subject}` document for the commit being reviewed
- Injects the build node ID into the Fix task curl template as `"causes":["{buildNodeID}"]`

### 7. `pkg/runner/critic.go` — `writeCritiqueArtifact()`
- Added `causeIDs []string` parameter
- Passes causes to `AssertClaim()` so critique claims are linked to the build they review

### 8. `pkg/runner/reflector.go` — Reflector Operate path
- Looks up current `Critique:` and `Build:` node IDs before calling Operate
- Injects them as `causes` in both the `assert` (lesson) and `intend` (document) curl templates

### 9. `pkg/runner/reflector.go` — Reflector Reason path
- Added `readFromGraphNode()` helper (returns full `*api.Node` including ID)
- Refactored `readFromGraph()` to delegate to `readFromGraphNode()`
- `runReflectorReason()` now collects `iterationCauses` from critique/build node IDs
- `AssertClaim()` (lessons) and `CreateDocument()` (reflection) now pass `iterationCauses`
- `appendReflection()` accepts `causeIDs []string` and forwards to `CreateDocument`

## Coverage
- `go.exe build -buildvcs=false ./...` — passes
- `go.exe test ./...` — all 13 packages pass
- New test: `TestCreateTaskSendsCauses` in `cmd/post/main_test.go`
