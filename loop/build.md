# Build Report — Child Completion Gate

## Gap
268/478 done tasks (56%) had incomplete children — board integrity was unreliable because the "complete" op had no enforcement.

## Changes

### `site/graph/store.go`
- Added `ErrChildrenIncomplete` sentinel error
- Modified `UpdateNodeState`: when `state == StateDone`, queries for incomplete children (`state != 'done'`) before the UPDATE. Returns `ErrChildrenIncomplete` if any exist. Zero-child nodes pass through unchanged.

### `site/graph/handlers.go`
- `handleOp` case `"complete"`: checks `errors.Is(err, ErrChildrenIncomplete)` → 422 Unprocessable Entity
- `handleNodeState`: same check when `newState == StateDone` → 422 Unprocessable Entity

### `site/graph/store_test.go`
- Added `errors` import
- Added `TestUpdateNodeStateChildGate`: creates parent + child task, verifies parent completion is rejected with `ErrChildrenIncomplete` while child is incomplete, then verifies it succeeds after child is completed

## Verification
- `go.exe build -buildvcs=false ./...` — exit 0
- `go.exe test ./...` — all pass (graph, auth, handlers)

## Design decision
Enforcement is in `UpdateNodeState` (store layer), not the handler layer — single enforcement point means any future caller (API, CLI, agent) gets the gate automatically.
