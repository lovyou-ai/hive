# Build Report — Fix: Architect Operate path missing causes

## Gap
`buildArchitectOperateInstruction` never received the milestone ID, so the curl template the LLM executes omitted `causes`. Since claude-cli implements `IOperator`, this is the production path. The test covered only the fallback (Reason) path.

## Changes

### `pkg/runner/architect.go`

**Call site**: Extract `milestoneID` before calling `buildArchitectOperateInstruction`, pass it as a third argument.

**`buildArchitectOperateInstruction`** signature changed from `(context, spaceSlug string)` to `(context, spaceSlug, milestoneID string)`.

When `milestoneID != ""`, a `causesSuffix` of `,"causes":["<ID>"]` is injected into the curl payload template, so the LLM's generated commands produce:

```json
{"op":"intend","kind":"task","title":"...","description":"...","priority":"high","causes":["milestone-42"]}
```

When no milestone is present (Scout fallback path), `causesSuffix` is empty and the payload is unchanged.

### `pkg/runner/architect_test.go`

Added `mockCaptureOperator` — implements `intelligence.Provider` + `decision.IOperator`. Captures `OperateTask.Instruction` for assertion.

Added `TestRunArchitectOperateInstructionIncludesCauses` — creates a milestone, runs `runArchitect` with the capture operator, asserts the instruction contains `"causes":["milestone-42"]`.

## Verification

```
go.exe build -buildvcs=false ./...   ✓ no errors
go.exe test -buildvcs=false ./...    ✓ all pass
```

New test: `TestRunArchitectOperateInstructionIncludesCauses` — PASS
Existing test: `TestRunArchitectSubtasksHaveCauses` — PASS (fallback path unaffected)
