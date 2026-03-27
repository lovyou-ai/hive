# Build: Bridge KindClaim graph nodes to MCP knowledge index

## Gap

`knowledge_search` returned zero results despite claims being asserted to the graph via the `assert` op. The two systems were disconnected: KindClaim nodes live in Postgres (graph store), but the MCP knowledge server reads only markdown files from the filesystem.

## What Changed

### `cmd/post/main.go`
- Added `syncClaims(apiKey, baseURL, outPath string) error`
  - Fetches KindClaim nodes from `/app/hive/knowledge?tab=claims&limit=200`
  - Formats them as markdown (title, state, author, body per claim)
  - Writes to `loop/claims.md`
- Called from `main()` after `syncMindState`, non-fatal

### `cmd/mcp-knowledge/main.go`
- Added `{"claims.md", "..."}` entry to `buildHiveLoop()`
- The `os.Stat` guard means claims.md only appears in the tree when it exists (after the first `close.sh` run that syncs claims)

### `cmd/post/main_test.go`
- `TestSyncClaimsWritesFile` — verifies claims are fetched and written as markdown
- `TestSyncClaimsEmptyDoesNotWrite` — verifies no file written for empty claims response

### `cmd/mcp-knowledge/main_test.go` (new file)
- `TestBuildHiveLoopIncludesClaimsWhenPresent` — tree includes `loop/claims` when file exists
- `TestBuildHiveLoopOmitsClaimsWhenAbsent` — tree omits it when file absent
- `TestHandleSearchFindsClaims` — `knowledge_search` finds content from claims.md
- `TestHandleGetClaims` — `knowledge.get("loop/claims")` returns full claim content

## Verification

```
go.exe build -buildvcs=false ./...  ✓ (no errors)
go.exe test ./...                   ✓ (all pass, 9 new tests)
```

## Flow After This Fix

1. Scout/Reflector assert claims via `c.AssertClaim(slug, title, body)` → Postgres
2. `close.sh` runs `cmd/post` → `syncClaims` fetches from API → writes `loop/claims.md`
3. MCP knowledge server indexes `loop/claims.md` as `loop/claims` topic
4. `knowledge_search "iteration gap title"` finds the claim

Invariant VERIFIED: after `close.sh`, `knowledge_search` for any asserted claim title will return results.
