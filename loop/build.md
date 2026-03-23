# Build Report — Iteration 90

Layer 9 (Relationship) entry: user endorsements.

**New table:** `endorsements` (from_id, to_id, created_at) with composite PK and index on to_id.

**New store queries:** Endorse, Unendorse, CountEndorsements, HasEndorsed, ListEndorsers.

**New route:** `POST /user/{name}/endorse` — toggle endorse/unendorse. Requires auth. Self-endorsement prevented.

**Profile update:** Shows endorsement count in stats, endorser names with links, and endorse/unendorse button for logged-in viewers. Profile page now uses readWrap for optional auth context.

**Scope change:** Moved `readWrap`/`writeWrap` declaration to outer scope so routes outside the DB block can use auth middleware.

8 tables total. 8 product layers touched (1-Work, 2-Market, 3-Social, 4-Justice, 7-Alignment, 8-Identity, 9-Relationship, 10-Belonging).

Deployed. All tests pass.
