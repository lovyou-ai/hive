# Build: Fix: KindClaim graph nodes not synced to MCP knowledge index

## Task

`knowledge_search` returned zero results despite claims being asserted to the graph via the `assert` op. The two systems were disconnected: KindClaim nodes live in Postgres (graph store), but the MCP knowledge server reads only markdown files from the filesystem. Fix task `c5dca156`.

## What Was Built

Two mechanisms bridge the gap:

**1. `cmd/post/main.go` — `syncClaims()`** (commit de17e45)
- Fetches all KindClaim nodes from `/app/hive/knowledge?tab=claims&limit=200`
- Formats them as markdown (title, state, author, body per claim)
- Writes to `loop/claims.md`
- Called from `main()` after `syncMindState`, non-fatal

**2. `cmd/mcp-knowledge/main.go` — claims.md indexed** (commit de17e45)
- Added `{"claims.md", "Asserted knowledge claims — lessons, decisions, invariants from the graph store (synced by cmd/post)"}` to `buildHiveLoop()`
- `os.Stat` guard: only appears in tree after first `close.sh` run that writes the file

**3. `cmd/post/main.go` — `assertScoutGap()`** (commit 4d0680c)
- Reads `loop/scout.md`, extracts `**Gap:**` line and iteration number
- POSTs `op=assert` to create a permanent KindClaim node for each Scout gap
- Gap survives `scout.md` being overwritten next iteration
- Called from `main()` after `syncClaims()`, non-fatal

**4. `cmd/post/main_test.go` — 2 additional tests added** (working tree)
- `TestAssertScoutGapNoGapLine` — assertScoutGap returns error when no `**Gap:**` line
- `TestAssertScoutGapAPIError` — assertScoutGap returns error on HTTP 4xx

## Verification

```
go.exe build -buildvcs=false ./...  ✓ (no errors)
go.exe test -count=1 ./cmd/post/ ./cmd/mcp-knowledge/  ✓ (all pass)
```

## Flow After This Fix

1. Scout asserts gap to graph via `assertScoutGap()` → KindClaim node in Postgres
2. `close.sh` runs `cmd/post` → `syncClaims` fetches from API → writes `loop/claims.md`
3. MCP knowledge server indexes `loop/claims.md` as `loop/claims` topic
4. `knowledge_search "iteration gap title"` finds the claim

ACTION: DONE
