# Build: Fix: claims.md sync broken — Lessons 126-148 missing from MCP index

## Gap

`syncClaims()` in `cmd/post/main.go` queried `/app/hive/knowledge?tab=claims` which
filters for `kind=claim`. All 1500+ graph nodes are `kind=task`, so the endpoint
returned 0 nodes and `claims.md` was never updated past Lesson 125.

## Fix

Replaced the knowledge endpoint call with board endpoint queries. The board returns
`kind=task` nodes — where lessons and critiques are actually stored.

**New approach in `syncClaims()`:**
1. Query `/app/hive/board?q=Lesson ` → board nodes with "Lesson" in title/body
2. Query `/app/hive/board?q=Critique:` → board nodes with "Critique:" in title/body
3. Filter client-side: only keep nodes whose title starts with `"Lesson "` or `"Critique:"`
4. Dedup by ID, sort oldest-first, write to `loop/claims.md`

## Files Changed

- `cmd/post/main.go`
  - Added `net/url`, `sort`, `time` imports
  - Added `boardNode` struct (subset of board Node needed for claims.md)
  - Added `claimTitlePrefixes` var (`["Lesson ", "Critique:"]`)
  - Replaced `syncClaims()` implementation: board queries instead of knowledge endpoint
  - Added `fetchBoardByQuery()`: fetches `/app/hive/board?q=<prefix>` as JSON
  - Added `hasClaimPrefix()`: title prefix guard (prevents non-claim tasks from leaking in)

- `cmd/post/main_test.go`
  - Updated `TestSyncClaimsWritesFile`: mock `/app/hive/board`, verify lesson+critique appear in order
  - Updated `TestSyncClaimsEmptyDoesNotWrite`: mock board returning `{"nodes": []}`, file not written
  - Added `TestSyncClaimsFiltersNonClaimNodes`: board returns mix; non-prefix node excluded
  - Updated `TestSyncClaimsClaimWithNoMetadata`: lesson node with empty state/author
  - Updated `TestSyncClaimsMultipleCauses`: lesson node with two causes, both joined
  - Updated `TestSyncClaimsWritesCauses`: lesson node with causes, **Causes:** line present

## Verification

```
go.exe build -buildvcs=false ./...   OK
go.exe test -buildvcs=false ./...    OK  (all 13 packages, cmd/post: 1.109s)
```

## Acceptance

After next `close.sh` run, `mcp__knowledge__knowledge_search` for "Lesson 148" will
return the relevant claim because `claims.md` will be re-synced from board nodes
with "Lesson " title prefix — capturing all 148+ lessons previously invisible.
