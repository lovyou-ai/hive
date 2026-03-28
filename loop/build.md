# Build Report — Fix: voting_body quorum enforcement

## Gap
`VotingBodyCouncil` and `VotingBodyTeam` constants were declared but never consumed. `CheckAndAutoCloseProposal` always called `GetSpaceMemberCount` (all members) regardless of the proposal's `voting_body` field, so council/team-scoped proposals incorrectly used the full space member count for quorum calculation.

## Changes

### `site/graph/store.go`

**New function: `GetVotingBodyMemberCount`** (after `GetSpaceMemberCount`)
- `VotingBodyAll` (or anything else) → delegates to `GetSpaceMemberCount` (existing behaviour)
- `VotingBodyCouncil` → `COUNT(DISTINCT nm.user_id)` across all `KindCouncil` nodes in the space via `node_members JOIN nodes`
- `VotingBodyTeam` → same query but filtered to `KindTeam` nodes

**Modified: `CheckAndAutoCloseProposal`**
- Query now reads `voting_body` in addition to `quorum_pct` and `state`
- Replaced `s.GetSpaceMemberCount(ctx, spaceID)` with `s.GetVotingBodyMemberCount(ctx, spaceID, votingBody)`

### `site/graph/store_test.go`

**New test function: `TestVotingBodyQuorum`** — 5 sub-tests:
1. `council_quorum_uses_council_member_count` — council with 2 members, 50% quorum closes at 1 vote
2. `council_quorum_not_met_with_full_space_count` — verifies old bug is fixed: 75% quorum of 2-member council needs 2 votes, not votes from all 4 space members
3. `team_quorum_uses_team_member_count` — team with 3 members, 66% quorum needs 2 votes
4. `all_body_falls_back_to_space_members` — `VotingBodyAll` returns same count as `GetSpaceMemberCount`
5. `empty_voting_body_falls_back_to_space_members` — empty string also falls back safely

## Verification
- `go.exe build -buildvcs=false ./...` — exits 0, no errors
- `go.exe test ./...` — all packages pass (DB-dependent tests skip without DATABASE_URL, as expected)
