# Build: Normalize LLM response before parsing — strip fences, guard zero-value

## Gap
LLM responses wrapped in markdown code fences were reaching `parseArchitectSubtasks` unparsed. Subtasks with empty titles could reach `CreateTask` violating the API contract. The markdown fallback silently consumed preamble text when strict markers were present but strict parsing failed.

## Changes

### `pkg/runner/architect.go`

**`normalizeArchitectResponse(content string) string`** (new function)
- Strips opening code fence line (` ```json `, ` ```text `, or plain ` ``` `)
- Strips closing ` ``` `
- Called at the top of `parseArchitectSubtasks` before any parsing

**`parseArchitectSubtasks`**
- Calls `normalizeArchitectResponse` first
- Adds debug log when `SUBTASK_TITLE:` markers are present but strict parsing returns 0 tasks — exposes format mismatches rather than silently falling through to markdown

**`runArchitect`**
- Guards against subtasks with empty titles before calling `CreateTask`
- Logs and skips empty-title subtasks

## Verification

- `go.exe build -buildvcs=false ./...` — clean
- `go.exe test ./...` — all pass including `TestParseArchitectSubtasks/fence-wrapped_response`
