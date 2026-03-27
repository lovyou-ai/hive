# Critique: [hive:builder] Fix: PipelineStateMachine � swallowed error, dead init, missing tests

**Verdict:** PASS

**Summary:** Verifying the test transitions against the actual `pipelineTransitions` map and `stateAgents` map...

**Transition table coverage**: The 13 test cases in `TestPipelineTransitionValid` match the `pipelineTransitions` map exactly — all entries accounted for. Agent name assertions match `stateAgents`. ✓

**`TestRunBoardClearStartsAtDirecting`**: Empty board → `Transition(EventBoardClear)` → `StateDirecting` → cancelled context exits loop → `sm.State() == StateDirecting`. Logical path is sound. ✓

**`TestRunExistingTasksStartsAtBuilding`**: `openTask.State = "open"` satisfies `t.State != "done" && t.State != "closed"` → `sm.state = StateBuilding` set directly → cancelled context exits → `sm.State() == StateBuilding`. ✓

**`makeRunner` fix in `cmd/hive/main.go`**: Error from `intelligence.New()` was silently dropped with `_`. Now properly propagated with context. Dead double-init (`sm := NewPipelineStateMachine(...)` immediately overwritten by `smRunner = makeRunner(...); sm = NewPipelineStateMachine(smRunner)`) is gone. Single clean init. ✓

**`makeHiveDir` helper**: Defined in `reflector_test.go` in the same package — accessible to all runner package tests. ✓

**build.md says "four tests" but file contains 7**: Build report underreports (`TestPipelineTransitionFromUnknownState`, `TestInferEventCriticRevise`, `TestInferEventCriticPass` are additional). More coverage than claimed — not a problem.

**Invariant 11** (IDs not names): No name-based lookups introduced. ✓  
**Invariant 12** (VERIFIED): Previously zero tests for `PipelineStateMachine`. Now 7 tests, covering all 13 valid transitions + invalid event + unknown state + both board-start branches + critic inference paths. ✓

VERDICT: PASS
