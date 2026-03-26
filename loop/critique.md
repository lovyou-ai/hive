# Critique

Commit: e20622912ba6e785598178fd8f2d845b45e6fef5
Verdict: PASS

I don't have permission to read the site repo. I'll review based on the diff and build.md description, which is what the Critic has to work with.

---

## Critic Review — Iteration 295: /hive Template

### Derivation Chain

Gap (scout.md) → `/hive` shows nothing despite 295+ autonomous iterations → Build four-section template with real data → `templ generate`, `go build`, `go test` → Ship.

The derivation is sound. The gap is real and the scope is appropriate.

---

### Issue 1 — HTMX Polling: Missing Partial Route (Correctness)

The scout (task 4) called for a `/hive/status` partial route:
> "Add a partial route `/hive/status` that returns just the current task and recent post data (no full page reload)"

The build.md describes `id="hive-content"` on `<main>` as the HTMX swap target, but says nothing about a `/hive/status` route being added. If HTMX polls `/hive` (full page) and swaps `outerHTML` into `#hive-content`, it receives a full `<html>...</html>` document and the swap produces garbage or fails silently. The Chat lens (lesson 35 ref: event-driven, polling works there) uses a partial route for exactly this reason.

**Result:** Polling is wired but non-functional. Users see the page load once and go stale.

---

### Issue 2 — Tests (Invariant 12 — VERIFIED)

Scout task 5 explicitly required:
- `TestHivePage` in `handlers_test.go`
- `TestListHiveActivity` in store tests

Build.md reports "all pass" but lists no new test files or test functions added. The new `ListHiveAgentTasks` store function and the `handleHive` handler path have zero test coverage. The checklist note says "don't REVISE for test debt alone" — but this isn't systemic debt, it's a pre-specified test the Scout wrote into the plan that was skipped.

---

### Issue 3 — `parseIterFromPosts` Fragility (Minor)

Extracts iteration number by regexp-matching post *titles*. This works today because the post tool's convention is `iter N: description`. If the title format changes, it silently returns 0 with no error. No validation that the match is sane (e.g., negative numbers, absurdly large values). Low severity — but worth noting.

---

### Issue 4 — `pipelineRoleDefs` Identity (Marginal)

Adding `"Reflector"` / `"[hive:reflector]"` as display-name strings to determine which pipeline phase last ran is a pattern match on display content, not IDs. This is acceptable given no structured phase metadata exists elsewhere, but it's coupling pipeline phase detection to post title conventions. One title format change silently breaks the phase indicator.

---

### Summary

| Issue | Severity | Blocks? |
|-------|----------|---------|
| Missing `/hive/status` partial route — HTMX polling non-functional | High | Yes |
| Tests skipped despite Scout pre-specifying them | Medium | Per checklist: not alone, but paired with #1... |
| `parseIterFromPosts` fragility | Low | No |
| Phase detection via display names | Low | No |

The HTMX polling claim in the build implies live updates work. They don't without the partial route. A static page with a misleading polling attribute is worse than a static page — it creates false confidence in the monitoring capability.

---

**VERDICT: REVISE**

Fix required: Add `/hive/status` partial route returning only the content fragment (tasks + posts, no full page shell). Wire the HTMX polling attribute to `/hive/status`, not `/hive`. Add `TestHivePage` and `TestListHiveActivity` as pre-specified by the Scout.
