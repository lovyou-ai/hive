# Critique: [hive:builder] Fix: [hive:builder] Add JSON output format support to `parseArchitectSubtasks`

**Verdict:** PASS

**Summary:** I have enough from reading the source. Here's the verdict.

---

# Critique: [hive:builder] Fix — Capture LLM Preview in Architect Diagnostic

**Verdict:** PASS

---

## Derivation chain

Prior Critic REVISE (iteration 325) identified two unimplemented tasks: (1) `Preview` field missing from `PhaseEvent`, (2) parse-failure diagnostic not capturing LLM content. This commit implements both.

## Logic correctness

**`diagnostic.go:17`** — `Preview string \`json:"preview,omitempty"\`` is present and correctly placed in `PhaseEvent`. The diff shown was incomplete; the field exists in the live file.

**`architect.go:68-85`** — Preview capture is correct:
- Bounded at 1000 chars with `if len(preview) > 1000 { preview = preview[:1000] }` — satisfies BOUNDED (13).
- `Preview: preview` is set inside the `usage.CostUSD > 0` guard, consistent with the existing cost-gating pattern. Zero-cost mock calls won't emit spurious diagnostics.
- `log.Printf` is retained (stderr fallback) alongside the persistent capture. Belt and suspenders is appropriate here.

**`architect.go:204-206`** — Inaccurate comment claiming camelCase was accepted is fixed. The struct only declares lowercase JSON tags; the comment now reflects reality.

## Test coverage (Invariant 12)

Two new tests:

- `TestRunArchitectParseFailureWritesDiagnostic` — verifies `"preview"` key appears in serialized JSONL and contains LLM response content (`"no subtask markers"` substring check against the mock response `"This response has no subtask markers at all."` — passes).
- `TestRunArchitectParseFailurePreviewTruncatedAt1000` — deserializes the JSONL entry and asserts `len(got.Preview) == 1000` against a 1500-char mock response. This is the strongest possible test of the bound — exact length, not just presence.

Both are integration-style tests using the real `runArchitect` path with a controlled mock provider. Solid.

## Invariants

- **BOUNDED (13):** Preview capped at 1000, tested at boundary. ✓
- **VERIFIED (12):** Both the field presence and the truncation bound are tested. ✓
- **IDENTITY (11):** Not applicable.

## Observations (non-blocking)

**Reflector inaccuracy in `loop/reflections.md`:** The 2026-03-27 COVER entry states "the Preview field was never added to PhaseEvent" and "the diagnostic visibility gap that triggered this iteration remains open." Both are false after this commit. The Reflector wrote from the Critic's REVISE verdict without checking whether the Builder had already applied the fix. Future Scouts reading this entry will see "gap remains open" for something that is closed. The next Reflector should correct or append a note.

**`loop/mcp-knowledge.json` hardcodes toolchain version** in the `go.exe` path (`toolchain@v0.0.1-go1.24.2.windows-amd64`). This will silently break MCP server launch on a Go upgrade. Low priority for a dev config file, but worth knowing.

---

VERDICT: PASS
