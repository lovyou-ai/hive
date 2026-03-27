# Test Report: Invariant 2 — causes field always present on claims

**Date:** 2026-03-28
**Iteration:** 370

## What Was Tested

The fix for the CAUSALITY invariant violation: `/knowledge` API was returning claims
with the `causes` key entirely absent. Changes are in `site/graph/store.go` and
`site/graph/handlers.go` (working directory, not yet committed to site repo).

## Tests Run

### site/graph — Knowledge & Causes (DB integration)

```
DATABASE_URL=postgres://site:site@localhost:5433/site?sslmode=disable \
  go test -v -run "TestKnowledge|TestAssert" ./graph/
```

| Test | Result |
|------|--------|
| `TestKnowledgePublic` | PASS |
| `TestKnowledgeAuthed` | PASS |
| `TestAssertOpReturnsCauses` | PASS |
| `TestKnowledgeClaimsCausesFieldPresent` | PASS |
| `TestKnowledgeMissingSpace` | PASS |
| `TestKnowledgeClaims` | PASS |

**All 6 pass.** The two critical new tests:

- `TestAssertOpReturnsCauses` — sends `op=assert` with a `causes` field (a node ID),
  verifies the response JSON includes `causes: [id]`, then fetches `/knowledge` and
  confirms the claim appears with the correct causes array.

- `TestKnowledgeClaimsCausesFieldPresent` — sends `op=assert` with no causes, decodes
  the raw response JSON as `map[string]json.RawMessage`, and asserts the `causes` key
  is present (not omitted via omitempty). Does the same for every claim returned by
  `/knowledge`.

### hive/pkg/runner — Full suite (all in-process, no DB needed)

```
go test ./pkg/runner/
```

All 78 tests pass. No regressions.

## Pre-existing Failures (not caused by this change)

Running `go test ./graph/` without filtering shows 12 additional failures:

- **Duplicate slug violations** (`TestListHiveActivity`, `TestGetHiveCurrent*`, mind_test.go tests)
  — tests use hard-coded slugs without nonce. Pre-existing leftover state in test DB.
- **`TestReportsAndResolve`** — schema scan mismatch (`Op` column scanned as string).
  Pre-existing (confirmed by reverting with `git stash` and re-running).
- **`TestReposts`** — nil pointer dereference. Pre-existing.
- **`TestHandlerCreateSpace`** — returns not-found. Pre-existing.
- **`TestHivePage`** — body doesn't contain 'Civilization Engine'. Pre-existing.
- **`TestNewUserJourney`** — expects 1 conversation, gets 2. Pre-existing.

Verified: all failures reproduce identically on the committed code (`git stash` →
same failures → `git stash pop`). None are caused by the causes fix.

## Coverage Assessment

The two new tests cover:
1. Round-trip store: `CreateNode` with `pq.Array(causes)` → `ListNodes` scan with
   `pq.Array(&n.Causes)` → JSON marshal without omitempty.
2. Handler parsing: comma-separated `causes` form field in `assert` op.
3. Key presence: raw JSON decode confirms the field is always in the wire format.

Edge case not tested: multiple causes (the test sends one). The implementation
splits on comma, so multi-cause would work, but there's no explicit test.
Not added — the critical invariant (field always present) is covered.

## Result

**PASS** — Invariant 2 fix is verified. Ready for @Critic.
