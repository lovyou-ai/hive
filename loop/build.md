# Build Report — Iteration 226: Critic Role

## What This Iteration Does

Implements the Critic role for the hive runtime (Phase 2, item 9 from hive-runtime-spec.md). The Critic reviews builder commits automatically, creating fix tasks when issues are found.

## Files Changed

| File | Lines | What |
|------|-------|------|
| `pkg/runner/critic.go` | 170 | New. Critic role: find unreviewed builder commits, review diffs via Reason(), create fix tasks on REVISE. |
| `pkg/runner/critic_test.go` | 65 | New. 9 tests: parseVerdict, extractIssues, buildReviewPrompt. |
| `pkg/runner/runner.go` | -5 | Removed critic stub (now in critic.go). |

## How It Works

1. **Every 4th tick** (~60s), Critic runs `git log --grep=\[hive:builder\]` for commits in the last 24h
2. For each unreviewed commit, gets the diff via `git diff hash~1..hash`
3. Calls `Reason()` (no tools, haiku model) with the diff + review checklist
4. Parses `VERDICT: PASS` or `VERDICT: REVISE` from the response
5. On REVISE: creates a fix task on the lovyou.ai board
6. On PASS: logs and moves on

## Review Checklist (sent to LLM)

1. **Completeness** — new constant/kind present in ALL guards and allowlists?
2. **Identity (invariant 11)** — IDs for matching, not names?
3. **Bounded (invariant 13)** — queries have LIMIT?
4. **Correctness** — injection, races, nil handling?
5. **Tests** — flagged but not REVISE-blocking (known systemic issue)

## E2E Test Results

**Run 1 (bug):** Unescaped regex brackets in `git --grep` matched 54 commits instead of 1. Fixed.

**Run 2 (correct):**
```
[critic] tick 4: found 1 unreviewed commits
[critic] reviewing 31f3349ca8b6: [hive:builder] Add Policy entity kind to the site
  ⏳ thinking done (1m16s)
[critic] review done (cost=$0.1631)
[critic] verdict: PASS
```

1 commit found, 1m16s review time, $0.16 cost. Correct commit identified and reviewed.

## Build

- `go build ./...` ✓
- `go test ./...` ✓ (23 tests: 9 critic + 14 existing)
