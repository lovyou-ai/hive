# Test Report: PipelineStateMachine fixes

## What Was Tested

Three gaps in `pkg/runner/pipeline_state_test.go` after the Builder's fixes to `PipelineStateMachine`.

### 1. `TestPipelineTransitionFromUnknownState`
- **What:** `Transition()` called from a state not in `pipelineTransitions` (the first `!ok` branch at line 115-117 of `pipeline_state.go`)
- **Verified:** Returns an error; state is unchanged
- **Why it matters:** The existing `TestPipelineTransitionInvalid` only tested an invalid *event* from a valid state. This tests an invalid *state* entirely — a distinct code path that was unexercised.

### 2. `TestInferEventCriticRevise`
- **What:** `inferEvent("critic")` with `critique.md` containing `VERDICT: REVISE`
- **Verified:** Returns `EventCritiqueRevise`
- **Why it matters:** The REVISE→StateBuilding loop is the fix cycle. It was untested — no test confirmed the critic verdict influenced the next state.

### 3. `TestInferEventCriticPass`
- **What:** `inferEvent("critic")` with `critique.md` containing `VERDICT: PASS`
- **Verified:** Returns `EventCritiquePass`
- **Why it matters:** Completes the critic verdict coverage; confirms the default pass path.

## Side Note: `hasFixes` Dead Code

The `hasFixes` variable in `Run()` can never be true without `hasOpen` also being true (the `hasFixes` condition is a strict subset of the `hasOpen` condition). The `hasOpen || hasFixes` guard is therefore equivalent to just `hasOpen`. Not a correctness bug, but dead code worth cleaning up.

## Results

```
=== RUN   TestPipelineTransitionValid              PASS
=== RUN   TestPipelineTransitionInvalid            PASS
=== RUN   TestRunBoardClearStartsAtDirecting       PASS
=== RUN   TestRunExistingTasksStartsAtBuilding     PASS
=== RUN   TestPipelineTransitionFromUnknownState   PASS  [new]
=== RUN   TestInferEventCriticRevise               PASS  [new]
=== RUN   TestInferEventCriticPass                 PASS  [new]

go test ./pkg/runner/ → ok (3.054s, all tests pass)
```

## Coverage

- Valid transitions: all 13 covered ✓
- Invalid event in valid state: covered ✓
- Invalid state (no transitions): covered ✓ [new]
- Board clear → StateDirecting: covered ✓
- Open tasks → StateBuilding: covered ✓
- Critic REVISE verdict: covered ✓ [new]
- Critic PASS verdict: covered ✓ [new]
- `makeRunner` error propagation (main.go): not unit-testable (CLI integration point)

## Status

PASS — all tests clean.

@Critic ready for review.
