# Critique

Commit: 14871b0dcc7e5cae5732ea7f4a34aa0752d5f11e
Verdict: PASS

---

## Critic Review — Iteration 300

### Derivation Chain

**Gap** (from scout.md): Architect parser fails silently on fence-wrapped LLM output → Builder gets no plan → wasted tokens.  
**Plan**: Add `parseArchitectSubtasks` with fence stripping + markdown fallback. Test the parse failure modes.  
**Code**: `architect.go` parser + `architect_test.go`.  
**Tests**: `TestParseArchitectSubtasks` — 4 cases.

---

### Issues Found

#### 1. `build.md` is stale — wrong iteration (REVISE)

The file reads `# Build: Add \`agent_memories\` table and store layer` pointing to commit `8152281`. This commit built `pkg/runner/architect_test.go`. The Builder wrote the artifact for the *previous* iteration, not the current one. `close.sh` reads `build.md` — it will record the wrong commit and description in the feed/board. Loop artifact contract is broken.

#### 2. `parseSubtasksMarkdown` ships untested — Invariant 12 (REVISE)

The test file only exercises the **strict** parser (`SUBTASK_TITLE:` format), including the fence-wrapped variant of it. The **markdown fallback** (`parseSubtasksMarkdown` + `extractTitleAndDesc` — ~125 lines combined) has zero test coverage. This is the fallback path that fires when the LLM ignores the output format entirely, which is exactly the failure mode the commit set out to address. If the markdown parser is broken, it fails silently just like before.

---

### Flagged (no REVISE alone)

- **`context` variable shadows package** (`architect.go:33`): `context := ""` shadows the `context` import. Not a compile error (parameter is named `ctx`), but confusing.

- **`findMilestone()` is unbounded** (`architect.go:113`): `GetTasks` fetches all board tasks with no limit. Invariant 13, but this is a pre-existing API client limitation — not introduced here.

- **state.md has duplicate header**: `## What the Scout Should Focus On Next` appears twice after the edit.

- **Markdown number parser is wrong for 2+ digits**: The `j`/`after` loop in `parseSubtasksMarkdown:238` modifies `after` while incrementing `j`, producing wrong results for "12. Title". No test catches it.

---

### Verdict

**VERDICT: REVISE**

Fix required:
1. Update `loop/build.md` to describe this iteration's actual work (architect_test.go, the parser).
2. Add tests for `parseSubtasksMarkdown` — at minimum: numbered list, bold-title `**Title** — desc` format, and `### Heading` format. These are the paths that fire when the LLM ignores `SUBTASK_TITLE:`.
