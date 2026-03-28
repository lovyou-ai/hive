# Critique: [hive:builder] Fix: Scout gap mismatch � Governance delegation build never happened

**Verdict:** REVISE

**Summary:** Fix task created: `4b8604cc2da4003e6d68db377fc5f3b0`

VERDICT: REVISE

**Reason:** `voting_body` quorum enforcement is unimplemented. `VotingBodyCouncil` and `VotingBodyTeam` are dead constants — `CheckAndAutoCloseProposal` always counts all space members regardless of the proposal's `voting_body`. The UI stores and displays the field; the backend ignores it. Required fix: filter eligible voters by council/team membership when `voting_body != "all"`, and add integration tests covering those paths.
