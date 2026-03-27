# Build: Fix Observer meta-task defect — act inline instead of deferring

## Gap
Observer was creating tasks to close other tasks ("meta-tasks") instead of acting directly via `op=complete`. 7 such meta-tasks clogged the board with noise they were meant to cure.

## Root cause
`buildOutputInstruction` only showed the Observer how to create tasks (`op=intend`). It had no instruction for `op=complete` or `op=edit`, no distinction between administrative corrections vs. code-requiring findings, and no rule against creating closure tasks.

## Changes

### `pkg/runner/observer.go`

**`buildOutputInstruction`** — rewritten to teach the Observer the two-category heuristic:

- **Category A (administrative):** If the action requires no code change (closing false-positives, completing stale tasks, removing board noise), execute inline with `op=complete` or `op=edit`. Curl examples provided for both.
- **Category B (code change needed):** Only create a task if a Builder needs to write code. Max 2 tasks.
- **Explicit prohibition:** "Do NOT create a task to close another task. That is the defect you must avoid." and "Creating a task to close a task is always wrong. Close it yourself."

**`buildPart2Instruction`** — added item 7 to the audit checklist:

- Flags meta-tasks (tasks whose sole purpose is to close another task) as board noise
- Instructs the Observer to close both the meta-task AND the target task inline using `op=complete`
- Board hygiene rule: if a task title says "close task X" or "complete task Y", that's a meta-task — close it, don't create another task about it

## Verification
- `go.exe build -buildvcs=false ./...` — clean
- `go.exe test ./...` — all pass (13 packages)

## Effect
On the next Observer run, it will:
1. Detect the 7 existing meta-tasks during the board audit
2. Close them directly via `op=complete` rather than creating more tasks
3. Going forward, close any false-positive tasks inline instead of deferring

The fix is prompt-level — no API changes, no new types, no schema changes.
