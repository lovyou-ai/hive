# Build: Fix Critic REVISE items (observer + reflector diagnostics)

**Status:** DONE

## Changes

### `pkg/runner/observer.go`
- Added `"mcp__knowledge__knowledge_search"` to `AllowedTools` in `runObserver` so the instruction's `knowledge.search` reference is valid (was a lie — tool not in allowed list).
- Added `if apiKey == ""` guard in `runObserver` — logs a warning when `LOVYOU_API_KEY` is not set.
- Refactored `buildObserverInstruction` into three functions: `buildObserverInstruction`, `buildPart2Instruction`, `buildOutputInstruction`. When `apiKey == ""`, Part 2 and the output section skip authenticated curl commands and substitute a skip notice or text-format fallback.

### `pkg/runner/reflector.go`
- Added `r.appendDiagnostic(PhaseEvent{Phase: "reflector", Outcome: "empty_sections", Preview: ...})` in the `emptySections` early-return path in `runReflectorReason`. This was missing — `TestRunReflectorEmptySectionsDiagnostic` expected `diagnostics.jsonl` to be written on partial LLM output but the early return was silent.

### `loop/critique.md`
- Removed duplicate bare `VERDICT: REVISE` line at bottom. The `**Verdict:** REVISE` header at the top is the canonical location.

## Verification
- `go.exe build -buildvcs=false ./...` — clean
- `go.exe test -buildvcs=false ./...` — all pass (7 packages with tests)
