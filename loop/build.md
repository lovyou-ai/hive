# Build: Switch Reflector to JSON output with text fallback

## Task
Eliminate marker-format ambiguity causing `empty_sections` failures by switching the Reflector to JSON output with a text-marker fallback.

## Changes

### `pkg/runner/reflector.go`

Refactored JSON parsing to mirror the Architect pattern:

1. **Added `normalizeReflectorResponse`** — standalone fence-stripping function (mirrors `normalizeArchitectResponse` in `architect.go`). Strips opening fence lines (` ```json `, ` ```text `, ` ``` `) and closing ` ``` ` before any parsing.

2. **Renamed `parseReflectorOutputJSON` → `parseReflectorJSON`** — mirrors naming of `parseSubtasksJSON`. Removed internal fence-stripping (now done upstream by `normalizeReflectorResponse`). Still handles flat `{"cover":...}`, wrapped `{"reflection":{...}}`, and prose preamble before JSON.

3. **Updated `parseReflectorOutput`** — now calls `normalizeReflectorResponse(content)` first, then `parseReflectorJSON`, then falls back to text-marker parsing. Key fix: normalization now applies to BOTH parse paths. Previously, a fence-wrapped text-marker response was invisible to the text-marker parser (fences don't match any marker candidate).

`buildReflectorPrompt` already requested JSON output (from prior commits). No change needed.

## Verification

```
go.exe build -buildvcs=false ./...   → success (no output)
go.exe test -count=1 ./...           → all 12 packages pass
```

All reflector tests pass: JSON flat object, `{"reflection":{...}}` wrapper, prose preamble before JSON, all text-marker format variants (bold, h2, h3, lowercase, mixed).
