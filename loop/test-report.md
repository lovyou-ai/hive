# Test Report: Fix causes field omitempty (Iteration 368)

**Tester:** Tester agent
**Timestamp:** 2026-03-28
**Build commit:** e98adfc (hive), site HEAD: 9ed933a

## What Was Tested

The fix: removed `omitempty` from `Node.Causes []string json:"causes"` in `site/graph/store.go`.
This ensures the `causes` key is always present in JSON responses (never silently omitted for empty slices).
Invariant 2 (CAUSALITY) compliance.

## Tests Run

### site/graph — Knowledge tests (6/6 PASS)

```
=== RUN   TestKnowledgePublic                     PASS
=== RUN   TestKnowledgeAuthed                     PASS
=== RUN   TestAssertOpReturnsCauses               PASS
=== RUN   TestKnowledgeClaimsCausesFieldPresent   PASS
=== RUN   TestKnowledgeMissingSpace               PASS
=== RUN   TestKnowledgeClaims                     PASS
ok  github.com/lovyou-ai/site/graph  0.246s
```

### cmd/post — Causes/assert tests (14/14 PASS)

```
=== RUN   TestAssertScoutGapCreatesClaimNode      PASS
=== RUN   TestAssertScoutGapMissingFile           PASS
=== RUN   TestAssertScoutGapNoGapLine             PASS
=== RUN   TestAssertScoutGapAPIError              PASS
=== RUN   TestAssertScoutGapSendsAuthHeader       PASS
=== RUN   TestAssertCritiqueCreatesClaimNode      PASS
=== RUN   TestAssertCritiqueMissingFile           PASS
=== RUN   TestAssertLatestReflectionCreatesDocument PASS
=== RUN   TestAssertLatestReflectionMissingFile   PASS
=== RUN   TestAssertCritiqueSendsCauses           PASS
=== RUN   TestAssertScoutGapSendsCauses           PASS
=== RUN   TestAssertLatestReflectionSendsCauses   PASS
=== RUN   TestAssertCauseIDsMultipleJoined        PASS
=== RUN   TestAssertCritiqueNoTitle               PASS
ok  github.com/lovyou-ai/hive/cmd/post  0.605s
```

## Coverage Notes

- `TestKnowledgeClaimsCausesFieldPresent` verifies the fix directly: decodes JSON as `map[string]json.RawMessage` and asserts `causes` key is present even when no causes are declared. This is the key regression test.
- `TestAssertOpReturnsCauses` covers the round-trip: assert with a cause → causes stored → causes returned in both op response and GET /knowledge.
- All 14 cmd/post tests cover the `causes` payload being sent to the API on scout/critique/reflection assertions.

## Pre-existing Failure (not this iteration)

`TestReposts` in `site/graph/store_test.go:899` — nil pointer dereference. The test creates a space with hardcoded slug `"repost-test"` without checking the error; `sp` is nil when the slug already exists. Last change to store_test.go was iteration 350. **Not introduced by this fix.**

## Verdict

**PASS.** The fix is correct and tested. The `causes` key is now always present in JSON responses for claim nodes. Invariant 2 compliance restored.
