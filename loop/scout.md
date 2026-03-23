# Scout Report — Iteration 88

## Gap: assignee stores display name, not user ID

The `nodes.assignee` column stores the display name (e.g., "Hive", "Matt") instead of a user ID. This is the last name-as-identifier bug from the identity fix (iter 48-49). The identity fix added `author_id` and `actor_id` but didn't address `assignee`.

Consequences:
- The dashboard query (`ListUserTasks`) has to resolve the user's name and match on it
- If a user changes their display name, their task assignments silently break
- Two-step lookup (resolve name → match) instead of single ID match
- Violates invariant 11 (IDENTITY)

## What "Filled" Looks Like

- Add `assignee_id` column to nodes
- Update all handlers that set assignee to also set assignee_id
- Update queries to match on assignee_id
- Keep `assignee` as display-only (like `author` alongside `author_id`)
- Backfill existing rows from users table

## Approach

Same pattern as the author_id migration (iter 48-49):
1. Add `assignee_id` column with default ''
2. Update handlers: intend, assign, claim ops set both assignee (name) and assignee_id (ID)
3. Update dashboard query to match on assignee_id
4. Update Mind triggers to use assignee_id
5. Backfill existing data
