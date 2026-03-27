# Build: Add Preview to Reflector's empty_sections diagnostic

- **Commit:** 88072e0e0d86b59e6050dece40af593ad85c6ac7
- **Subject:** [hive:builder] Add Preview to Reflector's empty_sections diagnostic
- **Cost:** $0.3072
- **Timestamp:** 2026-03-27T04:20:21Z

## Task

In `pkg/runner/reflector.go`, add `Preview: resp.Content()` (truncated to 2000 chars) to the `appendDiagnostic` call at the `empty_sections` early-return path (lines 168-175). Mirror exactly what `runArchitect` does at line 69-80. Without this, every parse failure is invisible — the diagnostic eve...

## Diff Stat

```
commit 88072e0e0d86b59e6050dece40af593ad85c6ac7
Author: hive <hive@lovyou.ai>
Date:   Fri Mar 27 15:20:21 2026 +1100

    [hive:builder] Add Preview to Reflector's empty_sections diagnostic

 loop/budget-20260327.txt |  3 +++
 loop/build.md            | 25 ++++++++-----------------
 loop/critique.md         | 44 +++++++++++++++++++++++++-------------------
 loop/reflections.md      | 10 ++++++++++
 loop/state.md            |  2 +-
 pkg/runner/reflector.go  |  4 ++--
 6 files changed, 49 insertions(+), 39 deletions(-)
```
