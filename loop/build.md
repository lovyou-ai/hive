# Build: Fix: syncClaims uses knowledge endpoint instead of board search

## What Was Built

Replaced `syncClaims`'s board-query loop with a single fetch from the knowledge endpoint (`/app/hive/knowledge?tab=claims`), which returns all 187 claims without the ~68-node server-side cap that was truncating board search results.

## Files Changed

### `cmd/post/main.go`
- **`syncClaims`**: Replaced the two-pass board query loop (one per prefix: "Lesson ", "Critique:") with a single call to `fetchKnowledgeClaims`. Title-prefix filter (`hasClaimPrefix`) kept to exclude non-Lesson/Critique claims.
- **`fetchKnowledgeClaims`** (new): Fetches all claims from `/app/hive/knowledge?tab=claims`, returns `[]boardNode`. Parses `{"claims": [...]}` JSON response shape.
- Updated stale comments on `boardNode`, `claimTitlePrefixes`, and `syncClaims` that said "the knowledge endpoint returns 0".

### `cmd/post/main_test.go`
- Updated 6 failing tests to serve `/app/hive/knowledge` returning `{"claims": [...]}`:
  - `TestSyncClaimsWritesFile`, `TestSyncClaimsEmptyDoesNotWrite`, `TestSyncClaimsFiltersNonClaimNodes`, `TestSyncClaimsClaimWithNoMetadata`, `TestSyncClaimsMultipleCauses`, `TestSyncClaimsWritesCauses`
- Renamed `TestSyncClaimsDeduplicatesAcrossQueries` → `TestSyncClaimsDeduplicatesNodes`.
- Replaced `TestSyncClaimsSecondQueryFails` with `TestSyncClaimsKnowledgeEndpointFails`.
- Added 3 new tests for `fetchKnowledgeClaims`: `TestFetchKnowledgeClaimsReturnsNodes`, `TestFetchKnowledgeClaimsSendsAuthHeader`, `TestFetchKnowledgeClaimsHTTPError`.

## Verification

```
go.exe build -buildvcs=false ./...  → OK
go.exe test ./...                   → all pass
```

## Root Cause

Board search (`/app/hive/board?q=Lesson&limit=500`) is server-side capped at ~68 rows regardless of the `limit` param. The knowledge endpoint (`/app/hive/knowledge?tab=claims`) calls `ListNodes(KindClaim)` with no row cap — returns all 187 claims. Result: claims.md now captures 102 lessons instead of 4.

ACTION: DONE
