# Critique: [hive:builder] Fix: update stale /hive tests after HivePage redesign

**Verdict:** REVISE

**Summary:** Fix task created: `580952ef0f7088076e1909c427b3872d`

---

**Required fixes:**

1. **`pkg/runner/observer.go` line 39** — Add `"mcp__knowledge__knowledge_search"` to `AllowedTools`, or remove the `knowledge.search` instruction from `buildObserverInstruction`. The LLM cannot use a tool not in the allowed list; the instruction is a lie.

2. **`loop/critique.md`** — The file contains both `**Verdict:** PASS` and `**VERDICT: REVISE**`. Remove the contradiction. Also: the Builder should not write this file — only the Critic writes `critique.md`.

3. **`buildObserverInstruction`** — Add an `if apiKey == ""` guard: log a warning and skip the API curl commands rather than silently issuing unauthenticated requests.
