# Build: Fix: add tests for buildPart2Instruction and buildOutputInstruction

## Task

Invariant 12 (VERIFIED): `observer.go` was refactored with `buildPart2Instruction` + `buildOutputInstruction` + `buildObserverInstruction` with apiKey=="" skip paths but tests were incomplete — not table-driven and missing `buildObserverInstruction` top-level coverage.

## What Was Built

Rewrote `pkg/runner/observer_test.go` with table-driven tests:

- `TestBuildPart2Instruction` — 2 cases: (a) empty apiKey → skip text, no curl; (b) set apiKey → curl with key+slug embedded
- `TestBuildOutputInstruction` — 2 cases: (a) empty apiKey → TASK_TITLE text format, no curl; (b) set apiKey → curl with key+slug, no text format
- `TestBuildObserverInstruction` — 2 new cases covering top-level format: (a) empty apiKey → Part 2 skip + text output section; (b) set apiKey → 2× curl auth headers (part2 + output), key+slug present throughout

## Verification

```
go.exe build -buildvcs=false ./...   ✓
go.exe test ./...                     ✓ all pass
```

`pkg/runner`: `ok github.com/lovyou-ai/hive/pkg/runner 4.009s`

ACTION: DONE
