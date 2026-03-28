# Critique: [hive:builder] Close orphaned subtasks when parent completes

**Verdict:** PASS

**Summary:** Fix task created: `b39219e9` — two items to address:

1. **`store_test.go:429-430`** — update the `TestUpdateNodeStateChildGateMultipleChildren` doc comment from "still blocks parent completion" to describe the new cascade-close behavior.
2. **`store.go:211`** — remove `ErrChildrenIncomplete` (no code path returns it; it's a dead sentinel that misleads callers).
