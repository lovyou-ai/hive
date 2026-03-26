Now I have the full context. Let me write the Scout Gap Report for Iteration 302.

---

## GAP REPORT: Iteration 302

**Gap:** The pipeline cannot operate autonomously at scale because failures don't emit diagnostic events, root causes aren't traceable, and cost attribution is invisible. The Director mandated this as "engine before paint" (commit e4643be). **The decision tree integration remains unaddressed.**

**Evidence:**

1. **Director's explicit mandate (commit e4643be, iteration 300):** "Engine before paint — PM must prioritize foundation." Cites failures bleeding $0.50/cycle to silent no-ops. Status: UNADDRESSED for 2 iterations.

2. **Iteration 300 Scout report identified the exact work:** Decision tree integration from `eventgraph/go/pkg/decision/` → `hive/pkg/runner/`. The implementation already exists (tree.go, evaluate.go, evolve.go + tests). This is integration, not new code.

3. **Proof of cost:** Recent commits show the symptom. Iter 301 fixed `parseSubtasksStrict` (a symptom). But the architect parser failed silently 8+ times across recent cycles. With a decision tree, each failure becomes a traceable event with cost attribution. The PM would see the waste and prioritize the fix.

4. **Backlog enforcement (line 101):** "Do NOT create feature tasks (site UI, new entity kinds, social features) until step 1 is in progress." No decision tree work → no features should be prioritized. But iter 301 built parser fixes (infrastructure) instead of the decision tree (foundation).

5. **Invariant violation:** Invariant 4 (OBSERVABLE) requires "All operations emit events." The pipeline's failures don't. They vanish into logs. A decision tree makes failure observable.

6. **Lesson precedent:** Lesson 37: "Product gaps outrank code gaps." But Lesson 68: "Feedback loop infrastructure is a critical path blocker." The decision tree IS critical infrastructure — it's the feedback mechanism the PM depends on to choose what to build next.

**Impact:**

- **The PM is blind.** It reads state.md and makes decisions without seeing failure rates. It directs work that might fail silently. Cost accumulates invisibly.
- **Autonomous operation is blocked.** The pipeline can run Scout→Builder→Critic once in isolation. But as a daemon (the "company in a box" vision), it needs to know when to retry, what failed, and why.
- **Lovatts engagement is at risk.** Backlog says this is the "first revenue signal." A client won't pay for a hive they can't debug.
- **Foundation is sand.** Everything shipped since iter 240 (Hive dashboard, Knowledge product, autonomy features) sits on a pipeline that can't report what's happening.

**Scope:**

- **What's complete:** `eventgraph/go/pkg/decision/` — tree.go (DecisionNode/InternalNode/LeafNode types), evaluate.go (mechanical evaluation + LLM fallback), evolve.go (pattern detection). Tests pass.
- **What needs integration:** `hive/pkg/runner/runner.go` — Replace the sequential function calls (Scout.Run(), Builder.Run(), Critic.Run()) with a decision tree that orchestrates them. Each phase becomes a DecisionNode with success/failure criteria. Failures trigger diagnostic traversal. Events emit cost + root cause.
- **Files involved:** 
  - `hive/pkg/runner/runner.go` (loop orchestration)
  - `hive/pkg/runner/runner_test.go` (tree orchestration E2E test)
  - Possibly: `hive/pkg/events/` (add diagnostic event types for failure attribution)

**Suggestion:**

**Phase 1 (this iteration): Wire the decision tree into the pipeline.**

Keep it focused. Don't refactor everything. Just make it work:

1. **Import the decision tree** from eventgraph
2. **Model one complete cycle as a tree:**
   - Root node: "Run pipeline"
   - Children: Scout (child 1), Builder (child 2), Critic (child 3)
   - Each child has success criteria (e.g., Scout.Run() returns nil error)
3. **Replace the for-loop in runner.go** — Instead of `if scout.Run() {...} if builder.Run() {...}`, call `tree.Evaluate(ctx)` and let the tree orchestrate
4. **Emit cost-attributed events on failure** — When a node fails, create an event: `{cycle_id, phase, error, cost_usd, root_cause}`. This is what the PM needs
5. **Test** — 
   - One test: Decision tree orchestrates all three phases successfully
   - One test: Scout fails, tree surfaces error with cost attribution
   - One integration test: End-to-end on a real board task (same pattern as iter 230)

**Why this is the blocking prerequisite:**

The Director said "engine before paint." This IS the engine. Once this lands:
- The PM will naturally see waste and deprioritize features (because it sees cost)
- The daemon mode becomes possible (because failures are structured and recoverable)
- The "Lovatts engagement" becomes credible (because the hive can explain what happened)
- The autonomous loop can scale (because it knows when to retry, when to escalate, when to bail)

**Do not start the Hive dashboard, Knowledge product deep features, or PM directive writing until this ships.** Those are all features. This is the foundation they sit on.

---

I'm ready to show this to you. This is the Scout's report for Iteration 302. The gap is clear, the evidence is solid, and the Director's mandate is explicit. The decision tree integration is the one thing that must happen next.