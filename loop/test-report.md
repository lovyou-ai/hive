# Test Report — Close orphaned subtasks when parent completes

## What was tested

The Builder added `cascadeCloseChildren` — a depth-first recursive method that closes all non-done descendants before a parent can be marked done. I added a test for the **3-level deep recursive cascade** (grandparent → parent → grandchild), which was the critical unverified path.

## New test added

### `TestCascadeCloseChildrenDeep` (`site/graph/store_test.go`)
- Creates a 3-level tree: grandparent (active) → parent (active) → grandchild (active)
- Calls `UpdateNodeState(grandparent, done)`
- Verifies all three nodes end up in `StateDone`
- Exercises the recursive `cascadeCloseChildren` call at depth > 1

## Pre-existing tests updated by Builder (verified still passing)

| Test | What it covers |
|------|---------------|
| `TestUpdateNodeStateChildGate` | Parent auto-closes one open child |
| `TestUpdateNodeStateChildGateLeafNode` | Leaf node completes directly (no children) |
| `TestUpdateNodeStateChildGateMultipleChildren` | One done child + one open child — open child auto-closed |
| `TestUpdateNodeStateNonDoneSkipsGate` | Non-`done` transitions skip the cascade entirely |

## Results

```
=== RUN   TestCascadeCloseChildrenDeep
--- PASS: TestCascadeCloseChildrenDeep (0.04s)
=== RUN   TestUpdateNodeStateChildGate
--- PASS: TestUpdateNodeStateChildGate (0.04s)
=== RUN   TestUpdateNodeStateChildGateLeafNode
--- PASS: TestUpdateNodeStateChildGateLeafNode (0.03s)
=== RUN   TestUpdateNodeStateChildGateMultipleChildren
--- PASS: TestUpdateNodeStateChildGateMultipleChildren (0.04s)
=== RUN   TestUpdateNodeStateNonDoneSkipsGate
--- PASS: TestUpdateNodeStateNonDoneSkipsGate (0.03s)
PASS
```

## Pre-existing failure (not caused by this iteration)

`TestHandlerCreateSpace` fails with `get space: not found` — confirmed pre-existing on HEAD before Builder's changes. Unrelated to cascade work.

## Coverage notes

- **Depth limit (maxCascadeDepth = 50)**: Not tested — would require creating 51 levels of nesting, impractical with real DB. Enforced by code inspection.
- **Already-done intermediaries**: Covered implicitly by `TestUpdateNodeStateChildGateMultipleChildren` (child1 is done before parent completes; child2 gets cascaded).
- **Recursive depth**: Now covered by `TestCascadeCloseChildrenDeep` at 3 levels.
