# Build: cmd/post: dedup loop header tasks (Iteration N) on board

- **Commit:** 537beb41b7bf5d730eec1bf26f324f7b36424261
- **Subject:** [hive:builder] cmd/post: dedup loop header tasks (Iteration N) on board
- **Cost:** $0.8811
- **Timestamp:** 2026-03-29T09:33:55Z

## Task

In hive/cmd/post/main.go, extend createTask dedup guard to cover ALL titles, not just Fix: prefixed ones. Before calling intend op, call findExistingTask(apiKey, baseURL, title) for any non-empty title. If found, comment instead of creating. This stops 11+ duplicate Iteration N and Target repo tasks...

## What Was Built

Build passes, all tests pass. The dedup guard in `createTask` now fires unconditionally for all non-empty titles — "Iteration N", "Target repo", and any other repeated titles will comment on the existing task instead of creating a new one.

ACTION: DONE

## Diff Stat

```
commit 537beb41b7bf5d730eec1bf26f324f7b36424261
Author: hive <hive@lovyou.ai>
Date:   Sun Mar 29 20:33:54 2026 +1100

    [hive:builder] cmd/post: dedup loop header tasks (Iteration N) on board

 cmd/post/main.go         | 13 ++++++------
 cmd/post/main_test.go    | 17 +++++++++------
 loop/budget-20260329.txt |  1 +
 loop/build.md            | 55 +++++++++++++++++++++++++++---------------------
 4 files changed, 49 insertions(+), 37 deletions(-)
```
