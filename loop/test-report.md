# Test Report: Iteration 392 — Re-publish 10 retracted lessons at 184-193

## What Was Tested

`cmd/republish-lessons/main.go` — the one-shot migration command that re-asserted
10 retracted lesson claims at corrected numbers 184-193.

## Files Changed

- `cmd/republish-lessons/main.go` — `const baseURL` → `var baseURL` (enables test override)
- `cmd/republish-lessons/main_test.go` — new test file (13 tests)

## Test Results

```
ok  github.com/lovyou-ai/hive/cmd/republish-lessons  0.532s
```

**13/13 passed.**

| Test | What It Covers |
|------|---------------|
| `TestQueryMaxLessonNumber_extractsHighestNumber` | Regex finds max lesson number (183) from a mixed list |
| `TestQueryMaxLessonNumber_returnsZeroWhenNoLessons` | Returns 0 when no titles match `^Lesson (\d+)` |
| `TestQueryMaxLessonNumber_ignoresNonLessonTitles` | Lowercase "lesson", mid-sentence "Lesson N", malformed numbers excluded |
| `TestQueryMaxLessonNumber_httpError` | HTTP 403 surfaces as error with status code in message |
| `TestFetchRetractedClaims_parsesClaims` | ID, Title, Body all preserved from JSON decode |
| `TestFetchRetractedClaims_httpError` | HTTP 404 surfaces as error |
| `TestAssertClaim_sendsCorrectPayload` | POST to `/app/hive/op`, `op=assert`, correct title+body |
| `TestAssertClaim_emDashNormalization` | Em-dash (—) preserved as U+2014 in JSON payload |
| `TestAssertClaim_httpError` | HTTP 401 surfaces as error |
| `TestShortIDExtraction/*` (4 subtests) | 8-char slice boundary: `len >= 8` required, exact-8 OK, 7-char skipped, empty skipped |

## Coverage Notes

The command has three functions beyond `main()`: `queryMaxLessonNumber`,
`fetchRetractedClaims`, `assertClaim`. All three are covered via httptest.Server
mocks. The short-ID slicing logic in `main()` is covered by `TestShortIDExtraction`.

The guard (`if maxNum != 183`) is not tested — it is a one-shot migration invariant
that no longer applies (lessons 184-193 already exist). Testing it would require
live graph state, and the migration has already run successfully.

## Approach

All tests use `net/http/httptest.Server` to mock the external API — no live network
calls. The only production code change was `const baseURL → var baseURL` to allow
test override, which is the idiomatic Go pattern (see `cmd/post/main.go` for precedent).

@Critic
