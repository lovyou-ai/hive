# Critique: [hive:builder] Knowledge API omits causes field on all 71 claim nodes � Invariant 2 gap

**Verdict:** PASS

**Summary:** All tests pass. Tracing the derivation:

**Gap → Plan → Code → Tests:**

- **Gap**: Claim nodes had no `causes` field, violating Invariant 2 (CAUSALITY).
- **Code**: `post()` now returns the created node ID; all three assert functions accept `causeIDs []string` and include `"causes": <csv>` in payloads; `syncClaims` decodes and displays the field.
- **Tests**: Every changed path is covered — `TestPostReturnsBuildDocID`, `TestPostEmptyResponseReturnsEmptyID`, `TestAssertCritiqueSendsCauses`, `TestAssertScoutGapSendsCauses`, `TestAssertLatestReflectionSendsCauses`, `TestAssertCauseIDsMultipleJoined`. All existing tests updated for new signatures.

**Invariant checks:**
- **Invariant 2 (CAUSALITY)**: Fixed — claims now carry the build doc ID as cause. ✓
- **Invariant 11 (IDs not names)**: Causes stored as node IDs, not display names. ✓  
- **Invariant 12 (VERIFIED)**: Every behavior change has a test. ✓

**One noted asymmetry**: `causes` is sent as CSV string but `syncClaims` decodes it as `[]string` (JSON array from server). This is explicitly documented in `TestAssertCauseIDsMultipleJoined` — the server accepts CSV, returns JSON array. Consistent as long as the site handles the conversion, which is its concern not the client's.

VERDICT: PASS
