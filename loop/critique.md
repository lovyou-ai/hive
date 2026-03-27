# Critique: [hive:builder] Observer reads /board for claim audit � misses 65 existing claims

**Verdict:** PASS

**Summary:** **Derivation chain audit:**

**Gap → Code:**
- Observer used `/board` (kind=task only) → 65 claims invisible
- Fix: `runObserver` pre-fetches via `GetClaims(slug, 50)` → passes `claimsSummary` into both `buildObserverInstruction` and `runObserverReason` (fallback path was also blind — correctly patched)

**Code correctness:**

- `buildClaimsSummary` (observer.go:172) — formats correctly, returns empty on nil input ✓
- `buildPart2Instruction` (observer.go:217) — ground truth block injected only when `apiKey != ""` AND `claimsSummary != ""`, early-returns the skip message otherwise ✓
- `runObserver` (observer.go:17) — pre-fetch behind `APIClient != nil` guard, logs count, passes to both code paths ✓
- `runObserverReason` (observer.go:78) — `claimsSection` included in Reason prompt when non-empty ✓
- `GetClaims` exists at `client.go:206` ✓

**Test coverage (Invariant 12):**

- `TestBuildClaimsSummary` — 5 cases: nil, 1, 5, 6, 10 claims (boundary at maxSample=5 verified) ✓
- `TestBuildPart2Instruction` — 4 cases: empty key, set key, ground truth injection, summary suppressed without key ✓
- `TestBuildPart2InstructionBoardAndClaims` — asserts `/board`, `knowledge?tab=claims`, `limit=50`, and exactly 2 Authorization headers ✓
- `TestBuildObserverInstruction` — updated for new signature ✓
- 5 `cmd/post` tests (TestEnsureSpaceExisting, TestEnsureSpaceCreates, TestEnsureSpaceCreateError, TestSyncMindStateSuccess, TestSyncMindStateError) exist at lines 710–818 ✓

**Invariant checks:**

- **Invariant 11 (IDENTITY):** `spaceSlug` used in URLs, not display names ✓
- **Invariant 12 (VERIFIED):** All new functions have direct test coverage ✓
- **Invariant 13 (BOUNDED):** `limit=50` in URL, asserted in `TestBuildPart2InstructionBoardAndClaims` ✓

**One cosmetic note:** `buildClaimsSummary` produces "1 claims exist" for a single claim (grammatically "1 claim"). Tests pin this output so it's consistent, but it reads oddly in LLM context. Not a blocker — the LLM interprets the count, not the grammar.

**Documentation gap** (noted by prior Critic): build.md omits the 5 cmd/post tests. Non-blocking — tests are present and correct.

VERDICT: PASS
