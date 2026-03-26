# Build: Capture LLM Preview in Architect Diagnostic

## What Changed

### `pkg/runner/diagnostic.go`
- Added `Preview string \`json:"preview,omitempty"\`` field to `PhaseEvent`.

### `pkg/runner/architect.go`
- Increased preview capture from 500 → 1000 chars (matching the spec).
- Populated `Preview` field in the `PhaseEvent` emitted on parse failure, so the LLM response content is written to `diagnostics.jsonl` and survives the run.
- Fixed inaccurate comment on `jsonSubtask` that claimed camelCase field names were accepted (only lowercase JSON tags are declared).

### `pkg/runner/architect_test.go`
- Added `encoding/json` import.
- Extended `TestRunArchitectParseFailureWritesDiagnostic` to assert `"preview"` field is present and contains LLM response content.
- Added `TestRunArchitectParseFailurePreviewTruncatedAt1000` to verify the preview is capped at 1000 chars.

## Verification

```
go.exe build -buildvcs=false ./...  → clean
go.exe test ./...                   → all 12 packages pass
```

## Why

The original failure was a 1,282-token LLM response that was permanently lost because it wasn't captured in the diagnostic. The JSON parser fix (prior commit) prevents future failures of that specific format. This change closes the remaining gap: if any format causes a parse failure, the LLM content (first 1000 chars) is now written to `diagnostics.jsonl` so PM/Scout can diagnose it without re-running.
