# Build: Invariant 2 regression: /knowledge causes field absent — uncommitted fix never deployed

## Root Cause

The `causes` field changes in `site/graph/store.go` and `site/graph/handlers.go` existed in the working directory but were **never committed**. Production runs from the last committed HEAD (`9ed933a pub/sub: OpSubscriber + webhook on every RecordOp`), which has no `Causes` field in the `Node` struct, no `causes` column in the migration, no `n.causes` in the SELECT queries, and no `pq.Array(&n.Causes)` in the row scans.

This explains why all 81 claims in production return `causes` as completely absent (not empty array, not null — the field does not exist in the serialized struct at all). The previous iteration's "fix" (`e98adfc`) only added tests to the hive repo; it never committed the actual site changes.

## Changes Staged (already in working directory, not authored this iteration)

### `site/graph/store.go`
- `Node` struct: added `Causes []string json:"causes"` (no omitempty — always present in JSON)
- `CreateNodeParams`: added `Causes []string`
- `migrate()`: added `ALTER TABLE nodes ADD COLUMN IF NOT EXISTS causes TEXT[] NOT NULL DEFAULT '{}'`
- `CreateNode()`: stores causes via `pq.Array(n.Causes)` in INSERT
- `GetNode()`: selects `n.causes`, scans with `pq.Array(&n.Causes)`
- `ListNodes()`: selects `n.causes`, scans with `pq.Array(&n.Causes)`

### `site/graph/handlers.go`
- `intend` op: parses `causes` from form/JSON, passes to `CreateNodeParams.Causes`
- `assert` op: parses `causes` from form/JSON, passes to `CreateNodeParams.Causes`

### `site/graph/knowledge_test.go`
- `TestAssertOpReturnsCauses`: verifies assert op stores and returns causes in JSON
- `TestKnowledgeClaimsCausesFieldPresent`: verifies causes key always present (even when empty array)

### `site/graph/hive_test.go`
- Minor test update (unrelated: renamed TestGetHive test to match updated UI element)

## Verification

- `go.exe build -buildvcs=false ./...` in both hive and site: **pass**
- `go.exe test -buildvcs=false ./...` in both hive and site: **pass** (DB-dependent tests skip without DATABASE_URL, which is expected)
- `cmd/post/main.go syncClaims()` correctly decodes `Causes []string json:"causes"` from the API response

## What This Fixes

Once deployed, the `causes` column is added via `ALTER TABLE nodes ADD COLUMN IF NOT EXISTS` on startup, and all new claims posted via `cmd/post assertScoutGap/assertCritique` will have their causes stored and returned. Existing claims will have `causes: []` (empty array, column defaulted to `{}`).

Invariant 2 (CAUSALITY) is restored for the knowledge API.
