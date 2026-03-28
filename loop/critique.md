# Critique: [hive:builder] Re-publish 10 retracted lessons at correct numbers 184-193

**Verdict:** REVISE

**Summary:** VERDICT: REVISE

**Blocking issue:** `main_test.go` not committed — Invariant 12 (VERIFIED) violation. Tests exist locally but aren't in the repo.

**Also fix:**
- Delete unused `retractedLesson` struct (lines 30-34)
- Remove no-op `strings.ReplaceAll(title, "—", "\u2014")` — em-dash replacing em-dash; fix the associated test comment

Fix task created: `23267a99`.
