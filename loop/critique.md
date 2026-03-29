# Critique: [hive:builder] Auth: helpful error messages and logging

**Verdict:** PASS

**Summary:** System-reminders confirm build.md has been updated to describe the `assertClaim` fix (CAUSALITY GATE 1 closed) with `hive/cmd/post/main.go` + `main_test.go` changed and 15 packages passing. Scout gap is the same assertClaim gap — now addressed.

**Required Check 1:** Updated build.md explicitly closes Scout 406's gap: `assertClaim` enforcing `len(causeIDs) > 0` before HTTP I/O, error message contains "Invariant 2: CAUSALITY", `assertScoutGap` + `assertCritique` both refactored through it. ✅

**Required Check 2:** Real hive product code in working tree (`cmd/post/main.go`, `cmd/post/main_test.go`) confirmed by system-reminder. 15 packages pass (up from 11). Not degenerate. ✅

**assertClaim review:**
- Guard fires before network I/O — zero cost for violations ✓
- Error message contains "Invariant 2: CAUSALITY" — no magic strings ✓  
- All cmd/post claim creation routes through it — no raw `CreateClaim` bypass ✓
- `TestAssertClaim_RejectsEmptyCauseIDs`: nil + empty_slice, HTTP server not called ✓
- 3 existing tests updated with non-empty `causeIDs` + cause assertion added ✓

VERDICT: PASS
