# Scout Report — Iteration 226

## Gap Identified

**No automated code review.** The builder ships code but nobody checks it. Iter 225's allowlist miss was caught manually. The Critic role is Phase 2 item 9 from hive-runtime-spec.md.

The builder follows patterns (99.1% accuracy). The Critic enforces completeness — checking that every code path touched by a change is updated. This is the difference between "almost works" and "works."

## Plan

Implement `runCritic` in the runner:

1. On each tick, scan `git log` for recent `[hive:builder]` commits not yet reviewed
2. Get the diff for each unreviewed commit
3. Call `Reason()` (no tools, fast, cheap — haiku) with the diff + review checklist
4. If REVISE: create a fix task on the board, comment on the original task
5. If PASS: comment "reviewed" on the task
6. Track reviewed commits to avoid re-reviewing

The Critic uses `Reason()` (not `Operate()`) — it reads diffs and thinks, it doesn't write code. This keeps it fast and cheap (haiku model).

## Priority

**P0** — Critic is the quality gate that turns "almost correct" into "correct."
