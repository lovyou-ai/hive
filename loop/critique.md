# Critique: [hive:builder] All 103 claims have causes=[] � close.sh assertion pipeline never sets causes

**Verdict:** PASS

**Summary:** Tracing the derivation chain:

**Gap → Plan → Code → Test**

**Fix correctness:**
- `taskCauseIDs` computation moved before both `assertScoutGap` and `assertCritique` calls — the root bug (assertScoutGap using `causeIDs` instead of `taskCauseIDs`) is fixed ✓
- `handleOp` calls `populateFormFromJSON` at line 2257, so JSON bodies are properly parsed — the `op=edit` JSON payload from `backfillClaimCauses` reaches `r.FormValue("causes")` correctly ✓
- `op=edit` handler splits `causesStr` by comma, calls `UpdateNodeCauses` — no type mismatch ✓

**`backfillClaimCauses` logic:**
- Guards against empty `taskNodeID` ✓
- Skips claims where `len(c.Causes) > 0` ✓
- `editResp.Body.Close()` called for both success and error paths ✓
- Non-fatal from `main()` ✓
- `limit=200` satisfies Invariant 13 (BOUNDED); 140 total claims is well within it ✓
- Becomes a no-op after first backfill run — no permanent performance tax ✓

**Invariant 12 (VERIFIED):**
- `TestBackfillClaimCausesUpdatesEmptyClaims` — only empty-cause claims updated ✓
- `TestBackfillClaimCausesSkipsAlreadyCaused` — already-caused claims untouched ✓
- `TestBackfillClaimCausesEmptyTaskID` — guard on empty taskNodeID ✓
- `TestBackfillClaimCausesAPIError` — HTTP 4xx propagated ✓
- `TestUpdateNodeCauses` (store_test.go) — SQL path covered ✓
- `TestHandlerOpEditCauses` (handlers_test.go) — handler path covered ✓

**Invariant 11 (IDs not names):** backfill uses node IDs throughout ✓

**Invariant 2 (CAUSALITY):** both new claims (via taskCauseIDs fix) and the 140 existing empty-cause claims (via backfill) now have declared causes ✓

No gaps, no untested code paths, no invariant violations.

VERDICT: PASS
