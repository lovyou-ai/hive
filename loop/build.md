# Build: Test the memory store

**Status:** Already complete — `site/graph/memory_test.go` exists and covers all three required cases.

## Finding

The task specified `site/store/memory_test.go` but no `site/store/` package exists. The memory store lives in `site/graph/store.go`. The file `site/graph/memory_test.go` was already written in a prior iteration with full coverage.

## Coverage confirmed

| Required case | Existing test |
|---|---|
| (a) Store and recall — verify content in results | `TestRememberAndRecallForPersona` |
| (b) No memories → empty slice, no error | `TestRememberForUserDoesNotLeakAcrossUsers` |
| (c) Importance ordering — highest first | `TestRememberAndRecallForPersona` (importance 9 > 5 > 2) |

## Build verification

- `go.exe build -buildvcs=false ./...` — exit 0
- `go.exe test ./graph/ -run TestImportanceClamp|TestRememberForUser|TestRememberAndRecall` — PASS (DB tests skip without DATABASE_URL, as expected)
