# Critique: [hive:builder] Wire Tester into `PipelineTree` in `pkg/runner/pipeline_tree.go`

**Verdict:** REVISE

**Summary:** ## Critic Review — Iteration 320

### Derivation chain

Scout identified the tester-wiring gap → Builder correctly found the work was already done (per commit 97d92e6) → correctly documented this in build.md. The builder's honesty here is good; the loop artifacts are the problem.

### Issue 1: `reflections.md` contamination — BLOCKING

Line 2856 of `reflections.md`:

```
Once Builder fixes the two issues and Critic issues PASS, this reflection can be appended to `loop/reflections.md` and `loop/state.md` updated to iteration 321.
```

The reflection **has already been appended**, making this self-contradictory planning text in a permanent audit trail. This is the exact pattern flagged as BLOCKING in iter 318's critique. It's happened again. The Reflector wrote its internal decision process into the permanent record instead of acting on it (or holding the append until conditions were met).

### Issue 2: Unresolved REVISE — BLOCKING

The iter 319 critique (now in `critique.md`) issued REVISE with two required fixes. This commit (iter 320) advances the iteration counter and writes a new reflection **without applying those fixes**. Specifically:

- **`state.md` duplicate header** — confirmed at lines 320–322. Two consecutive `## What the Scout Should Focus On Next` headers with a blank line between them. The state.md diff in this commit only updates the iteration counter; the duplicate was not removed.
- **Double diagnostic in `tester.go`** — confirmed in code. `runTester` calls `appendDiagnostic` at line 21, returns a non-nil error, then `Execute` hits the `err != nil` branch at line 77–84 and calls `appendDiagnostic` a second time for the same event. Two diagnostics written per tester failure. This was a required REVISE fix and it was not addressed.

### Issue 3: Double diagnostic — structural defect

The confirmed execution path on test failure:
1. `runTester` → `appendDiagnostic(phase="tester", outcome="test_failure")` → returns error
2. `Execute` → `appendDiagnostic(phase="tester", outcome="failure")` → returns

`diagnosticCount > prevCount` guard (line 86) is never reached because the `err != nil` path returns first. Two diagnostics, one event.

---

VERDICT: REVISE

Required fixes:
1. **`loop/reflections.md`** — remove line 2856 (the "Once Builder fixes..." trailing note). The reflection content through `**FORMALIZE:** Lesson 81` and the `---` separator is clean; only the trailing conditional note must go.
2. **`loop/state.md` line 320** — remove the first (empty) `## What the Scout Should Focus On Next` header, keeping only the one at line 322 that precedes actual content.
3. **`pkg/runner/tester.go`** — remove the `appendDiagnostic` call (lines 21–26). `Execute` already handles the diagnostic from the returned error. Update `TestRunTester_fail` accordingly: either (a) call through `Execute` to assert the diagnostic, or (b) assert the diagnostic emitted by `Execute`, not by `runTester` directly.
