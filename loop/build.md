# Build: Prevent Fix: title compounding — dedup before creating fix task

## What Was Built

Added deduplication to `createTask` in `cmd/post/main.go` to prevent "Fix: Fix: Fix: X" title compounding on the board.

## Root Cause

`createTask` always called `op=intend` unconditionally. When a build failed Critic review and the next iteration's `build.md` title was "Fix: [original title]", a new task was created. Repeated failures compounded: "Fix: Fix: X", "Fix: Fix: Fix: X", etc. This produced 95 compound-title tasks.

## Files Changed

### `cmd/post/main.go`

Added three new functions:

- **`stripFixPrefixes(title string) string`** — strips all leading "Fix: " prefixes to get the core title. "Fix: Fix: X" → "X".
- **`findExistingTask(apiKey, baseURL, coreTitle string) (string, error)`** — searches the board for a task whose title (after stripping Fix: prefixes) matches `coreTitle`. Uses existing `fetchBoardByQuery` with client-side title comparison.
- **`addTaskComment(apiKey, baseURL, nodeID, body string) error`** — posts `op=respond` on an existing task node (adds a comment/follow-up instead of a new task).

Modified **`createTask`** to run the dedup check before creating:
1. Strip all "Fix: " prefixes from `title` to get `coreTitle`
2. If `coreTitle != title` (title had at least one Fix: prefix): query board for existing task
3. If found: call `addTaskComment`, return existing task ID — no new task created
4. If not found or board API fails: fall through to normal `op=intend` creation (dedup is best-effort)

Added **`upgradeTaskPriority(apiKey, baseURL, nodeID, priority string) error`** — sends `op=edit` with a new priority value.

In **`main()`**: added a call to upgrade task `468e0549` from low to high (covers 95 nodes, not the 26 originally estimated).

### `cmd/post/main_test.go`

Added 8 new tests:
- `TestStripFixPrefixes` — 8 cases including double/triple Fix: prefix, no prefix, case-sensitivity
- `TestAddTaskCommentSendsRespondOp` — verifies op=respond, correct node_id and body
- `TestAddTaskCommentAPIError` — API error propagated
- `TestFindExistingTaskMatchesCoreTitle` — finds task with matching stripped title
- `TestFindExistingTaskNoMatch` — returns empty when no title match
- `TestCreateTaskDeduplicatesFixTask` — double Fix: prefix triggers comment on existing task, no intend op
- `TestCreateTaskNoDedup` — non-Fix title always creates new task, board not queried
- `TestCreateTaskDeduplicatesBoardAPIError` — board API failure is non-fatal, falls through to normal creation

## Verification

```
go.exe build -buildvcs=false ./...   ✓ no errors
go.exe test ./...                     ✓ all pass (13 packages)
```
