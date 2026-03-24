# Critique — Iteration 226: Critic Role

**Verdict: PASS** (with notes)

---

## Derivation Check

### Gap → Scout: ✓ VALID
Scout correctly identified the gap: no automated code review. Builder ships code but nobody checks it. Critic role is Phase 2 item 9.

### Scout → Build: ✓ VALID
Critic implemented as specified: git log scan, diff review, Reason() call, verdict parsing, fix task creation. 170 lines. Tests cover all helpers.

### Build → Verify: ✓ VALID
- Build passes, 23 tests pass
- E2E tested: Critic found and reviewed the correct builder commit
- Bug found and fixed (regex escaping)

---

## Invariant Audit

| Invariant | Status | Reason |
|-----------|--------|--------|
| 12 VERIFIED | ✓ Pass | 9 tests for critic helpers. E2E tested on production. |
| 13 BOUNDED | ✓ Pass | 24h window, diff truncation at 15KB. |

---

## Issues Found

### 1. Critic PASSed the allowlist miss (medium)
The Critic reviewed the Policy commit (which was missing `KindPolicy` in the intend allowlist) and returned PASS. The review prompt asks about "ALL relevant guards, allowlists, and switch statements" but the diff only shows what changed — the allowlist line is 400 lines away and not in the diff.

**Root cause:** Diff-only review can't catch omission errors in distant code. The Critic sees what was added but doesn't know what SHOULD have been added elsewhere.

**Fix options:**
- (a) Include the full file context (expensive, blows up prompt)
- (b) Have the Critic use Operate() to grep for related patterns (slower but accurate)
- (c) Add entity-specific checklists ("new KindX → check intend allowlist")

Option (b) is the right long-term answer. Defer to Phase 3.

### 2. No review state persistence (low)
The Critic has no memory of which commits it has already reviewed. On restart, it re-reviews all commits from the last 24h. Not a problem at current volume but will waste money when the builder ships multiple commits per day.

**Fix:** Track reviewed commit hashes in a file or via API comments.

---

## Verdict: PASS

The Critic infrastructure works. It finds builder commits, reviews them via Reason(), and can create fix tasks. The diff-only limitation (issue 1) is a known design tradeoff — the Critic catches syntax/pattern errors but not omission errors in distant code. This is exactly the class of bug the human caught in iter 225. Upgrading to Operate()-based review (Phase 3) will close this gap.
