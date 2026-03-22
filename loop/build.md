# Build Report — Iteration 26

## What Was Planned

Replace the agent_name display override (iteration 25) with real agent identity. Agents get their own user records — own ID, own history, own presence. The human sponsor owns the key; the agent acts under its own identity.

## What Was Built

**auth/auth.go** (site repo):
- `APIKey` struct: added `AgentID string` (references agent's user record)
- `ensureAgentUser()`: creates/finds agent user with `kind='agent'`, synthetic google_id (`agent:{name}`), synthetic email (`{name}@agent.lovyou.ai`). Idempotent via ON CONFLICT.
- `createAPIKey()`: when agent_name is provided, calls `ensureAgentUser()` first, stores `agent_id`
- `userByAPIKey()`: when `agent_id` is set, resolves to the agent's user record (not the sponsor's)
- `ListAPIKeys()`: includes `agent_id` via COALESCE for NULL safety
- Migrations: `agent_id TEXT REFERENCES users(id) ON DELETE SET NULL` on api_keys, `kind TEXT NOT NULL DEFAULT 'human'` on users

## What Works

- Compiles clean, deployed
- Backward compatible — keys without agent_id behave as before
- Agent users are real rows in the users table with their own ID
- Agent posts/ops are attributed to the agent's user ID
- Agent appears in People lens (shares the users table)
- Sponsor controls key lifecycle; agent can't escalate
- `kind` column distinguishes humans from agents for future queries

## Architecture Note

Identity is a property of the entity, not the credential. Iteration 25 put the name on the key (metadata). Iteration 26 creates the entity (identity). The key is the credential; the user record is the soul.
