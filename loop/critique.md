# Critique: [hive:builder] Scout phase must assert gap as claim node � currently writes only flat file

**Verdict:** PASS

**Summary:** **Derivation chain:**
- Gap: Scout gaps lost on every iteration — `scout.md` overwritten, no graph node
- Fix: `assertScoutGap()` reads `scout.md`, parses `**Gap:**` line + iteration number, POSTs `op=assert` to `/app/hive/op`
- Tests: 6 tests total — 2 parse helpers (table-driven), 1 HTTP integration, 1 missing file, 1 no-gap-line, 1 API error

**Implementation check:**

`extractGapTitle` (line 375–385): scans line by line, checks `HasPrefix("**Gap:**")`, strips prefix, trims space. Correct for all 3 test cases including extra whitespace. ✓

`extractIterationFromScout` (line 364–372): regex `Iteration (\d+)`, returns first match or "unknown". Simple and correct. ✓

`assertScoutGap` (line 325–362): reads relative path `loop/scout.md` (consistent with how other functions in this file work), validates gap is non-empty, builds body with iteration + gap, POSTs `op=assert`. Error handling consistent with `syncClaims` pattern. ✓

`pipeline_state.go`: `time` import added — build confirms it compiles. ✓

**Invariant 11:** `assertScoutGap` passes `gapTitle` as the claim title/body — a display string, not an ID used for lookup. No JOIN on names. ✓

**Invariant 12:** Every new code path tested. The diff introduces 4 new functions; the test file covers all four, including error paths (missing file, no gap line, API 4xx). Build passes, all tests green. ✓

**One observation:** The `assert` op posts without `"kind": "claim"` in the payload. Whether the server infers `KindClaim` from the `assert` op or requires explicit `kind` is a server concern outside this diff — consistent with how other operations work in this codebase (e.g. `createTask` also omits `kind`). Non-blocking.

VERDICT: PASS
