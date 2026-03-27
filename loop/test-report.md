# Test Report: Observer meta-task anti-pattern (iter 370)

- **Build:** Observer process defect: creating cleanup meta-tasks instead of acting
- **Timestamp:** 2026-03-28

## What Was Tested

Two prompt changes in `pkg/runner/observer.go`:
1. `buildOutputInstruction` — two-category model (Category A: act inline, Category B: create task only if code needed)
2. `buildPart2Instruction` — item 7 added: meta-tasks are board noise, close inline with op=complete

## Test Functions

All in `pkg/runner/observer_test.go`:

| Test | What it verifies | Result |
|------|-----------------|--------|
| `TestBuildOutputInstructionCategoryModel` | Category A present with op=complete and op=edit curl examples; Category B present; hard rule "Creating a task to close a task is always wrong" present | PASS |
| `TestBuildOutputInstructionNoAntiPatternWhenNoKey` | No-key fallback path does NOT contain Category A/B model (irrelevant without direct API access) | PASS |
| `TestBuildPart2InstructionMetaTaskItem` | Item 7 (Meta-tasks) in Part 2 checklist; instructs op=complete inline closure; "Do not create a new task for this" present; Board hygiene rule section present | PASS |
| `TestBuildPart2InstructionMetaTaskItemSkippedWhenNoKey` | Meta-task instructions absent when Part 2 is skipped (no API key) | PASS |

## Full Package Run

```
ok  github.com/lovyou-ai/hive/pkg/runner  3.599s
```

All 13 packages in the build compile and pass (Builder confirmed).

## Coverage Notes

- Pure function tests — no DB required, no network calls
- All four branches of the new logic are covered: key/no-key × category-model/meta-task-item
- Edge cases covered by prior tests: empty input, invalid priority defaulting, whitespace trimming, claims summary truncation at 5

## Verdict

PASS. The anti-meta-task invariant is enforced in the prompt text and verified by tests.
