# Critique — Iteration 90

## Verdict: APPROVED

- **Identity (inv 11):** Endorsements use from_id/to_id (user IDs). Profile resolved by name for URL but endorsement logic uses IDs. CHECK.
- **Tests (inv 12):** No tests for endorsement queries. ISSUE — consistent with recent pattern of adding tests in batches.
- **Bounded (inv 13):** ListEndorsers has LIMIT 20. CHECK.
- **Explicit (inv 14):** Endorsements table has explicit PK and index. CHECK.

## DUAL
Endorsements are separate from the ops table (which is space-scoped). This is architecturally correct — endorsements are user-to-user relationships, not space operations. They shouldn't be in the ops table because they don't belong to any space.
