# Build: All 103 claims have causes=[] � close.sh assertion pipeline never sets causes

- **Commit:** db6dc49ac1cae71465b5761aed459c6409a7ae87
- **Subject:** [hive:builder] All 103 claims have causes=[] � close.sh assertion pipeline never sets causes
- **Cost:** $2.4943
- **Timestamp:** 2026-03-28T08:16:16Z

## Task

Every claim on the knowledge board (103 lessons + 37 critiques) has causes=[]. This is a systemic Invariant 2 (CAUSALITY) violation. Root cause: close.sh does not set causes when calling op=assert for either lessons or critiques.

Fix required:
1. Lessons: before asserting a lesson claim, close.sh m...

## What Was Built

All done. Here's a summary of what was built:

## Changes

**`site/graph/store.go`** — Added `UpdateNodeCauses(ctx, nodeID, causes)` that issues `UPDATE nodes SET causes = $1 WHERE id = $2`. Required by the backfill path.

**`site/graph/handlers.go`** — Extended `op=edit` to accept an optional `causes` field. Validation changed from `body required` to `body OR causes required`. Calls `UpdateNodeCauses` when causes is provided. Records op for observability.

**`cmd/post/main.go`** — Three fixes:
1. Moved `taskCauseIDs` computation before `assertScoutGap` — both it and `assertCritique` now use the task node ID (not build doc ID) as cause
2. Added `backfillClaimCauses` — fetches all claims, POSTs `op=edit causes=<taskNodeID>` for any with empty causes, skips ones already satisfied
3. Added `ID` field to the `syncClaims` struct (was silently dropped from JSON)

**`cmd/post/main_test.go`** — 4 new tests covering the backfill function.

The backfill runs automatically on every `close.sh` run. The 103 existing empty-cause claims will be linked to the current iteration's task node on the next run. Future claims (lessons and critiques) will have correct causes from creation.

ACTION: DONE

## Diff Stat

```
commit db6dc49ac1cae71465b5761aed459c6409a7ae87
Author: hive <hive@lovyou.ai>
Date:   Sat Mar 28 19:16:15 2026 +1100

    [hive:builder] All 103 claims have causes=[] � close.sh assertion pipeline never sets causes

 cmd/post/main.go         | 107 +++++++++++++++++++++++++++++++++++++++++---
 cmd/post/main_test.go    | 113 +++++++++++++++++++++++++++++++++++++++++++++++
 loop/budget-20260328.txt |   1 +
 loop/build.md            |  68 +++++++++++++---------------
 4 files changed, 245 insertions(+), 44 deletions(-)
```
