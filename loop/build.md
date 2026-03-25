# Build: Agent Memory Phase 4 — persistent memory extraction and injection

## What Changed

### `site/graph/mind.go`

1. **`extractAndSaveMemories()`** — new method. After each agent reply, called asynchronously in a goroutine. Makes a lightweight Claude call with the last 6 messages (3 exchanges), extracts up to 3 durable facts/preferences about the user, stores each via `RememberForPersona` with kind="fact" and importance=7.

2. **`replyTo()`** — replaced generic `"Spoke with this user in conversation"` memory save with `go m.extractAndSaveMemories(...)`. The extraction now runs asynchronously so it doesn't block the user's reply.

3. **`buildSystemPrompt()`** — changed `RecallForPersona` limit from 10 to 5 (task specifies top 5 memories).

### `site/graph/memory_test.go`

4. **`TestBuildSystemPromptInjectsMemories`** — fixed broken SQL insert. Was using `(id, role, prompt, created_at)` but `agent_personas` has `name` not `role`. Fixed to use `store.UpsertAgentPersona()` with proper fields. Dropped unused `db` variable (`_, store := testDB(t)`).

## Already In Place (no changes needed)

- `agent_memories` table with `kind`, `source_id`, `importance` columns — in `migrate()` schema
- `RememberForPersona()` and `RecallForPersona()` methods in `store.go`
- Memory injection in `buildSystemPrompt()` — `== MEMORIES ==` section
- `TestRememberAndRecallForPersona`, `TestRememberForPersonaDefaults`, `TestRememberForPersonaInvalidKind` in `memory_test.go`

## Verification

- `go build -buildvcs=false ./...` — clean, no errors
- `go test ./...` — all pass (memory tests skip without DATABASE_URL; unit tests for buildSystemPrompt, extraction pass)
