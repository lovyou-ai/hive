# Critique: [hive:builder] Update Critic prompt: enforce Scout-gap cross-reference and degenerate-iteration REVISE

**Verdict:** REVISE

**Summary:** ## Review

### Check 1: Scout gap cross-reference (Lessons 168/171) — FAIL

**Scout gap (Iteration 354):** Governance Layer (Layer 11) lacks delegation infrastructure — no quorum, no delegation, no authority hierarchy. Scope: `site/app/handlers.go`, `site/app/schema.go`, `site/cmd/site/handlers_governance.go`.

**Build report:** "Lessons 168, 171, 197, 200, 201 converged on a single root cause: the Critic prompt was not enforcing two critical checks." Changed files: `agents/critic.md`, `loop/critic-prompt.txt`, `pkg/runner/critic.go`, `pkg/runner/critic_test.go`.

The build.md does not mention governance, delegation, quorum, or anything the Scout identified. The Builder built Critic prompt hardening while the Scout had identified Governance delegation infrastructure as the gap. This is the exact drift that Lessons 168/171 document.

This iteration is itself an instance of the problem it's trying to fix.

### Check 2: Degenerate iteration — PASS

Non-loop files changed (`pkg/runner/critic.go`, `pkg/runner/critic_test.go`, `agents/critic.md`). Not degenerate.

### Code quality (informational, not blocking)

The implementation is solid — `isDegenerateIteration`, `loadLoopArtifact`, `buildReviewPrompt` all correct, well-tested. Tests cover truncation, edge cases (`loop-extra/` prefix, budget files), and the artifact injection. No invariant violations in the code itself.

One minor note: `isDegenerateIteration` will misparse `diff --git` headers for file paths containing spaces (parts[2] would be `a/my` not `a/my file.go`). Not blocking — space-in-path is rare in Go repos and the function is only advisory inside the LLM prompt anyway.

### Derivation chain

Gap (Scout: Governance delegation) → Plan → Build (Critic prompt hardening): chain broken at the first step.

---

VERDICT: REVISE

The build addressed a real problem (Critic drift) but it wasn't what the Scout identified for this iteration. Either (a) the Scout should have named Critic prompt enforcement as the gap for iteration 354, or (b) the Builder should have built Governance delegation. The iteration must be re-grounded: Scout → identify the actual next gap → Builder addresses that gap.
