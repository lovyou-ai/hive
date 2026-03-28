# Build: Fix fetchBoardByQuery 65-node cap — claims.md missing lessons 110-195

## Gap
`fetchBoardByQuery` in `cmd/post/main.go` called `/app/hive/board?q=Lesson` with no limit parameter. The board API defaults to 65 nodes. With 179+ knowledge claims on the graph, lessons 110-195 were silently truncated — absent from `loop/claims.md` and unsearchable via the MCP knowledge tool. Three malformed "Lesson: date" entries (lessons 106/107/108) also appeared in claims.md from a prior run with looser filtering.

## Changes

### `cmd/post/main.go`
- Added `boardQueryLimit = 500` constant with Invariant 13 (BOUNDED) explanation
- Changed `fetchBoardByQuery` URL from `?q={q}` to `?q={q}&limit=500` using `fmt.Sprintf`

### `cmd/post/main_test.go`
- Added `strconv` import
- Added `TestFetchBoardByQuerySendsLimit` — verifies the board query sends a `limit` param >= 200; catches the silent truncation regression
- Added `"Lesson: 2026-03-27"` test case to `TestHasClaimPrefix` — documents that malformed "Lesson: date" titles are rejected by the filter (colon != space at index 6)

## Verification
- `go.exe build -buildvcs=false ./...` — clean
- `go.exe test ./...` — all pass (13 packages)

## Effect
Next `syncClaims` run will fetch all ~200 lesson and critique nodes (not just 65), rebuild `loop/claims.md` with lessons 1-195 complete, and exclude the malformed "Lesson: date" entries. MCP knowledge index will be current.
