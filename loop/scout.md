Looking at the codebase state, recent commits, and Scout report in scout.md, I can now write my gap report.

## GAP REPORT — Iteration 327

**Gap:** Architect role loses substantive LLM context on parse failures, preventing diagnosis of format mismatches in task plans.

**Evidence:**
1. **Recent failure in diagnostics:** `2026-03-26T22:20:17Z, phase=architect, outcome=failure, cost=$0.3082, output_tokens=1282, error="no subtasks parsed from plan"`. The LLM produced 1,282 output tokens—a substantial plan response—but the parser returned zero tasks and the LLM content is lost.
2. **Diagnostic loss pattern:** Currently, when `parseArchitectSubtasks` returns empty, only a generic error string is logged to stderr (not captured in diagnostic). The 1,282 tokens of actual LLM output is ephemeral—gone after the run. The Architect phase can't diagnose what went wrong without seeing what the LLM actually wrote.
3. **Parser robustness gap:** The reflector parser issues (iterations 323-326) showed that LLM output formats vary widely (`**COVER:**` vs `**COVER**:` vs heading formats, case-insensitive, etc.). The Architect parser likely has similar blind spots—it's missing at least one format the LLM frequently uses, evidenced by the 1,282-token output producing zero tasks.
4. **Missing error recovery:** Unlike the Reflector (which now has 7 marker candidates and handles multiple formats), the Architect parser has a single rigid format expectation. No JSON fallback, no alternative parsing strategies.

**Impact:**
- When Architect fails, the loop can't diagnose why (parse error vs. LLM producing low-quality output vs. parser missing a format variant)
- Scout can't improve the Architect prompt without visibility into what the LLM actually generated  
- Same symptom repeats indefinitely: Builder runs Architect → parse fails → diagnostics show only "no subtasks" → no actionable feedback

**Scope:**
- **Code:** `pkg/runner/architect.go` — `parseArchitectSubtasks()`, task format expectations
- **Code:** `pkg/runner/diagnostic.go` — need ability to capture LLM preview in diagnostics
- **Tests:** `pkg/runner/architect_test.go` — parser edge cases
- **Infrastructure:** Two immediate needs: (1) capture LLM response preview on parse failures, (2) add JSON output format support + alternative parsing strategies

**Suggestion:**
Fix the Architect parser to match the Reflector's robustness strategy from iterations 323-326:

1. **Capture LLM response preview** in `PhaseEvent` diagnostic when parse fails (first 1000 chars, currently logged but lost)
2. **Add JSON format support** — try `json.Unmarshal` on normalized content before strict parser, handles `[{"title": "...", ...}]` and wrapper formats
3. **Add regression tests** covering formats the LLM likely outputs: prose with em-dash separators, JSON arrays, preamble+tasks, mixed formats
4. **Add `Preview` field** to `PhaseEvent` struct for diagnostic capture

This closes the visibility gap that lets Architect failures become invisible. It brings Architect to feature parity with Reflector on parser robustness.

---

**This is the single highest-priority infrastructure gap for the hive repository after Reflector stabilization. Architect failures have zero diagnostic context.**