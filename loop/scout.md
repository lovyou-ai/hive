# Scout Report — Iteration 206

## Gap: No Goal entity — projects have no "why"

**Source:** unified-spec.md, work-general-spec.md. Goal is the Plan mode's primary entity.

**Current state:** Projects contain tasks. But there's no higher-level entity that says WHY a project exists. No OKRs, no objectives, no measurable outcomes. Work has execution (Board) and grouping (Projects) but no direction (Goals).

**What's needed:**
1. `goal` node kind
2. Goals page: `/app/{slug}/goals` — list goals with child project/task counts, progress
3. Goals can contain projects and tasks (via parent_id)
4. Sidebar link

**Why Goal:** It activates the Plan mode. Goal → Project → Task is the hierarchy every org needs. "Ship v1" → "Auth project" → "Implement OAuth" → done. Without goals, projects are just folders. With goals, projects have purpose.

**Approach:** Exact same pattern as Projects (iter 205): 1 constant, 1 handler, 1 template. `intend` op already accepts `kind` param — add `goal` to the allowlist.
