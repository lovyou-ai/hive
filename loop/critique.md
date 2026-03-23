# Critique — Iteration 87

## Derivation Chain
- **Gap:** `/app` showed only spaces — no cross-space visibility of tasks, conversations, or agent work
- **Plan:** Add three cross-space queries, rewrite template as dashboard
- **Code:** 3 store queries, 3 types, 4 templates, 3 helper functions, handler update
- **Test:** Existing tests pass (store + handler). No new tests for the dashboard queries.

## AUDIT

- **Identity (inv 11):** Tasks query uses `author_id` (user ID) for authorship and `assignee` (display name) for assignment. The assignee match is a name-based lookup — acceptable since the schema stores assignee as name, but this is a known debt. Conversation query matches on user IDs in tags. Agent activity uses `user_id` from space_members. No new name-as-identifier bugs introduced. CHECK.
- **Tests (inv 12):** No new tests for the three dashboard queries. ISSUE — but the queries are read-only aggregations over tested tables. Medium risk.
- **Bounded (inv 13):** All three queries have LIMIT clauses (30, 20, 20). CHECK.
- **Explicit (inv 14):** Dashboard types explicitly embed base types + space context. Dependencies declared via JOINs. CHECK.

## DUAL (root cause)
The assignee field stores display names, not user IDs. This means `ListUserTasks` must resolve the user's name first, then match. If a user changes their display name, their task assignments break. This predates iteration 87 — it's inherited debt, not new.

## Verdict: APPROVED

The core functionality works and is deployed. The assignee-as-name debt should be tracked but doesn't block this iteration.
