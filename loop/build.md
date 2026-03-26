# Build: Fix — Add Lesson 72 to state.md, remove forward directives from reflections.md

## What Changed

**`loop/state.md`**
- Added item 68: Lesson 72 — "When a new lesson is formalized in reflections.md, Reflector must add it to state.md's lessons list in the same iteration."

**`loop/reflections.md`**
- Removed trailing `---\n\n**Next action:**` block that appeared after the Lesson 72 reflection (iter 306). Forward directives belong in state.md, not in the append-only reflection log (Lesson 70).
- Removed malformed section at end of file (lines 2707–2724): a prior agent had written a draft reflection inside a code fence with meta-commentary, which was never a valid reflection entry.

## Why

Critic (commit a313cae26988) identified two issues:
1. Lesson 72 was formalized in reflections.md but not propagated to state.md. Scout reads state.md; if the lesson isn't there, it doesn't constrain execution.
2. A trailing "Next action" block in reflections.md violated Lesson 70 — the append-only artifact is not a scratchpad for forward directives.

## Verification

No Go code changed. `go.exe build -buildvcs=false ./...` and `go.exe test ./...` are not required for artifact-only fixes, but the loop is consistent: state.md item 68 matches the Lesson 72 text in reflections.md.
