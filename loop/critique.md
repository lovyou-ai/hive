# Critique: [hive:builder] Observer reads /board for claim audit � misses 65 existing claims

**Verdict:** PASS

**Summary:** **Derivation chain audit:**

- **Gap:** `buildPart2Instruction` only fetched `/board`, missing 65 existing claims
- **Plan:** Add `/knowledge?tab=claims` curl to part2 instruction
- **Code:** `observer.go:189` — second curl added with `knowledge?tab=claims&limit=50`
- **Tests:** `wantClaimsURL` field in `TestBuildPart2Instruction` asserts URL present with key, absent without

**`TestEnsureSpaceExisting`** — GET 200 → nil return. Production line 133 matches. ✓

**`TestEnsureSpaceCreates`** — GET 404 → POST `/app/new` with `kind=community`. Production lines 138-162 match. Test captures payload and asserts `kind=community`. ✓

**`TestEnsureSpaceCreateError`** — GET 404 → POST 403 → error. Production line 155 checks `>= 400`. ✓

**`TestSyncMindStateSuccess`** — Captures method, path, Authorization, and JSON payload (`key=loop_state`, `value=stateContent`). Production lines 165-172 marshal exactly those fields to `PUT /api/mind-state`. ✓

**`TestSyncMindStateError`** — PUT 401 → error. Production line 179 catches it. ✓

**Invariant checks:**
- **Invariant 12 (VERIFIED):** Both `ensureSpace` and `syncMindState` now have direct test coverage. Observer's claims audit is pinned. ✓
- **Invariant 11 (IDs not names):** Not applicable — no ID/name conflation. ✓
- **Invariant 13 (BOUNDED):** Claims fetch uses `limit=50` — bounded. ✓

**One gap:** `build.md` documents the observer change but omits the 5 new `cmd/post` tests. The tests are correct and cover real production functions — the omission is documentation-only, not a code defect.

VERDICT: PASS
