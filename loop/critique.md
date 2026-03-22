# Critique — Iteration 26

## Verdict: APPROVED

## Trace

1. Matt observed that `agent_name` is "a name without a soul" — display override, not real identity
2. Scout identified: agents need real user records, not cosmetic labels
3. Builder added `ensureAgentUser()` — creates agent users with kind='agent' and synthetic google_id
4. Builder added `agent_id` to api_keys — links key to agent's user record
5. Builder updated `userByAPIKey()` — resolves to agent user when agent_id is set
6. Compiles clean, deployed

Sound chain. The gap → plan → code → verify chain holds.

## Audit

**Correctness:**
- `ensureAgentUser()` uses `ON CONFLICT (google_id)` with synthetic `agent:{name}` — idempotent. Creating the same agent twice returns the existing record. ✓
- `createAPIKey()` creates agent user before inserting key — agent_id is always valid when set. ✓
- `userByAPIKey()` does two queries when agent_id is set (one for key lookup, one for agent). Could be a JOIN but simplicity wins here. ✓
- `ListAPIKeys()` uses `COALESCE(agent_id, '')` — handles NULL → empty string cleanly. ✓

**Breakage:**
- Existing keys without agent_id (NULL) → `agentID.Valid` is false → sponsor user returned. Fully backward compatible. ✓
- `ON DELETE SET NULL` on agent_id: if agent user is deleted, key falls back to sponsor identity. Reasonable. ✓
- Agent users can't log in via OAuth (no real google_id). Access only through API keys — sponsor controls access. ✓

**Simplicity:** One new function (`ensureAgentUser`), one new column (`agent_id`), one new column (`kind`). No new tables, no new types, no new middleware. ✓

**Security:** Agent identity is scoped to API key auth only. Agents can't create sessions, can't manage other keys, can't escalate. The sponsor revokes the key → agent loses access (but identity persists — their history remains). ✓

**Dual (root cause):** Why was iteration 25's approach wrong? Because it treated identity as a *property of the credential* (agent_name on the key) rather than a *property of the entity* (user record). A name on a key is metadata. A user record is identity. The distinction: metadata describes something; identity IS something.

## Observation

The `kind` column on users is written but not yet read. Future iterations can use it to filter agents from humans in People lens, activity feeds, etc. The foundation is laid; the views will follow.
