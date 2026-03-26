Based on my review of the state, scout history, and backlog, I can identify the next critical gap. Let me write the gap report for iteration 299.

---

## SCOUT REPORT — Iteration 299

**Gap:** The hive cannot manage itself autonomously because its own operational artifacts live in markdown files instead of on the graph. Tasks are in `state.md`, specs are in `backlog.md`, decisions are in `reflections.md` — invisible to the system it builds. The hive talks about being a substrate for collective existence but doesn't use its own substrate.

**Evidence:**

1. **Backlog item explicitly calls this "URGENT" — Dogfooding:**
   - Lines 208-219 of `loop/backlog.md` state: "The civilization doesn't live in its own product. Tasks are in state.md, not on the board."
   - Current state: Hive tasks are tracked in `loop/state.md` (90-line document), not as graph nodes. Scout can't query them. Builder doesn't pick them from the board.
   - Proof: `pkg/runner/scout.go` reads `state.md` manually; `pkg/runner/builder.go` claims tasks from `loveyou.ai/board` but hive's own tasks aren't there.

2. **Specs are invisible to the system:**
   - Hive has 8 converged specs (work-general-spec.md, social-spec.md, unified-spec.md, etc.) but they live as markdown files in `loop/`.
   - Knowledge layer (iter 6) was implemented on the graph but the HIVE doesn't use it.
   - `loop/backlog.md` line 198-206: "Knowledge layer as the hive's brain (not markdown files)" — the insight is that specs should be Knowledge claims with assert/challenge/verify lifecycle, but they're not.
   - Proof: grep `loop/` for Knowledge assertions about hive specs returns zero.

3. **Lessons accumulate but aren't persisted:**
   - `loop/state.md` has 75 lessons, each discovered through iteration. These should be Knowledge claims (with evidence, verification status, rejection reason).
   - Currently: Lessons are documentation. They can't be queried, linked, or verified. They repeat because the system has no memory of them.
   - Proof: iter 298 scout.md line 32 states "Lessons 71, 73, 74 all document this pattern but it repeats unchanged."

4. **Escalations are advisory, not binding:**
   - When Scout flags critical work (iter 298 line 34-39: daemon mode, PR workflow, artifact writes), these become tasks in state.md, not board nodes.
   - No way to query: "what escalations are open?" or "how many cycles has this item been deferred?"
   - Related lesson 68 (iter 297 reflections): "Feedback loop infrastructure is a critical path blocker."
   - Proof: iter 298 scout lists 5 open escalations with no tracking mechanism to verify closure.

**Impact:**

- **The hive cannot be autonomous** — it can't query its own state via the API. It needs Matt to read state.md and say "next."
- **The hive cannot improve itself** — lessons aren't verifiable (is lesson 71 still true?), escalations aren't trackable (are these still blockers?).
- **The hive is not an example of the product** — it preaches that all activity should be on the graph, but its own activity isn't. This breaks the soul: "Take care of yourself" — the hive isn't taking care of itself.
- **Bus factor remains unfixable** — daemon mode (iter 298 backlog line 225-233) requires the hive to see its own board state. Without dogfooding, daemon mode can't work.
- **Lovatts engagement can't launch** — company-in-a-box (backlog line 32-45) requires the hive to be a working example. A civilization that doesn't use its own infrastructure isn't compelling.

**Scope:**

Hive repo:
- `loop/state.md` → migrate hive tasks to board (as priority-high task nodes)
- `loop/backlog.md` → migrate open items to board as projects/goals
- `loop/*.md` specs → migrate to Knowledge claims with assert ops
- `pkg/runner/scout.go` → read board instead of state.md for hive tasks
- `pkg/runner/runner.go` → hive tasks become like any other task (assignable, statusable, closeable)

Site repo:
- Hive task space configuration (what lenses are visible, which ops are allowed)
- Hive agent persona + board integration (the agent sees its own tasks)

**Suggestion:**

Build a **Dogfooding Phase 1: Task Migration** that moves the hive's operational work onto the board:

**Part 1 (Immediate):**
- Create a space called "Hive Operations" (or let the hive create its own, using the space creation flow)
- Migrate `loop/state.md`'s "What to build next" section → board as tasks (high priority, assigned to Builder agent)
- Migrate iter 298 escalations (daemon mode, PR workflow, artifact validation, Lovatts prep) → board as tasks with "escalation" tag
- Scout reads `loveyou.ai/spaces/{hive-space-id}/board` instead of state.md

**Part 2 (Foundation for dogfooding depth):**
- Converged specs (work-general-spec.md, social-spec.md) → Knowledge layer with assert ops, author_id = Architect, status = verified
- lessons → Knowledge claims (assert op, kind="lesson", importance=1-5 based on line number in state.md)
- Council sessions → Knowledge claims (from council.md files, linked to findings)

**Verification criteria:**
- Scout can query hive tasks via the board API
- Builder can claim a hive task from the board (not hardcoded in runner logic)
- Critic can see hive tasks in context of all tasks (no special cases for hive work)
- All five escalations from iter 298 are visible on the board with clear status

**Done criteria:** The hive's next iteration (300+) runs entirely against board state, not state.md. Scout reads the board, creates tasks on the board, Builder claims from the board. Zero hive tasks in state.md after completion.

**Priority:** URGENT. This is the prerequisite for daemon mode (iter 298 backlog), the prerequisite for Lovatts engagement, and the proof that the product works. Without it, the hive remains a toolchain pretending to be a civilization.

---

**Ready for Architect phase.**