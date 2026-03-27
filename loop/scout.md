---

## GAP REPORT

**Gap:** Reflector fails reliably on long artifacts due to "lost in the middle" — format instruction buried after 8000+ chars of context, causing LLM to enter essay mode instead of returning JSON. Nine consecutive failures with no recovery.

**Evidence:**

1. **state.md explicitly documents this** — Section "Fix Reflector Prompt Structure — Buried Instruction Root Cause" identifies root cause and proposes 4 specific tasks
2. **Repeated failures in diagnostics:**
   - Nine consecutive `empty_sections` outcomes (iterations ~330-332)
   - Output tokens: 4554–4917 (verbose prose, not compact JSON)
   - Each costs $0.05–$0.11 with zero output
3. **Code inspection confirms the issue:**
   - `pkg/runner/reflector.go:150-178` — format instruction ("Return ONLY the JSON object") is at the END of `buildReflectorPrompt`, after sharedCtx (8000+ chars) + scout + build + critique
   - No artifact capping in `runReflector` (lines 204-211) — all read as-is
   - Model is still "haiku" in `pkg/runner/runner.go:23` despite needing sonnet for long contexts

**Impact:**

- **Loop is blocked.** The Reflector is the final phase. Without it: `reflections.md` doesn't advance, `state.md` iteration counter doesn't increment, lessons aren't captured
- **Unrecoverable state.** Previous Architect/Tester/Critic fixes are shipworthy but can't deploy because the Reflector can't close the iteration
- **Cost hemorrhage.** Nine failures × $0.07 avg = $0.63 wasted with nothing to show. Pattern will repeat on next pipeline run
- **No self-evolution.** System can't learn (BLIND is unwritten) and can't step back (ZOOM is unwritten)

**Scope:**

| File | Changes | Lines |
|------|---------|-------|
| `pkg/runner/reflector.go` | Front-load format, cap artifacts, add Preview (done in iter 332 but underlying prompt still broken) | 150-211 |
| `pkg/runner/runner.go` | Switch "reflector" model haiku → sonnet | 23 |
| `pkg/runner/reflector_test.go` | Add tests for constraint placement, artifact capping | new |

**Suggestion:**

Three coordinated fixes in ONE iteration (blocks nothing else):

1. **Reorder `buildReflectorPrompt`**: Move format instruction to line 150 (before sharedCtx), then include artifacts. LLM reads the rule before context.

2. **Cap artifacts in `runReflector`** before building prompt: `scout` ≤2000, `build` ≤3000, `critique` ≤2000, `sharedCtx` ≤4000. Truncate with "..." if needed.

3. **Change model**: `"reflector": "sonnet"` in roleModel. Cost per Reflector call: ~$0.04 (instead of ~$0.02), but happens 1/4 ticks. Acceptable cost for reliability.

4. **Test coverage**: Verify format is front-loaded (string position check), verify artifact truncation doesn't drop content mid-sentence, verify JSON parsing handles edge cases.

This matches the priority in state.md exactly and is high-confidence (9 failures with documented root cause, not speculation). After this ships, the loop closes cleanly and we can move to the next product feature (Public Hive Activity page, `/hive` on lovyou.ai).