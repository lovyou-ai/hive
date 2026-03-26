# Build: Add JSON output format support to parseArchitectSubtasks

## Gap
`parseArchitectSubtasks` had no JSON parser. A 1,282-token LLM response returned a JSON array but produced zero parsed tasks because only the strict (SUBTASK_TITLE:) and markdown parsers were tried.

## Changes

### `pkg/runner/architect.go`
- Added `encoding/json` import
- Added `jsonSubtask` struct (`title`, `description`, `priority` fields)
- Added `parseSubtasksJSON(content string) []architectSubtask`:
  - Tries bare JSON array `[{...}]` first
  - Falls back to `{"tasks":[...]}` wrapper object
  - Returns nil on empty/invalid JSON
  - Normalizes unknown priorities to "high"
- Updated `parseArchitectSubtasks` to call `parseSubtasksJSON` first, before the strict parser

### `pkg/runner/architect_test.go`
- Added `TestParseSubtasksJSON`: 6 cases covering bare array, wrapper object, invalid JSON, empty array, empty string, unknown priority
- Added `TestParseArchitectSubtasksJSON`: integration test verifying JSON is tried before strict/markdown

## Verification
- `go.exe build -buildvcs=false ./...` — clean
- `go.exe test ./...` — all pass (12 packages)
