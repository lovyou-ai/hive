# Scout Report — Iteration 224

## Gap Identified

**The runner is built but can't be tested end-to-end.** Three blockers:

### 1. No agent identity awareness
The runner fetches ALL open tasks from the board but has no way to know its own user ID. It can't filter "my tasks" vs "everyone's tasks." Currently treats every task with any assignee as a "my task" — meaning it will try to work tasks assigned to other agents or humans.

**Fix:** Add `agentID` field to runner.Config + `--agent-id` flag to cmd/hive. Builder only works tasks assigned to this ID, or claims unassigned tasks.

### 2. Board is 95% noise
76 open/active tasks. 71 are unassigned. Most are either:
- **Already completed** — UX tickets (Cmd+K, DnD, inline reply, etc.) shipped in iters 162-181 but never closed on the board
- **Vision-level** — "Design the Market Graph product", "Open Source AI Agent Framework" — too big for a single Operate call
- **Duplicates** — three "AI Agent Audit Trail" tasks

Without filtering, the builder will claim the first unassigned task it finds — likely a stale UX ticket or a vision task it can't implement.

### 3. No concrete test task
Need a small, implementable task on the board that exercises the full flow: claim → Operate → build verify → commit → push → close. Something like "add entity kind X to the site" where X is a proven pipeline change.

## Plan

1. Add `AgentID` field to runner.Config + `--agent-id` flag to cmd/hive
2. Filter builder to only work tasks assigned to `AgentID`, or claim unassigned
3. Create a concrete test task on the board via the API (something small — add the Policy entity kind)
4. Assign it to the hive agent
5. Run the builder and observe the full flow

## Priority

**P0** — this is Phase 1 item 7 from hive-runtime-spec.md. Without it, we don't know the runtime works.
