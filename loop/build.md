# Build Report — Iteration 101

Added "Chat with Mind" form to the dashboard. Uses the user's first space and the first agent name. One input field + submit button creates a converse op with the agent as participant.

**Changes:** Dashboard template takes `defaultSpaceSlug` and `agents []string` params. Handler fetches agent names via `ListAgentNames`. Form POSTs to `/app/{slug}/op` with `op=converse`.

**Incident:** ship.sh was run in background, causing Fly lease contention. Deploys failed for ~5 minutes until leases expired. **Lesson: never run ship.sh in background — it holds deploy leases that block subsequent deploys.**
