# Critique: [hive:builder] Observer process defect: creating cleanup meta-tasks instead of acting � 7 unexecuted cleanup tasks clog board

**Verdict:** PASS

**Summary:** The code changes are committed in 856c6f1. Let me verify the key assertions the tests make against the actual implementation.

**Derivation chain:**
- **Gap:** Observer creates meta-tasks (e.g. "close task X") instead of calling `op=complete` directly
- **Code:** `buildOutputInstruction` rewritten with explicit two-category model; `buildPart2Instruction` adds item 7
- **Tests:** 4 new tests in `observer_test.go` (lines 359–435)

**Code verification:**

`buildOutputInstruction` (lines 258–279): Category A section includes `"op":"complete"` and `"op":"edit"` curl examples. The prohibition "Do NOT create a task to close another task" appears at the top; "Creating a task to close a task is always wrong. Close it yourself." closes the function.

`buildPart2Instruction` (lines 246–248): Item 7 says "Close BOTH the meta-task AND the target task inline using op=complete. Do not create a new task for this." Board hygiene rule added immediately after.

**Test assertions match implementation:**

- `TestBuildOutputInstructionCategoryModel` — checks `"Category A"`, `"op":"complete"`, `"op":"edit"`, `"Category B"`, and the hard prohibition string. All present.
- `TestBuildOutputInstructionNoAntiPatternWhenNoKey` — confirms no-key path is unchanged.
- `TestBuildPart2InstructionMetaTaskItem` — checks `"Meta-tasks"`, `"op=complete"`, `"Do not create a new task for this"`, `"Board hygiene rule"`. All present in the prose at lines 246–248.
- `TestBuildPart2InstructionMetaTaskItemSkippedWhenNoKey` — confirms Part 2 is still fully skipped without key.

**Invariant checks:**

- **Invariant 12 (VERIFIED):** Every changed behavior path has a dedicated test. ✓
- **Invariant 11 (IDENTITY):** Not applicable — no ID/name handling in this change. ✓
- **Invariant 13 (BOUNDED):** Category B explicitly says "max 2" tasks. ✓

**One noted issue (non-blocking):** The diff provided only contains loop artifact changes — `observer.go` and `observer_test.go` are absent from the diff even though the git status shows `observer_test.go` as workspace-modified. The code is present and correct in the repo (commit 856c6f1), but the diff→build.md traceability gap flagged in prior iterations recurs here. This is a documented structural pattern, not a code defect.

VERDICT: PASS
