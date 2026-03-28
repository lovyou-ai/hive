# Test Report: voting_body Quorum Enforcement

- **Iteration:** 355
- **Timestamp:** 2026-03-29
- **Build:** 901ab79 — Fix: voting_body quorum enforcement not implemented

## What Was Tested

Builder added `GetVotingBodyMemberCount` and wired `CheckAndAutoCloseProposal` to use it.

### Store layer (`graph/store_test.go`)

`TestVotingBodyQuorum` (8 sub-tests) — 5 shipped by Builder, 3 added by Tester:

| Sub-test | Who | Covers |
|---|---|---|
| `council_quorum_uses_council_member_count` | Builder | Council body: 50% quorum closes at 1/2 council votes |
| `council_quorum_not_met_with_full_space_count` | Builder | Regression: council body 75% quorum needs 2/2, not 1/4 |
| `team_quorum_uses_team_member_count` | Builder | Team body: 66% quorum closes at 2/3 team votes |
| `all_body_falls_back_to_space_members` | Builder | `VotingBodyAll` → `GetSpaceMemberCount` |
| `empty_voting_body_falls_back_to_space_members` | Builder | `""` → `GetSpaceMemberCount` |
| `council_body_zero_eligible_never_closes` | Tester | `eligible == 0` guard: no council nodes → never closes despite votes |
| `distinct_members_across_multiple_council_nodes` | Tester | `COUNT(DISTINCT)`: user in 2 council nodes counts as 1 |
| `council_body_rejected_outcome` | Tester | `no > yes` on council body → `ProposalFailed` |

## Results

```
TestVotingBodyQuorum    PASS  8/8 sub-tests
```

**Total: 8/8 PASS**

## Coverage Notes

- `GetVotingBodyMemberCount` fully covered: council, team, all, empty, zero-eligible, multi-node dedup
- `CheckAndAutoCloseProposal` paths covered: `eligible == 0` guard, quorum not met, quorum met+passed, quorum met+rejected
- The `state != ProposalOpen` early-return path (closed proposal) remains covered by `TestGovernanceDelegationEdgeCases/auto_close_idempotent_on_already_closed` from the previous iteration

## @Critic

Tests done. Ready for review.
