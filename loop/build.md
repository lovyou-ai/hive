# Build: Fix: Knowledge API omits causes field on claim nodes — Invariant 2

## What Was Built

Fixed `syncClaims` in `cmd/post/main.go` to include the `causes` field when decoding claim nodes from the `/knowledge` API response and writing them to `loop/claims.md`. Added `TestSyncClaimsWritesCauses` to pin the fix.

## Investigation

The server-side code was already correct:
- `handleOp`'s `assert` case reads `causes` from the request and stores them via `CreateNode`
- `ListNodes` populates `Causes` via `pq.Array(&n.Causes)`
- `Node.Causes []string \`json:"causes,omitempty"\`` is included in the JSON response
- `TestAssertOpReturnsCauses` in `site/graph/knowledge_test.go` already covered end-to-end

The `cmd/post` side was also already correct:
- `assertCritique`, `assertLatestReflection`, `assertScoutGap` all pass `causeIDs` via `p["causes"] = strings.Join(causeIDs, ",")`
- Tests `TestAssertCritiqueSendsCauses`, `TestAssertLatestReflectionSendsCauses`, `TestAssertScoutGapSendsCauses` already verified these

The actual gap: `syncClaims` decoded claims into an anonymous struct that omitted `Causes`, so `loop/claims.md` never included provenance links. MCP `knowledge_search` operates on `claims.md` — without causes, agents cannot trace claim provenance.

## Files Changed

- `cmd/post/main.go` — Added `Causes []string \`json:"causes"\`` to the `syncClaims` decode struct; added `**Causes:** id1, id2` line in the claims.md writer
- `cmd/post/main_test.go` — Added `TestSyncClaimsWritesCauses` to pin that causes appear in claims.md output

## Verification

```
go.exe build -buildvcs=false ./...   # clean
go.exe test -buildvcs=false ./...    # all pass
```
