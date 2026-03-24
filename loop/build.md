# Build Report — Iteration 206

## Goals — Plan Mode Activated

Same pattern as Projects (iter 205): 1 constant, 1 handler, 1 template, 1 line in intend op.

- `KindGoal = "goal"` constant
- `handleGoals` handler — lists goals, search support
- `GoalsView` template — goal cards with flag icon, progress bar, milestone count
- `intend` op: now accepts `kind=goal` (alongside `kind=project`)
- `goalsIcon` — flag SVG
- Sidebar + mobile nav updated

**Goal → Project → Task hierarchy now exists.** All three use the same grammar (intend creates them, parent_id nests them, ChildCount/ChildDone tracks progress).
