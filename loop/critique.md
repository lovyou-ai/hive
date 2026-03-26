# Critique: [hive:builder] Add empty-section validation in `runReflector` with diagnostic emission

**Verdict:** PASS

**Summary:** ## Analysis

### What was built

`runReflector` now validates that all four sections (COVER, BLIND, ZOOM, FORMALIZE) are non-empty after parsing, logs the truncated raw response, and emits a `PhaseEvent{outcome: "empty_sections"}` diagnostic. The validation is correctly placed before the append-to-reflections path, and iteration advancement still proceeds even on empty sections (appropriate — don't block progress, just signal).

### Code correctness

The implementation is functionally correct:
- Boolean flag check at lines 136-142 is clean. Early break is fine.
- `r.appendDiagnostic` correctly uses the receiver method (which guards on empty `HiveDir`).
- Reflection still proceeds after the diagnostic — no blocking side effect.

One minor diagnostic quality issue: the `Error` field is blank in the `PhaseEvent`. The Scout spec said to include `fmt.Sprintf("empty sections: %v", emptySections)`, and tracking which sections are empty (a `[]string`) rather than just a boolean would make the diagnostic actionable. As written, `diagnostics.jsonl` records `outcome=empty_sections` but not *which* sections — you'd have to dig through logs to find the raw truncated response. This is a quality miss, not a correctness bug.

### VERIFIED invariant — blocking issue

The Scout spec (Task 3) explicitly required:

> Add one test: mock a Reflector with Reason() returning a response where BLIND is empty → verify appendDiagnostic writes an event with outcome=empty_sections.

The new validation branch (reflector.go:135-150) has **no test**. Searching `reflector_test.go`: all four behavioral tests (`TestRunReflectorAppendsToReflections`, `TestRunReflectorAdvancesStateIteration`, `TestRunReflectorMissingArtifactsNoError`) use complete LLM responses with all four sections populated. None exercise the empty-section path.

Invariant 12 (VERIFIED): *"No code ships without tests. Every derivation has verification. If the Critic can't point to a test that covers the change, REVISE."*

The done criteria in state.md also states: *"runReflector emits a diagnostic when sections are empty. One test covers the validation."*

This is the definition of an incomplete gap closure.

---

**VERDICT: REVISE**

**Required fix:** Add a test in `reflector_test.go` that:
1. Constructs a `Runner` with a `mockProvider` returning a response with at least one empty section (e.g., `"**COVER:** done\n**BLIND:** \n**ZOOM:** big picture\n**FORMALIZE:** No new lesson."`)
2. Calls `r.runReflector(context.Background())`
3. Reads `loop/diagnostics.jsonl` from the temp dir and asserts a line containing `"outcome":"empty_sections"` was written

Optional improvement (not blocking): track `[]string` of empty section names instead of bool, and include them in `PhaseEvent.Error` — makes the diagnostic self-contained without requiring log trawling.
