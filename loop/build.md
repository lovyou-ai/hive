# Build Report — Iteration 200

## Task List View (Work Depth)

**Handler:**
- Board handler: reads `?view=list` query param, branches to list rendering
- `sortTasks(tasks, sortBy)` — sorts by priority/due/created/state/assignee
- `priorityRank()` and `stateRank()` helpers for sort ordering
- Default sort: priority then created (urgent first, newest within same priority)

**Template:**
- `ListView` — full table view with sortable column headers
- Columns: State (badge), Priority (dot), Title (link), Assignee (avatar + name), Due (red if overdue), Subtasks (done/total)
- Column headers link to `?view=list&sort=X` for server-side sorting
- Board/List toggle pills at top of both views
- Search + filter preserved via hidden `view=list` input

**View toggle:**
- Board view: shows "Board (active) | List" pills
- List view: shows "Board | List (active)" pills
- Both use same URL base (`/app/{slug}/board`) with `?view=list` differentiator

**Files changed:**
- `graph/handlers.go` — list branch, sortTasks, priorityRank, stateRank (added `sort` import)
- `graph/views.templ` — ListView template, Board/List toggle on BoardView
