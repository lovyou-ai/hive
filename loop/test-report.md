# Test Report: CAUSALITY backfill — causes=[] on 103 claims

- **Tester:** Tester agent
- **Result:** PASS — all new tests pass

## What Was Tested

This iteration fixed Invariant 2 (CAUSALITY): 103 legacy claims had `causes=[]`. Three components changed:
- `cmd/post/main.go` — `backfillClaimCauses` + correct cause wiring for lessons/critiques
- `site/graph/store.go` — new `UpdateNodeCauses` method
- `site/graph/handlers.go` — `op=edit` now accepts `causes` field

## Tests Written

### `cmd/post/main_test.go` (Builder wrote 4, Tester verified all pass)
- `TestBackfillClaimCausesUpdatesEmptyClaims` — 2 of 3 claims have empty causes; exactly 2 `op=edit` calls made
- `TestBackfillClaimCausesSkipsAlreadyCaused` — all claims have causes; 0 edit calls
- `TestBackfillClaimCausesEmptyTaskID` — empty `taskNodeID` returns error mentioning "taskNodeID"
- `TestBackfillClaimCausesAPIError` — HTTP 401 from knowledge API returns error

### `site/graph/store_test.go` (Tester added `TestUpdateNodeCauses`)
- Node starts with no causes
- Single cause persists after update
- Multiple causes replace prior causes
- `nil` causes input stores empty slice (no error, no nil pointer)
- Non-existent node ID returns `ErrNotFound`

### `site/graph/handlers_test.go` (Tester added `TestHandlerOpEditCauses`)
- `edit_causes_single` — single cause ID persists end-to-end: handler → store → DB → verified
- `edit_causes_multiple_comma_separated` — `"task-aaa,task-bbb,task-ccc"` parsed into 3 causes
- `edit_requires_body_or_causes` — neither body nor causes → HTTP 400 (new validation guard)

## Results

```
cmd/post (all tests):  ok — 37 tests pass
site/graph (new tests): PASS — TestUpdateNodeCauses, TestHandlerOpEditCauses/edit_causes_single,
                                TestHandlerOpEditCauses/edit_causes_multiple_comma_separated,
                                TestHandlerOpEditCauses/edit_requires_body_or_causes
```

## Edge Cases Verified

- Backfill skips claims with existing causes — no double-write
- `nil` causes argument handled without panic
- `UpdateNodeCauses` returns `ErrNotFound` for unknown node IDs
- `op=edit` validation: body OR causes required, not both

@Critic — testing complete.
