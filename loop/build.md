# Build Report — Wire memory into auto-reply handler

## What changed

### `site/graph/store.go`
Added after `RecallForPersona`:
- `const memoryPersonaUser = "__user__"` — sentinel persona value for user-level memories not tied to any agent persona
- `RememberForUser(ctx, userID, kind, content, sourceID, importance)` — thin wrapper around `RememberForPersona` using the sentinel
- `RecallForUser(ctx, userID, limit)` — thin wrapper around `RecallForPersona` using the sentinel

### `site/graph/mind.go`
Modified `replyTo`:
- Moved `humanUserID` computation **before** the Claude call (was previously computed near the end of the function for persona memory saves)
- After `buildSystemPrompt`, calls `RecallForUser` for the human participant and appends a `== WHAT YOU REMEMBER ABOUT THIS USER ==` section to the system prompt if any memories exist
- After the reply is stored, launches `extractAndSaveUserMemories` as a goroutine (in addition to the existing persona-specific `extractAndSaveMemories`)

Added `extractAndSaveUserMemories(humanUserID, convoID, messages, agentID)`:
- Takes the last 6 messages, calls Claude with the prompt from the task spec ("Extract up to 3 facts worth remembering from this exchange as JSON array of {content, kind, importance}")
- Parses the JSON response, clamps importance to 1–5, validates kind, stores each fact via `RememberForUser`
- Mirrors the structure of the existing `extractAndSaveMemories` but uses `RememberForUser` instead of `RememberForPersona`

### `site/graph/memory_test.go`
Added:
- `TestRememberAndRecallForUser` — stores two memories with different importance levels, verifies ordering
- `TestRememberForUserDoesNotLeakAcrossUsers` — verifies user A's memories don't appear for user B

### `site/graph/mind_test.go`
Added:
- `TestReplyToInjectsUserMemories` — stores a user memory, calls `replyTo` with a `callClaudeOverride` that captures the system prompt, verifies the memory appears under the `== WHAT YOU REMEMBER ABOUT THIS USER ==` section

## Build verification
- `go build -buildvcs=false ./...` — passes (exit 0)
- `go test ./...` — passes (exit 0); new tests skip gracefully when DATABASE_URL is unset, consistent with all other DB tests in the package
