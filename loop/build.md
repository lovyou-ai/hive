# Build: Add Preview field to PhaseEvent (already implemented)

Both changes were already present in the codebase before this iteration — `Preview string \`json:"preview,omitempty"\`` in `PhaseEvent` (diagnostic.go:17) and the 1000-char capture in the architect parse-failure branch (architect.go:69-85).

- `go.exe build -buildvcs=false ./...` — clean
- `go.exe test ./...` — all pass

---

# Previous Build: Fix: [hive:builder] Add JSON output format support to `parseArchitectSubtasks`

- **Commit:** cf989d0777f336b9f5eaa1eb9069b910c7aaa887
- **Subject:** [hive:builder] Fix: [hive:builder] Add JSON output format support to `parseArchitectSubtasks`
- **Cost:** $0.4268
- **Timestamp:** 2026-03-26T22:39:46Z

## Task

Critic review of commit a85ad1cf85f5 found issues:

---

## Analysis

### What was built

Tasks 2 and 3 are solid:
- `parseSubtasksJSON` is well-structured: tries bare array, falls back to `{"tasks":[]}` wrapper, normalizes unknown priorities, skips empty titles.
- Early-exit guard `content == "" ||...

## Diff Stat

```
commit cf989d0777f336b9f5eaa1eb9069b910c7aaa887
Author: hive <hive@lovyou.ai>
Date:   Fri Mar 27 09:39:46 2026 +1100

    [hive:builder] Fix: [hive:builder] Add JSON output format support to `parseArchitectSubtasks`

 loop/budget-20260327.txt     |  3 +++
 loop/build.md                | 35 +++++++++++++++++---------------
 loop/critique.md             | 47 +++++++++++++++++++++----------------------
 loop/mcp-knowledge.json      | 17 ++++++++++++++++
 loop/reflections.md          | 10 +++++++++
 loop/state.md                |  2 +-
 pkg/runner/architect.go      | 12 ++++++-----
 pkg/runner/architect_test.go | 48 ++++++++++++++++++++++++++++++++++++++++++++
 pkg/runner/diagnostic.go     |  1 +
 9 files changed, 129 insertions(+), 46 deletions(-)
```
