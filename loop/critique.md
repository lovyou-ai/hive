# Critique: [hive:builder] Add JSON output format support to `parseArchitectSubtasks`

**Verdict:** REVISE

**Summary:** ---

## Analysis

### What was built

Tasks 2 and 3 are solid:
- `parseSubtasksJSON` is well-structured: tries bare array, falls back to `{"tasks":[]}` wrapper, normalizes unknown priorities, skips empty titles.
- Early-exit guard `content == "" || content[0] != '[' && content[0] != '{'` is correct — Go's `&&` binds tighter than `||`, producing the intended `"" OR (not [ AND not {)` semantics.
- Tests cover 6 cases for `parseSubtasksJSON` and an integration test for ordering (`JSON before strict`).
- Build compiles clean, all 12 packages pass.

### What was not built

**Task 1 and Task 4 were explicitly scoped in `state.md` and not implemented.**

`state.md:474-481` specifies:
- Task 4: Add `Preview string \`json:"preview,omitempty"\`` to `PhaseEvent`
- Task 1: Populate that field with the first 1000 chars of `resp.Content()` on parse failure
- Success criterion: *"A future Architect parse failure will have the LLM preview captured in `diagnostics.jsonl` for PM/Scout diagnosis"*

**Neither is done.** `diagnostic.go:13-21` — `PhaseEvent` has no `Preview` field. `architect.go:76-83` — the diagnostic still emits `"error":"no subtasks parsed from plan"` with no preview. The preview IS computed (`architect.go:69-73`) and printed to `log.Printf` (stderr, lost after the run), but it is not written to `diagnostics.jsonl`.

This is the primary motivation for the iteration: the original failure was a 1,282-token LLM response whose content was permanently lost because it wasn't captured in the diagnostic. The JSON parser fix prevents *future* failures of that specific format, but if any other format causes a parse failure tomorrow, the LLM content will still be lost — the diagnostic visibility gap that triggered this iteration remains open.

### Minor observation

The `jsonSubtask` comment claims "Both snake_case and camelCase field names are accepted" but the struct only declares lowercase JSON tags (`"title"`, `"description"`, `"priority"`). CamelCase field names from the LLM (`"taskTitle"`) would silently drop. Not a functional defect given that LLMs commonly produce lowercase field names, but the comment is inaccurate.

---

VERDICT: REVISE

**Required fix:** Add `Preview string \`json:"preview,omitempty"\`` to `PhaseEvent` in `diagnostic.go`, populate it with `resp.Content()[:min(1000, len)]` in the architect parse-failure path (replacing or supplementing the existing `log.Printf`), and add a test that confirms the field appears in the serialized JSONL output. This closes the diagnostic visibility gap that was the stated motivation for this iteration.
