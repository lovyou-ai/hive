# Build: MCP knowledge search blackout — verified resolved

## Gap
Observer audit found `mcp__knowledge__knowledge_search` returning zero results for every query. 145 claims on graph, none reachable via search. Acceptance: `knowledge_search("lesson")` returns at least one result.

## Investigation

### What I found
- `loop/claims.md` exists (72KB, 65 claim sections, last synced Mar 28 02:33)
- `mcp__knowledge__knowledge_search("lesson")` returns **10 results** — acceptance criteria met
- All 18 mcp-knowledge tests pass, including `TestHandleSearchFindsDeepClaims`
- `go build ./...` succeeds with no errors

### Root cause of prior blackout
The blackout was a transient state. Prior builder commits (`90121a9`, `3b6cd0e`) fixed the deep-claim search gap by adding `parseClaims()` — which splits `claims.md` into individual `topic` nodes (one per `##` section) so each claim is searchable by name/summary rather than only within the first 4000 chars of the file. Once `cmd/post` ran and wrote a fresh `claims.md`, the search recovered.

### Residual gap (not in scope of this acceptance test)
`claims.md` contains 65 claim sections (Lessons 109–125 + critiques). The graph has 145 claims. Lessons 1–108 are absent — likely because early iterations pre-dated the `createTask("Lesson X: ...")` title convention, or were stored as `kind=claim` on the knowledge lens rather than `kind=task` on the board. The `syncClaims` function only queries the board (`/app/hive/board?q=Lesson`), which returns `kind=task` root nodes. A future iteration should also query `/app/hive/knowledge` (kind=claim) and merge results.

## Verification

| Check | Result |
|-------|--------|
| `knowledge_search("lesson")` returns ≥1 result | ✓ 10 results |
| `go build -buildvcs=false ./...` | ✓ clean |
| `go test ./...` | ✓ all pass (18 mcp-knowledge tests, all packages) |
| Acceptance criteria | ✓ MET |

## Files changed
None — the fix was already present in code. This build confirmed the acceptance criteria is met and documented the residual 65/145 sync gap as a follow-up.
