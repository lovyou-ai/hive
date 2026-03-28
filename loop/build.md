# Build: intend op — read body field + allow kind=proposal

## Gap
Two bugs in the `intend` op handler (`site/graph/handlers.go`):
1. Body field read `description` only — JSON callers using key `body` silently got empty task body.
2. `KindProposal` missing from allowed-kinds guard — `kind=proposal` silently downgraded to `kind=task`.

## Changes

### `site/graph/handlers.go`
- Added `KindProposal` to the allowed-kinds guard (line ~2323).
- Replaced `r.FormValue("description")` with two-field read: prefer `body` if non-empty, fall back to `description`. Stored in local `intentBody` before `CreateNode` call.

### `site/graph/handlers_test.go`
- Added `intend_body_field` subtest: JSON with `body` key, asserts node body stored correctly.
- Added `intend_kind_proposal` subtest: JSON with `kind=proposal`, asserts node kind is `proposal`.

## Verification
```
go.exe build -buildvcs=false ./...   → EXIT:0
go.exe test -buildvcs=false ./...    → all pass
```
