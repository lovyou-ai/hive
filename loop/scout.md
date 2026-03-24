# Scout Report — Iteration 200

## Gap: Task List view (Work depth — Linear's default)

**Source:** work-spec.md — List view replaces "Linear board, Jira table." work-product-spec.md — "List — tasks sorted/filtered/grouped (table)."

**Current state:** Work has Board (kanban) and Dashboard (My Work). No table/list view for bulk task scanning. Board is visual but doesn't let you sort by priority or scan 50 tasks quickly.

**What's needed:**
1. Route: `GET /app/{slug}/list` — new lens
2. Handler: load tasks, support sort (priority, state, due_date, created_at) and filter (state, priority, assignee)
3. Template: compact table rows with state indicator, priority dot, title, assignee, due date, subtask count
4. Sidebar: add "List" link in space navigation

**Why this first:** List is the most-used view in Linear. Power users scan tasks by priority, not by kanban column. This is the #1 missing Work feature for Linear parity.

**Approach:** Follow the Board handler pattern. Reuse existing `ListNodes` with kind=task filters. Add sort parameter to query. Render as compact rows, not cards.

**From the spec:**
```
List — tasks sorted/filtered/grouped (table)
```

**Risk:** Low. New route + handler + template. No schema changes. No new store methods needed — existing ListNodes covers it.
