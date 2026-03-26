# Build: parseSubtasksStrict multi-line description accumulation

## Gap
`parseSubtasksStrict` was reported to only capture text on the same line as `SUBTASK_DESCRIPTION:`, dropping continuation lines.

## Finding
The fix was already applied in a prior iteration. `pkg/runner/architect.go:211` uses an `inDesc bool` state variable that tracks whether the parser is inside a description block. Lines not starting with a `SUBTASK_` prefix are appended to `current.desc` while `inDesc == true` (lines 227–232). `inDesc` resets to `false` on blank lines, code fences, or new `SUBTASK_TITLE:` / `SUBTASK_PRIORITY:` lines.

The test `TestParseArchitectSubtasks/multi-line_description_spanning_two_lines` in `pkg/runner/architect_test.go:94` covers exactly this case and passes.

## Changes
No code changes required — fix already present and tested.

## Verification
- `go.exe build -buildvcs=false ./...` — clean
- `go.exe test -count=1 ./pkg/runner/...` — ok (1.049s)
