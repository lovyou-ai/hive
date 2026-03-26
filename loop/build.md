# Build: Enrich `writeBuildArtifact` with commit subject, diff stat, and task body

- **Commit:** (pending — not yet committed)
- **Subject:** (pending)
- **Cost:** $0.0000
- **Timestamp:** 2026-03-27T00:00:00Z

## Task

In `pkg/runner/runner.go`, add two helper methods `gitSubject()` (git log -1 --format=%s) and `gitDiffStat()` (git show --stat HEAD, truncated to 1000 chars) alongside the existing `gitHash()`. Rewrite `writeBuildArtifact` to include commit message, diff stat, and truncated task body (300 chars) in build.md so the Reflector has real substance to reflect on.

## Changes

- `pkg/runner/runner.go` — Added `gitSubject()` and `gitDiffStat()` helpers; rewrote `writeBuildArtifact` to emit subject, task body (truncated to 300 chars), and diff stat (truncated to 1000 chars) sections.
