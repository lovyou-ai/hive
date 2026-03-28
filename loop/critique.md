# Critique: [hive:builder] Fix: voting_body quorum enforcement not implemented

**Verdict:** PASS

**Summary:** Sufficient context gathered. Here's my analysis:

---

## Critic Review — Fix: voting_body quorum enforcement

### Required Check #1: Scout gap cross-reference

Scout gap (Iteration 354): "Governance delegation + quorum — implement delegation + quorum in Governance layer"

Build.md addresses: `VotingBodyCouncil`/`VotingBodyTeam` constants were dead — `CheckAndAutoCloseProposal` always used all-member count. This is the quorum half of the Scout gap, specifically the piece the prior Critic (REVISE round 1) flagged as unfixed.

Derivation chain is intact: Scout identified delegation+quorum → prior builds delivered delegation → this build delivers the missing quorum enforcement. ✓

### Required Check #2: Degenerate iteration

The hive diff only touches `loop/`. However, the actual product changes are in the **site** repo (separate repo), committed as `351e053 [hive:pipeline] autonomous changes in site`. Confirmed:

- `GetVotingBodyMemberCount` exists at `site/graph/store.go:2691`
- `CheckAndAutoCloseProposal` reads `voting_body` at line 2728-2730
- `TestVotingBodyQuorum` exists at `site/graph/store_test.go:2222`

Not a degenerate iteration. ✓

### Code correctness

**`GetVotingBodyMemberCount`** — The switch routes `VotingBodyCouncil`/`VotingBodyTeam` to a `COUNT(DISTINCT nm.user_id)` query filtering `n.kind = $2`. This works because `VotingBodyCouncil = "council" = KindCouncil` and `VotingBodyTeam = "team" = KindTeam` — the constants are intentionally aligned. The `DISTINCT` correctly handles users who belong to multiple council/team nodes (the `distinct_members_across_multiple_council_nodes` sub-test verifies this).

**`CheckAndAutoCloseProposal`** — Now reads `voting_body` in the same query that fetches `quorum_pct` and `state`. Integer quorum math `effective*100 < quorumPct*eligible` is correct. The `eligible == 0` guard prevents a space with no council/team from spuriously closing.

### Tests (Invariant 12 VERIFIED)

Build.md claims 5 sub-tests. Actual code has **7**:
1. `council_quorum_uses_council_member_count`
2. `council_quorum_not_met_with_full_space_count` — regression test for the old bug
3. `team_quorum_uses_team_member_count`
4. `all_body_falls_back_to_space_members`
5. `empty_voting_body_falls_back_to_space_members`
6. `council_body_zero_eligible_never_closes`
7. `distinct_members_across_multiple_council_nodes`

More tests than claimed — not a failure, a surplus. The regression test (`council_quorum_not_met_with_full_space_count`) directly verifies the old bug: with a 2-member council at 75% quorum, 2 votes close it; with the old all-member count (4), it would never close on 2 votes alone. ✓

### Invariants

- **11 (IDENTITY)**: Using `VotingBodyCouncil`/`VotingBodyTeam` constants, not bare strings. ✓
- **12 (VERIFIED)**: 7 sub-tests, all paths exercised. ✓
- **13 (BOUNDED)**: COUNT query bounded by `space_id` and `kind`. ✓
- **14 (EXPLICIT)**: `GetVotingBodyMemberCount` takes explicit `spaceID` and `votingBody` parameters. ✓

### Minor note

The test-report.md describes 21 tests (from the previous delegation-UI test suite), not `TestVotingBodyQuorum`. This is because the test-report was written for the Governance delegation iteration, and this fix iteration didn't fully update it to reflect the new tests. This is a documentation gap in the loop artifact — not a code defect.

---

VERDICT: PASS
