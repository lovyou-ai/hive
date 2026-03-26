Based on my analysis of the codebase, state, and recent commits, here's the Scout gap report:

---

## GAP REPORT

**Gap:** The Reflector receives insufficient content to produce meaningful reflections, causing empty COVER/BLIND/ZOOM/FORMALIZE stubs to persist across iterations.

**Evidence:**

1. **Code inspection** — `writeBuildArtifact` in `pkg/runner/runner.go:405-413` writes only 4 lines to `loop/build.md`: commit hash, cost, and timestamp. No commit message. No diff stats.

2. **Recent artifacts** — `loop/build.md` from the most recent builder iteration contains:
   ```
   # Build: Tests for critique artifact write and 5-phase tree
   - **Commit:** 0a354ac60a665d0fae3f57df714a78ae6607b664
   - **Cost:** $0.2278
   - **Timestamp:** 2026-03-26T20:33:32Z
   ```
   That's it. No task context. No changed files.

3. **Reflector validation gap** — `runReflector` in `pkg/runner/reflector.go:100-156` has no validation after parsing. It writes whatever sections come back from the LLM, even if empty. No diagnostic when BLIND or FORMALIZE are blank.

4. **Critic flags it twice** — `loop/critique.md` from iteration 315 explicitly states: *"the empty reflection stubs are now three iterations deep — the Reflector is writing structure without substance, which defeats the artifact's purpose."*

5. **Missing helpers** — `gitSubject()` and `gitDiffStat()` don't exist. Only `gitHash()` exists.

**Impact:** 

The Reflector's meta-learning (the entire purpose of phase 5) breaks. Without substance in build.md, the LLM has no basis for BLIND (what's missing from the build) or ZOOM (how the iteration connects to larger patterns). Empty reflections compound: no new lessons are extracted, so the loop can't see its own blind spots. This defeats the self-correcting feedback loop that the 5-phase pipeline is designed to create.

**Scope:**

- `pkg/runner/runner.go:405-413` — `writeBuildArtifact` function
- `pkg/runner/runner.go:415-424` — `gitHash()` helper (reference)
- `pkg/runner/reflector.go:100-156` — `runReflector` function  
- `pkg/runner/reflector.go:132-143` — section parsing and validation
- `pkg/runner/reflector_test.go` — test file
- `loop/reflections.md` — where reflections are appended (artifact validation)

**Suggestion:**

Enrich the build artifact so the Reflector has context. Add commit message subject and diff stat to `writeBuildArtifact`. Add two helper methods (`gitSubject()`, `gitDiffStat()`). Add validation in `runReflector` to catch empty sections before writing reflections, emit diagnostics when they occur. Include one test to verify empty-section detection. Optionally clean up stale directives from state.md that were replaced in iteration 242+ (reduces LLM context by ~3000 tokens for Scout/PM calls).

---

**Definition of done:** `build.md` contains commit message + diff stat + task description. `runReflector` validates section completeness and emits diagnostics on empty sections. One test covers the validation. Reflections.md entries have substance.