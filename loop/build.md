# Build Report — Iteration 239 (Fix 3)

## Gap
Critic flagged invariant 13 (BOUNDED) violation: `computeHiveStats(posts []Node)` accepted an unbounded slice. The handler passed a magic literal `20` to `ListHiveActivity`, but the bound was implicit — not named, not traceable to the consumer.

## What Changed

### `site/graph/handlers.go`
- Added `const maxHivePosts = 20` (above `HiveStats`) citing invariant 13 (BOUNDED).
- Updated `computeHiveStats` doc comment: "posts must be pre-bounded (callers must pass at most maxHivePosts entries)."
- Changed `h.store.ListHiveActivity(r.Context(), "", 20)` → `h.store.ListHiveActivity(r.Context(), "", maxHivePosts)`.

The bound was always enforced by the store LIMIT clause — this fix makes it **explicit, named, and traceable** between fetch and process.

## Verification

```
go.exe build -buildvcs=false ./...   → clean (no output)
go.exe test ./...                    → all pass (graph: 0.561s)
```

Hive tests passing:
- `TestParseCostDollars`
- `TestParseDurationStr`
- `TestComputeHiveStats`
- `TestGetHive_PublicNoAuth`
- `TestGetHive_RendersMetrics`
