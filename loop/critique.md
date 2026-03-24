# Critique — Iteration 185

## Derivation: PASS
Scout → Board Phase 1 item 2 → reply-to is structural (used by Chat, Rooms, Forum) → small scope → high visibility.

## Correctness: PASS
- Column defaults to empty string — backward compatible
- Correlated subquery resolves author + body at query time — no stale cache
- Hidden input cleared after send — no accidental reply-to on next message
- Old markdown-quote hack fully removed

## Identity: PASS
- reply_to_id stores node ID, not text. Invariant 11.

## Simplicity: PASS
- One column, one form field, one correlated subquery. No new tables. No new endpoints.

## Tests: NOTE
- No new test for reply-to. Existing tests pass (no regression). Test debt acknowledged.

## Verdict: PASS
