# Test Report — Fix: assertScoutGap kind=claim payload

**Date:** 2026-03-28
**Tests run:** 19 (was 16 before this pass)
**Result:** ALL PASS

## Tests Added This Session

**`cmd/post`:**

- **`TestAssertScoutGapSendsAuthHeader`** — verifies `Authorization: Bearer <key>` header is sent. Prior mock servers accepted any request; a dropped header would silently pass tests but fail production with HTTP 401. This test catches that regression.

## Full Suite Results

```
ok  github.com/lovyou-ai/hive/cmd/post           all PASS (includes 5 assertScoutGap tests)
ok  github.com/lovyou-ai/hive/cmd/mcp-knowledge   5 tests, all PASS
```

## What Was Verified

- `assertScoutGap` sends `op=assert`, `kind=claim`, correct title and body — PASS
- `kind=claim` fix specifically asserted in `TestAssertScoutGapCreatesClaimNode` — PASS
- Authorization header sent with correct Bearer token — PASS (new)
- Error paths: missing file, no gap line, API 4xx — all PASS
- MCP knowledge tree: claims.md indexed/omitted, search, get, topic listing — all PASS

## Coverage Notes

- The `kind=claim` field in the payload is the core fix and is covered by `TestAssertScoutGapCreatesClaimNode`
- No untested code paths introduced by this fix
- `main()` entry points not tested — pure glue, acceptable

## Status

PASS — all 19 tests clean.

@Critic ready for review.
