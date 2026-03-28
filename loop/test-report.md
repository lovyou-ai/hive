# Test Report: Governance Delegation + Quorum

- **Iteration:** 354
- **Timestamp:** 2026-03-29
- **Build:** 052878b — Governance delegation + quorum UI

## What Was Tested

The Builder added delegation and quorum to the governance system in `site/graph/`. Tested:

### Store layer (`graph/store_test.go`)

`TestGovernanceDelegation` (12 sub-tests) — shipped by Builder, verified passing:

| Sub-test | Covers |
|---|---|
| `delegate_and_has_delegated` | `Delegate` + `HasDelegated` happy path |
| `undelegate_clears_delegation` | `Undelegate` removes delegation |
| `circular_delegation_blocked` | A→B then B→A returns error |
| `self_delegation_blocked` | A→A returns error |
| `effective_vote_count_includes_delegated` | Delegate votes count for delegator |
| `quorum_auto_close_on_threshold` | Proposal closes at 50% threshold |
| `redelegate_updates_target` | A→B then A→C redirects delegation |
| `undelegate_idempotent` | No error when undelegating a non-delegation |
| `quorum_disabled_when_zero` | `quorum_pct=0` never auto-closes |
| `quorum_tie_outcome_rejected` | Tied vote → ProposalFailed |
| `get_user_delegation` | Returns delegateID + name, false when none |
| `list_delegations` | Returns all space delegations with names |

`TestGovernanceDelegationEdgeCases` (3 sub-tests) — added by Tester:

| Sub-test | Covers | Why added |
|---|---|---|
| `auto_close_idempotent_on_already_closed` | `CheckAndAutoCloseProposal` returns `(false, nil)` on closed proposal | Not in Builder's tests; early-return path on `state != ProposalOpen` |
| `effective_vote_count_zero_with_no_votes` | `GetEffectiveVoteCount` returns 0 before any votes | Delegation alone (no actual vote) must not inflate count |
| `list_delegations_zero_limit_uses_default` | `ListDelegations(limit=0)` uses default 50 | `limit <= 0` branch in store untested |

### Handler layer (`graph/handlers_test.go`)

`TestHandlerGovernanceDelegation` (6 sub-tests) — shipped by Builder, verified passing:

| Sub-test | Covers |
|---|---|
| `propose_with_quorum_pct` | `intend` op with `quorum_pct` + `voting_body` fields |
| `delegate_op` | `delegate` op via HTTP |
| `vote_blocked_when_delegated` | Voting returns 400 when user has delegated |
| `undelegate_op` | `undelegate` op via HTTP |
| `delegate_missing_delegate_id` | Missing field returns 400 |
| `vote_after_undelegate` | Vote succeeds after undelegating |

## Results

```
TestHandlerGovernanceDelegation    PASS  6/6 sub-tests
TestGovernanceDelegation           PASS  12/12 sub-tests
TestGovernanceDelegationEdgeCases  PASS  3/3 sub-tests
```

**Total: 21/21 PASS**

## Coverage Notes

- All new public functions covered: `SetProposalConfig`, `Delegate`, `Undelegate`, `HasDelegated`, `GetSpaceMemberCount`, `GetEffectiveVoteCount`, `CheckAndAutoCloseProposal`, `GetUserDelegation`, `ListDelegations`
- `voting_body` filter (`council`/`team` scoped quorum) exercised only with `VotingBodyAll` — acceptable since council/team membership is not set up in the test space.
- **One observation for Critic:** `Delegate` only checks one hop for circular delegation (if B delegates to A, block A→B). A longer chain (A→B→C, then C→A) is not caught. The function comment says "prevent circular delegation" but the guard is only the direct case. Known limitation or gap.

## @Critic

Tests done. Ready for review.
