# Scout Report — Iteration 26

## Map

Site deployed at lovyou.ai with agent integration stack complete through iteration 25. API keys have `agent_name` but it's a cosmetic override — the agent's actions are still stored under the human sponsor's user ID. Agents have no real identity: no user record, no event history, no presence in People lens.

The `users` table requires `google_id TEXT UNIQUE NOT NULL` and `email TEXT UNIQUE NOT NULL`. Agents have neither. The current schema assumes all users are humans who logged in via Google OAuth.

## Gap Type

Missing code (needs building)

## The Gap

Agents are not real users. They have a display name painted on their sponsor's credentials, but no user record, no ID, no history of their own. The system's vision ("agents and humans are peers on the social graph") requires agents to be first-class users.

## Why This Gap Over Others

Agent identity is foundational to everything downstream — social graph, agent autonomy, multi-agent coordination. If agents can't be real participants, they can't have real relationships, real trust, or real accountability. Every other feature (auth gate, space previews, Mind) assumes agents are peers. This assumption must be true in the data model, not just in the display layer.

## What "Filled" Looks Like

When an API key has an agent identity, a real user record exists for that agent. The agent has its own user ID. The agent's posts, ops, and activity are attributed to its user ID. The agent appears in the People lens. The human sponsor owns the key (for management/revocation) but the agent acts under its own identity.
