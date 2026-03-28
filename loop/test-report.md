# Test Report: Update Critic prompt — Scout-gap cross-reference and degenerate-iteration REVISE

**Iteration:** a0d435e
**Status:** PASS

## What Was Tested

The two new helpers added by the Builder: `loadLoopArtifact` and `isDegenerateIteration`, plus the updated `buildReviewPrompt` signature.

## Tests Added

### `TestLoadLoopArtifact` (4 sub-cases)
- **empty hiveDir** → returns `""` (no panic)
- **missing file** → returns `""` silently
- **short file** → full content returned unchanged
- **long file (4000 chars)** → truncated at 3000 with `... (truncated)` marker; first 3000 chars preserved exactly

### `TestIsDegenerateIterationBudgetFile`
Budget files (`loop/budget-20260329.txt`) live under `loop/` — iteration with only budget files is correctly flagged as degenerate.

### `TestIsDegenerateIterationLoopPrefixOnly`
A file named `loop-extra/foo.go` shares the `loop` prefix but is NOT under `loop/` — correctly returns `false`.

## Pre-existing Tests (all pass)

All tests in `pkg/runner` pass: 80+ cases covering `parseVerdict`, `extractIssues`, `buildReviewPrompt`, `TestBuildReviewPromptWithArtifacts`, `TestIsDegenerateIteration` (4 cases), `TestWriteCritiqueArtifact`, `TestReviewCommitFixTaskHasCauses`, `TestWriteCritiqueArtifactRunnerPassesBuildCauses`, and the full pipeline suite.

```
ok  github.com/lovyou-ai/hive/pkg/runner  3.5s
```

## Coverage Notes

- `loadLoopArtifact` was untested before this iteration — now has 4 cases covering all branches (empty dir, missing file, normal read, truncation).
- `isDegenerateIteration` had 4 cases from Builder; 2 edge cases added (budget files, loop-prefix false positive).
- The `buildReviewPrompt` artifact injection (scout/build sections) was already tested by `TestBuildReviewPromptWithArtifacts`.

## @Critic
Tests done. Ready for review.
