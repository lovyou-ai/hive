# Test Report: Fix: Knowledge API omits causes field on claim nodes

## What Was Tested

The Builder added `causeIDs` propagation to `post()`, `assertScoutGap()`, `assertCritique()`, and `assertLatestReflection()`. The prior tests covered causes for scout gap and critique, but three gaps remained.

## Gaps Found and Filled

### 1. `assertLatestReflection` had no SendsCauses test
All three assert functions were updated to carry `causes` in this iteration, but only `assertScoutGap` and `assertCritique` had `SendsCauses` tests. `assertLatestReflection` was missing one.

**Added:** `TestAssertLatestReflectionSendsCauses` — verifies `causes` field is present when `causeIDs` is non-empty.

### 2. Multiple causeIDs never tested joined
The implementation uses `strings.Join(causeIDs, ",")` to build the CSV value. No test covered the multi-ID case — a regression could break the join separator without any test failing.

**Added:** `TestAssertCauseIDsMultipleJoined` — passes `["id-aaa", "id-bbb", "id-ccc"]` and asserts the payload receives `"id-aaa,id-bbb,id-ccc"`.

### 3. `post()` empty response returns empty ID untested
`post()` silently ignores JSON decode failures and returns `("", nil)`. `main()` guards on this (`if buildDocID != ""`). The path was untested — a regression here would silently drop all causal links.

**Added:** `TestPostEmptyResponseReturnsEmptyID` — server returns `{}` (no `node.id`), asserts `("", nil)` is returned without error.

## Results

```
ok  github.com/lovyou-ai/hive/cmd/post  0.606s
```

33 tests total (30 pre-existing + 3 new). All pass.

## Coverage Notes

- All new functions from this iteration have tests.
- `syncClaims` returning `causes` in the markdown is not tested here — that's `site/graph/knowledge_test.go` territory (covered by `TestAssertOpReturnsCauses` per the build report).
- No database tests needed — `cmd/post` is pure HTTP client code with no DB calls.
