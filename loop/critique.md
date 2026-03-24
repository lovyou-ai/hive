# Critique — Iteration 200

## Task List View: PASS

**Correctness:**
- Sort by priority uses rank ordering (urgent=0, high=1, medium=2, low=3, none=4). ✓
- Sort by due handles nil dates (nil sorts last). ✓
- Sort by state: active first, then review, open, done. Sensible. ✓
- Default sort: priority then created. Matches Linear's default. ✓
- Overdue highlighting in due column (red for past-due, non-done tasks). ✓
- Blocker indicator on title row. ✓

**Template:**
- Table headers are sortable links. Clean. ✓
- Board/List toggle is consistent on both views. ✓
- Search preserves view=list via hidden input. ✓
- Compact rows — more tasks visible than Board view. ✓

**NOTE:** Sort is server-side (page reload per sort change). For a table view, client-side sort would be snappier. Acceptable at current scale. Could HTMX-ify later.

**Tests:** No new tests. sortTasks is a pure function that could be tested, but the sort logic is straightforward.

## Verdict: PASS

**Milestone:** Iteration 200. 200 iterations shipped to production.
