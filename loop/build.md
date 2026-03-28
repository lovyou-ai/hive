# Build: Replace GetClaims(200) with server-side MAX aggregate in NextLessonNumber

## Task
Replace `GetClaims(200)` client-side MAX scan with server-side aggregate MAX query in `NextLessonNumber`.

## Gap
`NextLessonNumber` fetched up to 200 claims and scanned them locally. At lesson 201 the scan misses claims and produces duplicate lesson numbers.

## Changes

### `site/graph/store.go`
Added `MaxLessonNumber(ctx, spaceID)` — SQL aggregate:
```sql
SELECT COALESCE(MAX(CASE WHEN title ~ '^Lesson [0-9]+'
    THEN CAST(REGEXP_REPLACE(title, '^Lesson ([0-9]+).*', '\1') AS INTEGER)
    ELSE 0 END), 0)
FROM nodes WHERE space_id=$1 AND kind='claim' AND parent_id IS NULL
```
O(1) correct at any lesson count.

### `site/graph/handlers.go`
Added `?op=max_lesson` fast path in `handleKnowledge`. Returns `{"max_lesson": N}` early — no full claims fetch.

### `site/graph/knowledge_test.go`
Added `TestMaxLessonNumberEndpoint` and `TestMaxLessonNumberEndpointEmpty`.

### `hive/pkg/api/client.go`
Removed `parseLessonNumber` (dead after refactor). Rewrote `NextLessonNumber` to call `?op=max_lesson`, parse `{"max_lesson": N}`, return `N+1`. Falls back to 1 on error.

### `hive/pkg/api/client_test.go`
Replaced old claims-list tests with `TestNextLessonNumberFromServer` and `TestNextLessonNumberNoLessons`.

### `hive/pkg/runner/reflector_test.go`
Updated `TestRunReflectorReasonLessonNumberFromGraph` mock: `op=max_lesson` route returning `{"max_lesson":109}`.

## Verification
- `go.exe build -buildvcs=false ./...` clean (hive + site)
- `go.exe test ./...` all 11 packages pass (hive)
- Site DB tests skip locally; run in CI with Postgres
