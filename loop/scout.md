I'm ready with the Scout report. Let me present it:

---

# Scout Report — Iteration 265

## Gap

**Test isolation failures block the autonomous loop** — three handler tests fail with duplicate slug constraint errors, violating **Invariant 12** (VERIFIED: "No code ships without tests").

## Evidence

- **build.md (iter 264, lines 40-42):** Three invite handler tests documented as failing:
  - `TestHandlerJoinViaInvite`
  - `TestHandlerCreateInviteHTMX`
  - `TestHandlerRevokeInvite`
  - All fail with duplicate slug constraint errors

- **scout.md (iter 264):** Explicitly escalated as **IMMEDIATE** blocking prerequisite — Invariant 12 violation

- **state.md (iter 265, lines 564-612):** Current directive lists **Task 1: Fix invite handler test isolation (REVISE condition, Lesson 47)** as start of iteration cycle

- **Lesson 47:** "REVISE before new work" — blocking invariant violations resolved at iteration start, not deferred

- **Lesson 64 (reflections.md):** Bottleneck synthesis requires binding response contracts. Gap persisted across iterations 262–265 despite Scout escalation

## Impact

- The test suite does not pass cleanly; `go test ./graph/...` fails on these three handlers
- Invariant 12 (VERIFIED) is violated — code ships but tests don't pass  
- Critic phase cannot verify changes with confidence — test foundation is broken
- Directive explicitly lists this as blocker to Role Membership (Task 2) and agent memory visibility (next priority)
- Each deferred iteration erodes trust in the loop's self-governance

## Scope

**Minimal and production-safe** — test setup only

| File | What | Risk |
|------|------|------|
| `site/graph/handlers_test.go` | Three test functions — generate unique slug suffixes per test run | **Low** — test-only changes |

**Root cause:** Tests reuse hardcoded slugs across functions sharing live Postgres. The `spaces` table unique constraint on `slug` causes collisions.

**Fix pattern:** Generate unique slug per test using `fmt.Sprintf("test-space-%d", time.Now().UnixNano())` or transaction rollback. Pattern already used elsewhere.

## Suggestion

Ship minimal fix:
1. Update three failing test functions to generate unique slugs per test
2. Verify: `go test -run "TestHandlerJoinViaInvite|..." ./graph/` → pass
3. Full suite: `go test ./graph/...` → pass  
4. Deploy: `./ship.sh "iter 265: fix invite handler test isolation"`

---

**Why now:** This is the second iteration Scout has escalated this as IMMEDIATE. Invariant 12 is the foundation of the autonomous loop. Until this passes, Critic cannot verify changes, and the loop's output is unreliable. This is a prerequisite, not a feature.